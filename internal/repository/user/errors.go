package user

import "ava/pkg/serrors"

var (
	ErrUserNotFound      = serrors.New("user not found")
	ErrUserAlreadyExists = serrors.New("user already exists")
)
