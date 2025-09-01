package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PaginatedResponse represents paginated response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// Created sends a created response
func Created(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusCreated, response)
}

// BadRequest sends a bad request response
func BadRequest(c *gin.Context, message string, err interface{}) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
	c.JSON(http.StatusBadRequest, response)
}

// NotFound sends a not found response
func NotFound(c *gin.Context, message string) {
	response := APIResponse{
		Success: false,
		Message: message,
	}
	c.JSON(http.StatusNotFound, response)
}

// InternalError sends an internal server error response
func InternalError(c *gin.Context, message string, err interface{}) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
	c.JSON(http.StatusInternalServerError, response)
}

// Conflict sends a conflict response
func Conflict(c *gin.Context, message string) {
	response := APIResponse{
		Success: false,
		Message: message,
	}
	c.JSON(http.StatusConflict, response)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, message string, items interface{}, total int64, page, pageSize int) {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	paginatedData := PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    paginatedData,
	}
	c.JSON(http.StatusOK, response)
}
