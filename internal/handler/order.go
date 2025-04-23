package handler

import (
	"ecom-go/internal/dtos"
	"net/http"
	"strconv"

	"ecom-go/internal/service"
	"ecom-go/pkg/errors"
	"ecom-go/pkg/http/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) Register(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrder)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var createOrderDTO dtos.CreateOrderDTO
	if err := c.ShouldBindJSON(&createOrderDTO); err != nil {
		response.Error(c, errors.NewBadRequestError("Invalid request payload", err))
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), &createOrderDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {

	orders, err := h.orderService.ListOrders(c.Request.Context(), 0, 10)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.NewBadRequestError("Invalid order ID", err))
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, order)
}
