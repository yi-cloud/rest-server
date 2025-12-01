package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/api/v1/services"
	"github.com/yi-cloud/rest-server/models"
	"gorm.io/gorm"
	"net/http"
)

type CommonHandler[T models.Common] struct {
	Service *services.CommonService[T]
}

func NewCommonHandler[T models.Common](service *services.CommonService[T]) *CommonHandler[T] {
	return &CommonHandler[T]{Service: service}
}

func (h *CommonHandler[T]) Find(c *gin.Context) {
	id, err := ParseUintParam(c, "id")
	if err != nil {
		return
	}
	ret, err := h.Service.Find(id)
	GinResponseData(c, ret, err)
}

func (h *CommonHandler[T]) Create(c *gin.Context) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.SetUserId(&req, uint(c.GetUint64("userId")))
	err := h.Service.Create(&req)
	GinResponseData(c, req, err)
}

func (h *CommonHandler[T]) Update(c *gin.Context) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := ParseUintParam(c, "id")
	if err != nil {
		return
	}
	req.SetId(&req, id)
	err = h.Service.Update(&req)
	GinResponseData(c, req, err)
}

func (h *CommonHandler[T]) Delete(c *gin.Context) {
	id, err := ParseUintParam(c, "id")
	if err != nil {
		return
	}
	GinResponseData(c, nil, h.Service.Delete(id))
}

func GetCommonInstance[T models.Common](db *gorm.DB) *CommonHandler[T] {
	return NewCommonHandler[T](services.NewCommonService[T](models.NewCommonRepository[T](db)))

}
