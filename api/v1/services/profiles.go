package services

import (
	"github.com/yi-cloud/rest-server/common"
	"github.com/yi-cloud/rest-server/models"
	"time"
)

type UserProfileService struct {
	profileRepo *models.UserProfileRepository
}

func NewUserProfileService(profileRepo *models.UserProfileRepository) *UserProfileService {
	return &UserProfileService{profileRepo: profileRepo}
}

type ProfileResponse struct {
	models.UserProfile
	AgeCount int `json:"ageCount"`
}

func (s *UserProfileService) GetUserProfile(userId uint) (any, error) {
	profile, err := s.profileRepo.FindOne(userId)
	if err == nil && profile != nil {
		resp := &ProfileResponse{
			UserProfile: *profile,
		}
		ageStr := profile.Age.String()
		if ageStr == "0001-01-01 00:00:00" {
			resp.AgeCount = 0
		} else {
			age, err := time.Parse(common.TimeFormat, ageStr)
			if err == nil {
				resp.AgeCount = time.Now().Year() - age.Year()
			}
		}

		return resp, nil
	}

	return profile, err
}
