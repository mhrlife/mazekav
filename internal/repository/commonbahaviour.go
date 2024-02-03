package repository

import (
	"context"
	"gorm.io/gorm"
	"mazekav/internal/entity"
)

type commonBehaviour[T entity.DBModel] struct {
	db *gorm.DB
}

func NewCommonBehaviour[T entity.DBModel](db *gorm.DB) CommonBehaviourRepository[T] {
	return &commonBehaviour[T]{
		db: db,
	}
}

func (c *commonBehaviour[T]) ByID(ctx context.Context, id uint) (T, error) {
	return c.ByField(ctx, "id", id)
}

func (c *commonBehaviour[T]) ByField(ctx context.Context, field string, id uint) (T, error) {
	var t T
	err := c.db.WithContext(ctx).Where(field+"=?", id).First(&t).Error
	return t, err
}

func (c *commonBehaviour[T]) Save(ctx context.Context, model *T) error {
	return c.db.WithContext(ctx).Save(model).Error
}
