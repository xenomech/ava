package repository

import (
	apitokenrepo "ava/api/internal/repository/apitoken"
	devicerepo "ava/api/internal/repository/device"
	flowrepo "ava/api/internal/repository/flow"
	hubrepo "ava/api/internal/repository/hub"
	membershiprepo "ava/api/internal/repository/membership"
	roomrepo "ava/api/internal/repository/room"
	scenerepo "ava/api/internal/repository/scene"
	sessionrepo "ava/api/internal/repository/session"
	tenantrepo "ava/api/internal/repository/tenant"
	tokenrepo "ava/api/internal/repository/token"
	userrepo "ava/api/internal/repository/user"

	"gorm.io/gorm"
)

type Repository struct {
	APIToken   apitokenrepo.Repository
	User       userrepo.Repository
	Tenant     tenantrepo.Repository
	Membership membershiprepo.Repository
	Session    sessionrepo.Repository
	Token      tokenrepo.Repository
	Flow       flowrepo.Repository
	Room       roomrepo.Repository
	Scene      scenerepo.Repository
	Hub        hubrepo.Repository
	Device     devicerepo.Repository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		APIToken:   apitokenrepo.NewRepository(db),
		User:       userrepo.NewRepository(db),
		Tenant:     tenantrepo.NewRepository(db),
		Membership: membershiprepo.NewRepository(db),
		Session:    sessionrepo.NewRepository(db),
		Token:      tokenrepo.NewRepository(db),
		Flow:       flowrepo.NewRepository(db),
		Room:       roomrepo.NewRepository(db),
		Scene:      scenerepo.NewRepository(db),
		Hub:        hubrepo.NewRepository(db),
		Device:     devicerepo.NewRepository(db),
	}
}
