package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/nimaibhat/GoCrudBookManagement/response"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware to handle panics and respond with a 500 status code
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				resp := response.NewResponse(map[string]interface{}{"error": "Internal Server Error"}, http.StatusInternalServerError)
				resp.SendResponse(c.Writer)
				c.Abort()
			}
		}()
		c.Next()
	}
}
