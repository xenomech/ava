package repository

import (
	userrepo "ava/internal/repository/user"

	"gorm.io/gorm"
)

type Repository struct {
	User userrepo.Repository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: userrepo.NewRepository(db),
	}
}
