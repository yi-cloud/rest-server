package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/yi-cloud/rest-server/pkg/config"
	"net/http"
	"strconv"
	"strings"
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
		if respErr, ok := err.(gophercloud.ErrUnexpectedResponseCode); ok {
			var respBody map[string]any
			sc = respErr.GetStatusCode()
			if err := json.Unmarshal(respErr.Body, &respBody); err != nil {
				ret = string(respErr.Body)
			} else {
				ret = respBody
			}
		} else {
			ret = gin.H{"error": err.Error()}
		}
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

func NewIdentityV3Client(c *gin.Context) (*gophercloud.ServiceClient, error) {
	osClient, err := openstack.NewClient(config.KeystoneOpt.EndPoint)
	if err != nil {
		return nil, err
	}

	osClient.SetToken(c.GetString("token"))
	return openstack.NewIdentityV3(osClient, gophercloud.EndpointOpts{})
}

func HasRoleassignment(roles []any, serviceName string) bool {
	for _, role := range roles {
		if _role, ok := role.(map[string]any); ok {
			roleName := _role["name"].(string)
			if strings.HasPrefix(roleName, serviceName) {
				return true
			}
		}
	}
	return false
}

func IsAdminRole(roles []any) bool {
	for _, role := range roles {
		if _role, ok := role.(map[string]any); ok {
			name := _role["name"].(string)
			if name == "admin" {
				return true
			}
		}
	}
	return false
}

func MakeLoginResponse(body any) (resp TokenResponse) {
	if _body, ok := body.(map[string]any); ok {
		token := _body["token"]
		if _body, ok := token.(map[string]any); ok {
			if user, ok := _body["user"].(map[string]any); ok {
				resp.User = user
			}
			if project, ok := _body["project"].(map[string]any); ok {
				resp.Project = project
			}
			if roles, ok := _body["roles"].([]any); ok {
				resp.Roles = roles
				resp.IsAdmin = IsAdminRole(roles)
			}

			if expireAt, ok := _body["expires_at"].(any); ok {
				resp.ExpiresAt = expireAt
			}

			if issuedAt, ok := _body["issued_at"].(any); ok {
				resp.IssuedAt = issuedAt
			}

			if catalogs, ok := _body["catalog"].([]any); ok {
				for _, catalog := range catalogs {
					if _catalog, ok := catalog.(map[string]any); ok {
						if _type, ok := _catalog["type"].(string); ok {
							serviceName := _catalog["name"].(string)
							if _type != "identity" && (resp.IsAdmin || HasRoleassignment(resp.Roles, serviceName)) {
								resp.Catalogs = append(resp.Catalogs, catalog)
							}
						}
					}
				}
			}
		}
	}
	return
}
