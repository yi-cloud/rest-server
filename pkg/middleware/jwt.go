package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
				logs.Logger.Debug("cookies: %v", cookie)
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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return RsaPublicKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userId", uint64(claims["uid"].(float64)))
		c.Set("userName", claims["aud"].(string))
		c.Set("mobilePhone", claims["phone"].(string))
		c.Set("userRole", claims["role"].(string))

		urole := c.GetString("userRole")
		logs.Logger.Debugf("user role: %s", urole)
		c.Next()
	}
}
