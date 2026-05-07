package middlewares

import (
	"log"
	"time"

	"hotel-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			log.Printf("Error occurred: %v", err)

			var appErr *errors.AppError
			switch e := err.(type) {
			case *errors.AppError:
				appErr = e
			default:
				appErr = errors.ErrInternalServer
			}

			c.JSON(appErr.Code, ErrorResponse{
				Code:      appErr.Code,
				Message:   appErr.Message,
				Timestamp: time.Now(),
			})
		}
	}
}
