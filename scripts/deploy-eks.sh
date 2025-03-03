#!/bin/bash
set -e

# Config values
CLUSTER_NAME="ecom-cluster"
REGION="ap-southeast-1"
ECR_REPOSITORY="ecom-go"
NAMESPACE="ecom-prod"
DB_NAME="ecomdb"
DB_USER="ecomuser"
DB_PASSWORD="ecompassword"

STORAGE_CLASS="gp2"

# Clean up
if [ "$1" == "--cleanup" ]; then
  echo "Cleaning up previous deployment..."
  eksctl delete cluster --name $CLUSTER_NAME --region $REGION
  exit 0
fi

echo "=== Creating EKS cluster: $CLUSTER_NAME ==="
eksctl create cluster \
  --name $CLUSTER_NAME \
  --region $REGION \
  --nodegroup-name my-nodes \
  --node-type t3.small \
  --nodes 2

echo "=== Ensuring AWS EBS CSI driver is installed ==="
eksctl create addon --name aws-ebs-csi-driver --cluster $CLUSTER_NAME --region $REGION --force || echo "EBS CSI driver already installed"

echo "=== Attaching IAM policy to allow EBS volume creation ==="
NODE_ROLE=$(aws iam list-roles --query "Roles[?contains(RoleName, 'eksctl-$CLUSTER_NAME-nodegroup')].RoleName" --output text)

if [[ -z "$NODE_ROLE" ]]; then
  echo "ERROR: Could not find node instance role for EKS nodes."
  exit 1
fi

echo "Node instance role: $NODE_ROLE"

aws iam attach-role-policy \
    --role-name "$NODE_ROLE" \
    --policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy

echo "Waiting for IAM policy attachment to propagate..."
sleep 30  # Give AWS time to propagate changes

echo "=== Updating kubeconfig for EKS cluster ==="
aws eks update-kubeconfig --region $REGION --name $CLUSTER_NAME

echo "=== Creating namespace if not exists ==="
kubectl get namespace $NAMESPACE >/dev/null 2>&1 || kubectl create namespace $NAMESPACE

# Base64 encode DB credentials
DB_USER_BASE64=$(echo -n $DB_USER | base64)
DB_PASSWORD_BASE64=$(echo -n $DB_PASSWORD | base64)

# Get ECR repository URI
echo "=== Setting up Amazon ECR ==="
ECR_REPO_URI=$(aws ecr describe-repositories --repository-names $ECR_REPOSITORY --region $REGION --query 'repositories[0].repositoryUri' --output text || \
  aws ecr create-repository --repository-name $ECR_REPOSITORY --region $REGION --query 'repository.repositoryUri' --output text)

# Login to ECR
echo "Logging in to Amazon ECR..."
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_REPO_URI

# Build and push Docker images
echo "Building and pushing API image..."
docker build -t $ECR_REPO_URI:api-latest -f ./docker/api/Dockerfile .
docker push $ECR_REPO_URI:api-latest

echo "Building and pushing Worker image..."
docker build -t $ECR_REPO_URI:worker-latest -f ./docker/worker/Dockerfile .
docker push $ECR_REPO_URI:worker-latest

# Deploy infrastructure (DB, Redis, RabbitMQ, ConfigMap)
echo "=== Deploying infrastructure components ==="
TEMPLATES=(
  "infra/k8s/db.yaml"
  "infra/k8s/redis.yaml"
  "infra/k8s/rabbitmq.yaml"
  "infra/k8s/configmap.yaml"
)

for template in "${TEMPLATES[@]}"; do
  echo "Applying $template..."
  sed -e "s#\${NAMESPACE}#$NAMESPACE#g" \
      -e "s#\${DB_NAME}#$DB_NAME#g" \
      -e "s#\${DB_USER_BASE64}#$DB_USER_BASE64#g" \
      -e "s#\${DB_PASSWORD_BASE64}#$DB_PASSWORD_BASE64#g" \
      -e "s#\${STORAGE_CLASS}#$STORAGE_CLASS#g" \
      "$template" | kubectl apply -f -
done

# Wait for PVC to be bound before deploying API/Worker
echo "=== Waiting for PostgreSQL PVC to be bound ==="
PVC_NAME="postgres-pvc"
for i in {1..12}; do
  PVC_STATUS=$(kubectl get pvc $PVC_NAME -n $NAMESPACE -o jsonpath='{.status.phase}' 2>/dev/null || echo "NotFound")
  if [[ "$PVC_STATUS" == "Bound" ]]; then
    echo "PVC $PVC_NAME is bound!"
    break
  fi
  echo "Waiting for PVC $PVC_NAME to bind... ($i/12)"
  sleep 10
done

# Wait for Postgres to be ready
echo "=== Waiting for PostgreSQL pod to be ready ==="
for i in {1..12}; do
  POSTGRES_READY=$(kubectl get pods -n $NAMESPACE -l app=postgres -o jsonpath='{.items[0].status.containerStatuses[0].ready}' 2>/dev/null || echo "false")
  if [[ "$POSTGRES_READY" == "true" ]]; then
    echo "PostgreSQL pod is ready!"
    break
  fi
  echo "Waiting for PostgreSQL pod to be ready... ($i/12)"
  sleep 10
done

# Deploy application components (API & Worker)
echo "=== Deploying API & Worker ==="
TEMPLATES=(
  "infra/eks/api.yaml"
  "infra/eks/worker.yaml"
  "infra/eks/autoscaling.yaml"
)

for template in "${TEMPLATES[@]}"; do
  echo "Applying $template..."
  sed -e "s#\${NAMESPACE}#$NAMESPACE#g" \
      -e "s#\${ECR_REPO_URI}#$ECR_REPO_URI#g" \
      "$template" | kubectl apply -f -
done

# Wait for API LoadBalancer IP
echo "=== Waiting for API LoadBalancer IP ==="
for i in {1..15}; do
  EXTERNAL_IP=$(kubectl get service api-service -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
  if [[ -n "$EXTERNAL_IP" ]]; then
    echo "API service is available at: http://$EXTERNAL_IP"
    break
  fi
  echo "Waiting for external IP... ($i/15)"
  sleep 10
done

if [[ -z "$EXTERNAL_IP" ]]; then
  echo "=== API LoadBalancer address not available yet. Check with: kubectl get service api-service -n $NAMESPACE ==="
fi

echo "=== Deployment completed successfully! ==="
echo "To view pod status: kubectl get pods -n $NAMESPACE"
echo "To view service status: kubectl get services -n $NAMESPACE"
echo "To view logs: kubectl logs -n $NAMESPACE deployment/api"
