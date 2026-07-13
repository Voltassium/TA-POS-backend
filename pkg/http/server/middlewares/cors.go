package middlewares

import "github.com/gin-gonic/gin"

func HandleCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:5173":      true,
			"http://localhost:4173":      true,
			"https://selipos.vercel.app": true,
			"https://ta-pos-frontend-fcbc4pdu6-dimas-premonos-projects.vercel.app": true,
		}

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Content-Length, Access-Control-Request-Method, Access-Control-Request-Headers, ngrok-skip-browser-warning")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
