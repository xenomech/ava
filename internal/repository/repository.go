package repository

import (
	sessionrepo "ava/internal/repository/session"
	userrepo "ava/internal/repository/user"

	"gorm.io/gorm"
)

type Repository struct {
	User    userrepo.Repository
	Session sessionrepo.Repository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:    userrepo.NewRepository(db),
		Session: sessionrepo.NewRepository(db),
	}
}
