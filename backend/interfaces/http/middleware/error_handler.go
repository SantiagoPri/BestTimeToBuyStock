package middleware

import (
	"backend/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware to handle application errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var status int
			var message string

			switch errors.GetCode(err) {
			case errors.ErrNotFound:
				status = http.StatusNotFound
				message = err.Error()
			case errors.ErrInvalidInput:
				status = http.StatusBadRequest
				message = err.Error()
			case errors.ErrUnauthorized:
				status = http.StatusUnauthorized
				message = err.Error()
			case errors.ErrForbidden:
				status = http.StatusForbidden
				message = err.Error()
			case errors.ErrConflict:
				status = http.StatusConflict
				message = err.Error()
			case errors.ErrNotAvailable:
				status = http.StatusServiceUnavailable
				message = err.Error()
			default:
				status = http.StatusInternalServerError
				message = "Internal server error"
			}

			c.JSON(status, gin.H{"error": message})
			c.Abort()
		}
	}
}
