package usecase

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type MessageUsecase interface {
	CreateMessage(ctx context.Context, req *model.RequestCreateMessage) (*model.Message, error)
	GetMessages(ctx context.Context, req *model.RequestGetMessages) ([]*model.Message, error)
	GetMessagesInDuration(ctx context.Context, req *model.RequestGetMessagesInDuration) ([]*model.Message, error)
	GetPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*model.Message, error)
	PatchMessage(ctx context.Context, messageID uuid.UUID, req *model.RequestPatchMessage) (*model.Message, error)
	PinnMessage(ctx context.Context, messageID uuid.UUID) error
	UnpinnMessage(ctx context.Context, messageID uuid.UUID) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
}
