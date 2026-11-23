package repository

import (
	"context"

	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
