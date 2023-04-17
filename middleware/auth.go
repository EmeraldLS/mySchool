package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mySchool/token"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("token") == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unauthorized access. No access token",
			})
			c.Abort()
			return
		}
		if err := token.ValidateToken(c.GetHeader("token")); err != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("An error occured. Error('%v')", err),
			})
			c.Abort()
			return
		}
	}
}
