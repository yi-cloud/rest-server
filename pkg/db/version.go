package db

import (
	"errors"
	"gorm.io/gorm"
)

var DBVersion string

type Version struct {
	gorm.Model
	Version string `json:"version" gorm:"not null;size:16"`
}

type VersionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) *VersionRepository {
	return &VersionRepository{db: db}
}

func (r *VersionRepository) IsAutoMigrate(oldVersion *Version) (bool, error) {
	var version Version
	err := r.db.AutoMigrate(&Version{})
	if err != nil {
		panic("auto migrate <version> table failed")
	}
	err = r.db.First(&version).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		version.ID = 1
		version.Version = "v1.0.0"
		err = r.db.Create(&version).Error
		return true, err
	}

	if err == nil && DBVersion != "" && version.Version != DBVersion {
		if oldVersion != nil {
			oldVersion.ID = version.ID
			oldVersion.Version = version.Version
		}
		version.Version = DBVersion
		err = r.db.Save(&version).Error
		return true, err
	}

	return false, err
}
