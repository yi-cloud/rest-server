package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/api/v1/services"
	"github.com/yi-cloud/rest-server/models"
	"github.com/yi-cloud/rest-server/pkg/db"
	"github.com/yi-cloud/rest-server/pkg/server"
)

type UserProfileHandler struct {
	profileService *services.UserProfileService
}

func NewUserProfileHandler(profileService *services.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{profileService: profileService}
}

func (h *UserProfileHandler) GetUserProfile(c *gin.Context) {
	profile, err := h.profileService.GetUserProfile(uint(c.GetUint64("userId")))
	GinResponseData(c, profile, err)
}

func GetUserProfileInstance(db *db.DBManager) *UserProfileHandler {
	return NewUserProfileHandler(services.NewUserProfileService(models.NewUserProfileRepository(db)))
}

func AddUserProfileRoutes(s *server.ApiServer) {
	g := s.AddGroup("/user-profiles", nil)
	s.AddRoute(g, "GET", "/:id", GetUserProfileInstance(db.DBInstance()).GetUserProfile)
}
