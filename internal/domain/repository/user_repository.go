package repository

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *model.RequestCreateUser) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUserByName(ctx context.Context, userName string) (*model.User, error)
	GetUsersByName(ctx context.Context, userNames []string) ([]*model.User, error)
	PatchUser(ctx context.Context, userName string, req *model.RequestPatchUser) (*model.User, error)
	DeleteUser(ctx context.Context, userName string) error
}
