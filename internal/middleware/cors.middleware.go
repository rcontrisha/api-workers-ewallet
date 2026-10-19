package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func Cors(c *gin.Context) {
	// simple cors
	allowedOrigins := []string{"http://localhost", "http://localhost:5173", "https://recordings-hammer-destination-township.trycloudflare.com"}
	if slices.Contains(allowedOrigins, c.GetHeader("Origin")) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	}
	c.Header("Access-Control-Allow-Headers", "Content-Type, XXX-Header, Authorization")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH, DELETE")

	// preflight cors
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
