package usecase

import (
	"context"
	"time"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type MessageUsecase interface {
	CreateMessage(ctx context.Context, req *model.RequestCreateMessage) (*model.Message, error)
	GetMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]*model.Message, error)
	GetMessagesInDuration(ctx context.Context, channelID uuid.UUID, start, end time.Time) ([]*model.Message, error)
	GetPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*model.Message, error)
	PatchMessage(ctx context.Context, messageID uuid.UUID, req *model.RequestPatchMessage) (*model.Message, error)
	PinnMessage(ctx context.Context, messageID uuid.UUID) error
	UnpinnMessage(ctx context.Context, messageID uuid.UUID) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
}
