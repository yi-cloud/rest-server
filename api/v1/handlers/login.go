package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/yi-cloud/rest-server/api/v1/services"
	"github.com/yi-cloud/rest-server/common"
	"github.com/yi-cloud/rest-server/models"
	"github.com/yi-cloud/rest-server/pkg/config"
	"github.com/yi-cloud/rest-server/pkg/db"
	"github.com/yi-cloud/rest-server/pkg/middleware"
	"net/http"
	"strings"
	"time"
)

type LoginParams struct {
	Name        string `json:"name"`
	MobilePhone string `json:"mobilePhone" binding:"required,min=11"`
	VerifyCode  string `json:"verifyCode" binding:"required,min=6,max=6"`
	NickName    string `json:"nickName"`
}

type DoctorLoginParams struct {
	LoginParams
	Hospital   string `json:"hospital" binding:"required"`
	IsTop3     *bool  `json:"isTop3"`
	JobTitle   string `json:"jobTitle" binding:"required"`
	Department string `json:"department"`
	AdeptField string `json:"adeptField"`
}

func GenerateToken(uid uint, name, mobilePhone, role string) (string, string, error) {
	now := time.Now()
	expSecond := config.TokenLifeTime
	if expSecond <= 0 {
		expSecond = 1800
	}
	expire := now.Add(time.Duration(expSecond) * time.Second)
	token := jwt.NewWithClaims(middleware.SigningMethod, jwt.MapClaims{
		"iss":   "rest-server",
		"iat":   now.Unix(),
		"exp":   expire.Unix(),
		"uid":   uint64(uid),
		"aud":   name,
		"role":  role,
		"phone": mobilePhone,
	})

	ret, err := token.SignedString(middleware.PrivateKey)
	return expire.Format(common.TimeFormat), ret, err
}

func Login(c *gin.Context) {
	var req LoginParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NickName == "" {
		req.NickName = req.MobilePhone
	}
	if req.Name == "" {
		req.Name = strings.ToLower(
			common.RandomStringWithSuffix("_"+req.MobilePhone[len(req.MobilePhone)-4:], 8))
	}

	var ret any
	var err error
	var uid uint
	var name string
	var mobile string

	ret, err = services.NewUserService(
		models.NewUserRepository(db.DBInstance())).GetUser(req.Name, req.MobilePhone, req.NickName)
	uid = ret.(*models.User).ID
	name = ret.(*models.User).Name
	mobile = ret.(*models.User).MobilePhone

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expire, token, err := GenerateToken(uid, name, mobile, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]any{"token": token, "expire": expire, "data": ret})
}
