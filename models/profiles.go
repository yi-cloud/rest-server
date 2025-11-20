package models

import (
	"errors"
	"github.com/yi-cloud/rest-server/common"
	"github.com/yi-cloud/rest-server/pkg/db"
	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserId            uint          `json:"userId" gorm:"primarykey;not null"`
	Name              string        `json:"name" gorm:"size:255"`
	Height            uint8         `json:"height"`
	Weight            float32       `json:"weight"`
	Age               common.MyTime `json:"age"`
	Gender            *uint8        `json:"gender"`
	PersonalSignature string        `json:"personalSignature" gorm:"size:40"`
}

func (UserProfile) SetId(p interface{}, val uint) {
	p.(*UserProfile).ID = val
}

func (UserProfile) SetUserId(p interface{}, val uint) {
	p.(*UserProfile).UserId = val
}

func (UserProfile) SetUserName(p interface{}, val string) {
	return
}

type UserProfileRepository struct {
	db *db.DBManager
}

func NewUserProfileRepository(db *db.DBManager) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) FindOne(userId uint) (*UserProfile, error) {
	var data *UserProfile
	err := r.db.GormDB().Model(UserProfile{}).Where("user_id = ?", userId).Scan(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return data, err
}

func init() {
	db.DBInstance().AddTabModel(&UserProfile{})
}
