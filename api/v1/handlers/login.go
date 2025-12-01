package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/yi-cloud/rest-server/pkg/config"
	"github.com/yi-cloud/rest-server/pkg/middleware"
	"net/http"
)

type LoginParams struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegenerateToken(tokenString string, roles []any) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return middleware.PublicKey, nil
	})

	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["roles"] = roles

	token = jwt.NewWithClaims(middleware.SigningMethod, claims)
	return token.SignedString(middleware.PrivateKey)
}

func Login(c *gin.Context) {
	var req LoginParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error

	// request keystone to get token
	osClient, err := openstack.AuthenticatedClient(context.TODO(), gophercloud.AuthOptions{
		IdentityEndpoint: config.KeystoneOpt.EndPoint,
		Username:         req.Name,
		Password:         req.Password,
		DomainName:       config.KeystoneOpt.DomainName,
		Scope:            config.GetAuthScope(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := osClient.GetAuthResult().(tokens.CreateResult)
	if result.Err != nil {
		GinResponseData(c, nil, result.Err, result.StatusCode)
		return
	}

	resp := MakeLoginResponse(result.Body)
	if config.Regenerate {
		resp.Token, err = RegenerateToken(osClient.TokenID, resp.Roles)
		if err != nil {
			result.StatusCode = http.StatusInternalServerError
		}
	} else {
		resp.Token = osClient.TokenID
	}
	GinResponseData(c, resp, err, result.StatusCode)
}
