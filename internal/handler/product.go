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

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Register(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var createProductDTO dtos.CreateProductDTO
	if err := c.ShouldBindJSON(&createProductDTO); err != nil {
		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), createProductDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.productService.ListProducts(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product ID"))
		return
	}

	product, err := h.productService.ViewProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product ID"))
		return
	}

	var updateProductDTO dtos.UpdateProductDTO
	if err := c.ShouldBindJSON(&updateProductDTO); err != nil {
		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, updateProductDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product ID"))
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}
