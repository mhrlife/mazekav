package repository

import (
	"gorm.io/gorm"
	"mazekav/internal/entity"
)

type usersRepository struct {
	CommonBehaviourRepository[entity.User]
}

func NewUsersRepository(db *gorm.DB) UserRepository {
	return &usersRepository{
		CommonBehaviourRepository: NewCommonBehaviour[entity.User](db),
	}
}
