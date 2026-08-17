package repository

import (
	membershiprepo "ava/internal/repository/membership"
	sessionrepo "ava/internal/repository/session"
	tenantrepo "ava/internal/repository/tenant"
	tokenrepo "ava/internal/repository/token"
	userrepo "ava/internal/repository/user"

	"gorm.io/gorm"
)

type Repository struct {
	User       userrepo.Repository
	Tenant     tenantrepo.Repository
	Membership membershiprepo.Repository
	Session    sessionrepo.Repository
	Token      tokenrepo.Repository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:       userrepo.NewRepository(db),
		Tenant:     tenantrepo.NewRepository(db),
		Membership: membershiprepo.NewRepository(db),
		Session:    sessionrepo.NewRepository(db),
		Token:      tokenrepo.NewRepository(db),
	}
}
