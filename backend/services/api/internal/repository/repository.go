package repository

import (
	devicerepo "ava/api/internal/repository/device"
	flowrepo "ava/api/internal/repository/flow"
	hubrepo "ava/api/internal/repository/hub"
	membershiprepo "ava/api/internal/repository/membership"
	roomrepo "ava/api/internal/repository/room"
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
	Room       roomrepo.Repository
	Hub        hubrepo.Repository
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
		Room:       roomrepo.NewRepository(db),
		Hub:        hubrepo.NewRepository(db),
		Device:     devicerepo.NewRepository(db),
	}
}
