package repository

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *model.RequestCreateUser) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	PatchUser(ctx context.Context, userID uuid.UUID, req *model.RequestPatchUser) (*model.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *model.RequestChangePassword) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
