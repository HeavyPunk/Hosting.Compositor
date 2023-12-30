package middleware_logging

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorLogger(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		fmt.Printf("Caughed error: %v\n", err)
	}

	c.JSON(http.StatusInternalServerError, "")
}
