package repository

import (
	"context"
	"mazekav/internal/entity"
)

type CommonBehaviourRepository[T entity.DBModel] interface {
	ByID(ctx context.Context, id uint) (T, error)
	ByField(ctx context.Context, field string, id uint) (T, error)
	Save(ctx context.Context, model *T) error
	// add more common behaviour
}

type UserRepository interface {
	CommonBehaviourRepository[entity.User]
}

type RestaurantRepository interface {
	CommonBehaviourRepository[entity.Restaurant]
	Nearby(ctx context.Context, lat, long, radius float64) ([]entity.Restaurant, error)
}
