package middleware

import (
	"cloudpan/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := model.GetSession(c.Request)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{"msg": "Get session error"})
			c.Abort()
			return
		}

		admin := session.Values["user"]
		if admin == nil {
			c.JSON(http.StatusAccepted, gin.H{"msg": "Unauthorized access!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
