package controller

import (
	"errors"
	"go-clean-architecture/internal/entity"
	"go-clean-architecture/internal/usecase"
	"go-clean-architecture/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userUseCase *usecase.UserUseCase
}

// NewUserController creates a new user controller instance
func NewUserController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// CreateUser handles POST /users
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := ctrl.userUseCase.CreateUser(c.Request.Context(), &user); err != nil {
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExists):
			response.Conflict(c, "User with this email already exists")
		case errors.Is(err, entity.ErrInvalidUserName), errors.Is(err, entity.ErrInvalidUserEmail):
			response.BadRequest(c, "Invalid user data", err.Error())
		default:
			response.InternalError(c, "Failed to create user", err.Error())
		}
		return
	}

	response.Created(c, "User created successfully", user)
}

// GetUser handles GET /users/:id
func (ctrl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	user, err := ctrl.userUseCase.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			response.NotFound(c, "User not found")
		case errors.Is(err, entity.ErrInvalidUserID):
			response.BadRequest(c, "Invalid user ID", err.Error())
		default:
			response.InternalError(c, "Failed to retrieve user", err.Error())
		}
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

// GetAllUsers handles GET /users
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	users, total, err := ctrl.userUseCase.GetAllUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to retrieve users", err.Error())
		return
	}

	response.Paginated(c, "Users retrieved successfully", users, total, page, pageSize)
}

// UpdateUser handles PUT /users/:id
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := ctrl.userUseCase.UpdateUser(c.Request.Context(), uint(id), &user); err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			response.NotFound(c, "User not found")
		case errors.Is(err, entity.ErrInvalidUserID), errors.Is(err, entity.ErrInvalidUserName), errors.Is(err, entity.ErrInvalidUserEmail):
			response.BadRequest(c, "Invalid user data", err.Error())
		default:
			response.InternalError(c, "Failed to update user", err.Error())
		}
		return
	}

	response.Success(c, "User updated successfully", user)
}

// DeleteUser handles DELETE /users/:id
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	if err := ctrl.userUseCase.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			response.NotFound(c, "User not found")
		case errors.Is(err, entity.ErrInvalidUserID):
			response.BadRequest(c, "Invalid user ID", err.Error())
		default:
			response.InternalError(c, "Failed to delete user", err.Error())
		}
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// ActivateUser handles PUT /users/:id/activate
func (ctrl *UserController) ActivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	if err := ctrl.userUseCase.ActivateUser(c.Request.Context(), uint(id)); err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			response.NotFound(c, "User not found")
		default:
			response.InternalError(c, "Failed to activate user", err.Error())
		}
		return
	}

	response.Success(c, "User activated successfully", nil)
}

// DeactivateUser handles PUT /users/:id/deactivate
func (ctrl *UserController) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	if err := ctrl.userUseCase.DeactivateUser(c.Request.Context(), uint(id)); err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			response.NotFound(c, "User not found")
		default:
			response.InternalError(c, "Failed to deactivate user", err.Error())
		}
		return
	}

	response.Success(c, "User deactivated successfully", nil)
}
