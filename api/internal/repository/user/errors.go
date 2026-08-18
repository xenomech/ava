package user

import "ava/api/internal/serrors"

var (
	ErrUserNotFound      = serrors.New("user not found")
	ErrUserAlreadyExists = serrors.New("user already exists")
)
