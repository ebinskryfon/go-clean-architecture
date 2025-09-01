package usecase

import (
	"context"
	"fmt"
	"go-clean-architecture/internal/entity"
	"go-clean-architecture/internal/usecase/interfaces"
)

// UserUseCase implements business logic for user operations
type UserUseCase struct {
	userRepo interfaces.UserRepository
}

// NewUserUseCase creates a new user use case instance
func NewUserUseCase(userRepo interfaces.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with business validation
func (uc *UserUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	// Business validation
	if !user.IsValid() {
		return entity.ErrInvalidUserName
	}

	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return entity.ErrUserAlreadyExists
	}

	// Create user
	return uc.userRepo.Create(ctx, user)
}

// GetUser retrieves a user by ID
func (uc *UserUseCase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	if id == 0 {
		return nil, entity.ErrInvalidUserID
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, entity.ErrInvalidUserEmail
	}

	return uc.userRepo.GetByEmail(ctx, email)
}

// GetAllUsers retrieves all users with pagination
func (uc *UserUseCase) GetAllUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, err := uc.userRepo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser updates an existing user
func (uc *UserUseCase) UpdateUser(ctx context.Context, id uint, user *entity.User) error {
	if id == 0 {
		return entity.ErrInvalidUserID
	}

	// Check if user exists
	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return entity.ErrUserNotFound
	}

	// Business validation
	if !user.IsValid() {
		return entity.ErrInvalidUserName
	}

	// Check email uniqueness if email is being changed
	if user.Email != existingUser.Email {
		emailUser, err := uc.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && emailUser != nil && emailUser.ID != id {
			return fmt.Errorf("email already taken by another user")
		}
	}

	user.ID = id
	return uc.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user by ID
func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) error {
	if id == 0 {
		return entity.ErrInvalidUserID
	}

	// Check if user exists
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	return uc.userRepo.Delete(ctx, id)
}

// ActivateUser activates a user
func (uc *UserUseCase) ActivateUser(ctx context.Context, id uint) error {
	user, err := uc.GetUser(ctx, id)
	if err != nil {
		return err
	}

	user.Activate()
	return uc.userRepo.Update(ctx, user)
}

// DeactivateUser deactivates a user
func (uc *UserUseCase) DeactivateUser(ctx context.Context, id uint) error {
	user, err := uc.GetUser(ctx, id)
	if err != nil {
		return err
	}

	user.Deactivate()
	return uc.userRepo.Update(ctx, user)
}
