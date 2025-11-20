package services

import (
	"github.com/yi-cloud/rest-server/models"
)

type CommonService[T models.Common] struct {
	Repo *models.CommonRepository[T]
}

func NewCommonService[T models.Common](repo *models.CommonRepository[T]) *CommonService[T] {
	return &CommonService[T]{Repo: repo}
}

func (s *CommonService[T]) Find(id uint) (*T, error) {
	return s.Repo.Find(id)
}

func (s *CommonService[T]) Create(v *T) error {
	return s.Repo.Create(v)
}

func (s *CommonService[T]) Update(v *T) error {
	return s.Repo.Update(v)
}

func (s *CommonService[T]) Delete(id uint) error {
	return s.Repo.Delete(id)
}
