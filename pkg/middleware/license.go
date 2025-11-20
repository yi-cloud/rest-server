package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/pkg/license"
	"net/http"
)

func CheckLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch license.CheckResult {
		case license.IsNone:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "license is none"})
			return
		case license.Expire:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "license is expire"})
			return
		case license.DecodeError:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "license decode failed"})
			return
		case license.ProductException:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "license is not special product"})
			return
		case license.SNError:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "license sn error"})
			return
		case license.ExceedClusters:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "license manage clusters is exceeded."})
			return
		}
	}
}
