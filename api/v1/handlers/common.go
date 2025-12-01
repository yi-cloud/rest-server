package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GinResponseOk(c *gin.Context, ret interface{}) {
	if ret != nil {
		c.JSON(http.StatusOK, ret)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func GinResponseData(c *gin.Context, ret interface{}, err error, statusCode ...int) {
	var sc int
	if len(statusCode) > 0 {
		sc = statusCode[0]
	} else {
		sc = http.StatusOK
	}

	if err != nil {
		ret = gin.H{"error": err.Error()}
		if sc <= 0 {
			sc = http.StatusBadRequest
		}
	}

	c.JSON(sc, ret)
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

func ParseBool(c *gin.Context, key, value string) (*bool, error) {
	ret := true
	switch value {
	case "true", "True", "1", "TRUE":
		return &ret, nil
	case "false", "False", "0", "FALSE":
		ret = false
		return &ret, nil
	case "":
		return nil, nil
	default:
		return nil, errors.New("bool type only is: true/True/1/TRUE or false/False/0/FALSE or None")
	}
}

func ParseUintParam(c *gin.Context, key string) (uint, error) {
	return ParseUint(c, key, c.Param(key))
}

func ParseUintQuery(c *gin.Context, key string) (uint, error) {
	return ParseUint(c, key, c.Query(key))
}

func ParseBoolParam(c *gin.Context, key string) *bool {
	b, _ := ParseBool(c, key, c.Param(key))
	return b
}

func ParseBoolQuery(c *gin.Context, key string) *bool {
	b, _ := ParseBool(c, key, c.Query(key))
	return b
}
