package user

import (
	"context"
	"errors"
	"time"

	"ava/internal/model"
	"ava/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error

	logger.Info("userRepository.CreateUser", logger.Any("user.ID", user.ID), logger.Any("error", err))

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrUserAlreadyExists
	}

	return err
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var users []*model.User

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(user).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"password":   hashedPassword,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()

	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"email_verified":    true,
		"email_verified_at": &now,
		"updated_at":        now,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
