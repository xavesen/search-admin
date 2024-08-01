package storage

import (
	"context"

	"github.com/xavesen/search-admin/internal/models"
)

type Storage interface{
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}