#!/bin/bash
set -e
set -x

# Config values
APP_NAME="ecom-go"
NAMESPACE="ecom-dev"
DB_NAME="ecomdb"
DB_USER="ecomuser"
DB_PASSWORD="ecompassword"

STORAGE_CLASS="standard"  # Minikube's default

# Edit these values directly if needed

echo "=== Deploying $APP_NAME to Minikube ==="
echo "Namespace: $NAMESPACE"

# Base64 encode the DB credentials
DB_USER_BASE64=$(echo -n $DB_USER | base64)
DB_PASSWORD_BASE64=$(echo -n $DB_PASSWORD | base64)

# Ensure Minikube is running
echo "Checking Minikube status..."
minikube status || minikube start

# Set docker env to use Minikube's Docker daemon
echo "Setting Docker environment to Minikube..."
eval $(minikube docker-env --shell=bash)

echo "Setting kubectl context to Minikube..."
kubectl config use-context minikube

# Create namespace if it doesn't exist
kubectl get namespace $NAMESPACE >/dev/null 2>&1 || kubectl create namespace $NAMESPACE

# Build Docker images
echo "Building API Docker image..."
docker build -t $APP_NAME-api:latest -f ./docker/api/Dockerfile .

echo "Building Worker Docker image..."
docker build -t $APP_NAME-worker:latest -f ./docker/worker/Dockerfile .

# Process and apply each YAML file
echo "Deploying infrastructure components..."

# Array of template files to process
TEMPLATES=(
  "infra/k8s/db.yaml"
  "infra/k8s/redis.yaml"
  "infra/k8s/rabbitmq.yaml"
  "infra/k8s/configmap.yaml"
  "infra/minikube/api.yaml"
  "infra/minikube/worker.yaml"
)

# Process each template
for template in "${TEMPLATES[@]}"; do
  echo "Applying $template..."
  sed -e "s/\${NAMESPACE}/$NAMESPACE/g" \
      -e "s/\${DB_NAME}/$DB_NAME/g" \
      -e "s/\${DB_USER_BASE64}/$DB_USER_BASE64/g" \
      -e "s/\${DB_PASSWORD_BASE64}/$DB_PASSWORD_BASE64/g" \
      -e "s/\${APP_NAME}/$APP_NAME/g" \
      -e "s#\${STORAGE_CLASS}#$STORAGE_CLASS#g" \
      "$template" | kubectl apply -f -
done

# Get service URL
echo "Waiting for service to get external IP..."
sleep 10
minikube service api-service -n $NAMESPACE --url

echo "=== Deployment to Minikube completed ==="
echo "To access the API, use the URL above"