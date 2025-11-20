package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/pkg/logs"
	"io"
	"net/http"
)

func PrintRawBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength <= 0 {
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("read request body error: %v", err)})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		logs.Logger.Infoln(string(body))
	}
}
