package storage

import (
	"context"

	"github.com/xavesen/search-admin/internal/models"
)

type StorageMock struct {
	Error	error
	Users	[]models.User
}

func (s *StorageMock) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (s *StorageMock) GetAllUsers(ctx context.Context) ([]models.User, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	return s.Users, nil
}

func (s *StorageMock) GetUser(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

func (s *StorageMock) DeleteUser(ctx context.Context, id string) error {
	return nil
}