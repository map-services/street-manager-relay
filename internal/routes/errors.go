package routes

import (
	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, code int, message string, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
