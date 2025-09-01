package repository

import (
	"context"
	"errors"
	"go-clean-architecture/internal/entity"
	"go-clean-architecture/internal/usecase/interfaces"

	"gorm.io/gorm"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		// Handle duplicate email error
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entity.ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetAll retrieves all users with pagination
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	result := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entity.ErrUserAlreadyExists
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}

// Delete soft deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entity.User{}).Count(&count)
	return count, result.Error
}
