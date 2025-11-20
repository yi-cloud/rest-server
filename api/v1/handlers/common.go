package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GinResponseOk(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func GinResponseData(c *gin.Context, ret interface{}, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ret)
}

func ParseUint(c *gin.Context, key, value string) (uint, error) {
	if value == "" {
		errmsg := fmt.Sprintf("%s is null", key)
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg})
		return 0, errors.New(errmsg)
	}

	i32, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, err
	} else if i32 < 0 {
		errmsg := fmt.Sprintf("Invalid param id: %d", i32)
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg})
		return 0, errors.New(errmsg)
	}

	return uint(i32), nil
}

func ParseUintParam(c *gin.Context, key string) (uint, error) {
	return ParseUint(c, key, c.Param(key))
}

func ParseUintQuery(c *gin.Context, key string) (uint, error) {
	return ParseUint(c, key, c.Query(key))
}
