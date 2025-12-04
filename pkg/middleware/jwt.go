package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/pkg/logs"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "login") {
			return
		}
		logs.Logger.Debugf("start token auth")
		authString := c.GetHeader("Authorization")
		if len(authString) == 0 {
			token := c.Request.Form.Get("token")
			if len(token) > 0 {
				authString = "Bearer " + token
			}
		}

		if len(authString) == 0 {
			cookies := c.Request.Cookies()
			for _, cookie := range cookies {
				logs.Logger.Debugf("cookies: %v", cookie)
				if strings.Contains(cookie.Name, "_access_token") {
					authString = "Bearer " + cookie.Value
					break
				}
			}
		}

		kv := strings.Split(authString, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			logs.Logger.Info("AuthString invalid:", authString)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unavailable auth token"})
			return
		}
		tokenString := kv[1]
		logs.Logger.Debugf("Authorization token: %s", tokenString)
		c.Set("token", tokenString)
		c.Next()
	}
}
