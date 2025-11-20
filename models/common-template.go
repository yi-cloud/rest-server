package models

import (
	"gorm.io/gorm"
)

type Common interface {
	SetId(interface{}, uint)
	SetUserId(interface{}, uint)
	SetUserName(interface{}, string)
}

type CommonRepository[T Common] struct {
	db *gorm.DB
}

func NewCommonRepository[T Common](db *gorm.DB) *CommonRepository[T] {
	return &CommonRepository[T]{
		db: db,
	}
}

func (r *CommonRepository[T]) Find(id uint) (*T, error) {
	var t *T
	res := r.db.Where("id = ?", id).Find(&t)
	if res.RowsAffected > 0 {
		return t, res.Error
	}
	return nil, res.Error
}

func (r *CommonRepository[T]) Create(v *T) error {
	return r.db.Create(v).Error
}

func (r *CommonRepository[T]) Update(v *T) error {
	return r.db.Model(v).Updates(v).Error
}

func (r *CommonRepository[T]) Delete(id uint) error {
	var t T
	t.SetId(&t, id)
	return r.db.Delete(&t, "id", id).Error
}
