package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//跨域设置
func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			// AllowAllOrigins:  true,
			AllowOrigins: []string{"*"},
			// AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowMethods:  []string{"*"},
			AllowHeaders:  []string{"*"},
			ExposeHeaders: []string{"Content-Length", "text/plain", "Authorization", "Content-Type"},
			// AllowCredentials: true,
			MaxAge: 12 * time.Hour,
		},
	)

}
