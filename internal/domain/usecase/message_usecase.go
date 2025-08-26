package usecase

import (
	"context"
	"time"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type MessageUsecase interface {
	CreateMessage(ctx context.Context, req *model.RequestCreateMessage) (*model.MessageDetail, error)
	GetMessage(ctx context.Context, messageID uuid.UUID) (*model.MessageDetail, error)
	GetMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]*model.MessageDetail, error)
	GetMessagesInDuration(ctx context.Context, channelID uuid.UUID, start, end time.Time) ([]*model.MessageDetail, error)
	GetPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*model.MessageDetail, error)
	PatchMessage(ctx context.Context, messageID uuid.UUID, req *model.RequestPatchMessage) (*model.MessageDetail, error)
	PinnMessage(ctx context.Context, messageID uuid.UUID) error
	UnpinnMessage(ctx context.Context, messageID uuid.UUID) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error
	ToMessageDetails(ctx context.Context, messages []*model.Message) ([]*model.MessageDetail, error)
}
