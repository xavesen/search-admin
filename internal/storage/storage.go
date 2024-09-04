package storage

import (
	"context"

	"github.com/xavesen/search-admin/internal/models"
)

type Storage interface{
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *models.User) error
	CreateFilter(ctx context.Context, filter *models.Filter) (*models.Filter, error)
	GetAllFilters(ctx context.Context) ([]models.Filter, error)
	DeleteFilter(ctx context.Context, id string) error
}