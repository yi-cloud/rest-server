package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/api/v1/services"
	"github.com/yi-cloud/rest-server/models"
	"github.com/yi-cloud/rest-server/pkg/db"
	"github.com/yi-cloud/rest-server/pkg/server"
	"net/http"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	GinResponseData(c, users, err)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	mobile := c.Query("mobile")
	if mobile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mobile is null"})
		return
	}
	user, err := h.userService.GetUserByMobile(mobile)
	GinResponseData(c, user, err)
}

func GetUserInstance(db *db.DBManager) *UserHandler {
	return NewUserHandler(services.NewUserService(models.NewUserRepository(db)))
}

func AddUserRoutes(s *server.ApiServer) {
	g := s.AddGroup("/users", nil)
	s.AddRoute(g, "GET", "", GetUserInstance(db.DBInstance()).GetUsers)
	s.AddRoute(g, "GET", "/:id", GetUserInstance(db.DBInstance()).GetUser)
}
