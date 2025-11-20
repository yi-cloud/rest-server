package models

import (
	"github.com/yi-cloud/rest-server/pkg/db"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null;size:64;"`
	MobilePhone string `json:"mobilePhone" gorm:"primarykey;unique;not null;size:40"`
	NickName    string `json:"nickName" gorm:"size:64"`
	Avatar      string `json:"avatar" gorm:"null;size:255"`
	QrCode      string `json:"qrCode" gorm:"not null;size:255"`
	Description string `json:"description" gorm:"null"`
}

type UserRepository struct {
	db *db.DBManager
}

func NewUserRepository(db *db.DBManager) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]User, error) {
	var users []User
	err := r.db.GormDB().Find(&users).Error
	return users, err
}

func (r *UserRepository) FindOne(mobilePhone string) (*User, error) {
	user := User{MobilePhone: mobilePhone}
	err := r.db.GormDB().First(&user, user).Error
	return &user, err
}

func (r *UserRepository) FindById(id uint) (*User, error) {
	user := User{Model: gorm.Model{
		ID: id,
	}}
	err := r.db.GormDB().First(&user, user).Error
	return &user, err
}

func (r *UserRepository) Create(name, mobilePhone, nickName, qrCode string) (*User, error) {
	user := User{MobilePhone: mobilePhone, Name: name, NickName: nickName, QrCode: qrCode}
	r.db.GormDB().Begin()
	err := r.db.GormDB().Create(&user).Error
	if err != nil {
		r.db.GormDB().Rollback()
		return nil, err
	}

	profile := &UserProfile{UserId: user.ID}
	err = r.db.GormDB().Model(profile).Create(profile).Error
	if err != nil {
		r.db.GormDB().Rollback()
		return nil, err
	}

	r.db.GormDB().Commit()
	return &user, err
}

func (r *UserRepository) Updates(userId uint, values interface{}) (*User, error) {
	user := &User{}
	err := r.db.GormDB().Model(user).Where("id = ?", userId).Updates(values).Error
	return user, err
}

func init() {
	db.DBInstance().AddTabModel(&User{})
}
