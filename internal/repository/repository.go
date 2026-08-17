package repository

import "gorm.io/gorm"

type Repository struct {
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{}
}
