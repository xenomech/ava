package user

import (
	"context"

	"ava/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}
