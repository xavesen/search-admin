package storage

import (
	"context"

	"github.com/xavesen/search-admin/internal/models"
)

type StorageMock struct {
	Error	error
	Users	[]models.User
	User 	models.User
	Filters	[]models.Filter
}

func (s *StorageMock) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	user.Id = "1"

	return user, nil
}

func (s *StorageMock) GetAllUsers(ctx context.Context) ([]models.User, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	return s.Users, nil
}

func (s *StorageMock) GetUser(ctx context.Context, id string) (*models.User, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	return &s.User, nil
}

func (s *StorageMock) DeleteUser(ctx context.Context, id string) error {
	return s.Error
}

func (s *StorageMock) UpdateUser(ctx context.Context, user *models.User) error{
	return s.Error
}

func (s *StorageMock) CreateFilter(ctx context.Context, filter *models.Filter) (*models.Filter, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	filter.Id = "1"

	return filter, nil
}

func (s *StorageMock) GetAllFilters(ctx context.Context) ([]models.Filter, error) {
	if s.Error != nil {
		return nil, s.Error
	}

	return s.Filters, nil
}

func (s *StorageMock) DeleteFilter(ctx context.Context, id string) error {
	return s.Error
}