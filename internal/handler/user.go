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

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register sets up routes for the user handler
func (h *UserHandler) Register(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("", h.List)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

// Create handles user creation
func (h *UserHandler) Create(c *gin.Context) {
	var createUserDTO dtos.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), createUserDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, user)
}

// GetByID handles retrieving a user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user ID"))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Update handles updating a user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user ID"))
		return
	}

	var updateUserDTO dtos.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), uint(id), updateUserDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Delete handles deleting a user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user ID"))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}

// List handles retrieving users with pagination
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	users, total, err := h.userService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, users, page, pageSize, total)
}
