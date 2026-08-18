package repository

import (
	devicerepo "ava/api/internal/repository/device"
	flowrepo "ava/api/internal/repository/flow"
	membershiprepo "ava/api/internal/repository/membership"
	sessionrepo "ava/api/internal/repository/session"
	tenantrepo "ava/api/internal/repository/tenant"
	tokenrepo "ava/api/internal/repository/token"
	userrepo "ava/api/internal/repository/user"

	"gorm.io/gorm"
)

type Repository struct {
	User       userrepo.Repository
	Tenant     tenantrepo.Repository
	Membership membershiprepo.Repository
	Session    sessionrepo.Repository
	Token      tokenrepo.Repository
	Flow       flowrepo.Repository
	Device     devicerepo.Repository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:       userrepo.NewRepository(db),
		Tenant:     tenantrepo.NewRepository(db),
		Membership: membershiprepo.NewRepository(db),
		Session:    sessionrepo.NewRepository(db),
		Token:      tokenrepo.NewRepository(db),
		Flow:       flowrepo.NewRepository(db),
		Device:     devicerepo.NewRepository(db),
	}
}
