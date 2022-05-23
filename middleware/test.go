package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func TestMiddleware() gin.HandlerFunc {
	// Do any setup work here
	return func(c *gin.Context) {
		fmt.Println("[TEST-MIDDLEWARE] We got the call")
		c.Next()
	}
}
