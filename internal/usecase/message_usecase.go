package usecase

import (
	"context"
	"time"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/gofrs/uuid"
)

type messageUsecase struct {
	messageRepo repository.MessageRepository
}

func NewMessageUsecase(messageRepo repository.MessageRepository) usecase.MessageUsecase {
	return &messageUsecase{
		messageRepo: messageRepo,
	}
}

func (m *messageUsecase) CreateMessage(ctx context.Context, req *model.RequestCreateMessage) (*model.Message, error) {
	return m.messageRepo.CreateMessage(ctx, req)
}

func (m *messageUsecase) GetMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]*model.Message, error) {
	if limit < 1 || limit > 1000 {
		return nil, model.ErrInvalidRequestLimit
	}
	return m.messageRepo.GetMessages(ctx, channelID, limit, offset)
}

func (m *messageUsecase) GetMessagesInDuration(ctx context.Context, channelID uuid.UUID, start, end time.Time) ([]*model.Message, error) {
	if start.After(end) {
		return nil, model.ErrInvalidTimeRange
	}
	return m.messageRepo.GetMessagesInDuration(ctx, channelID, start, end)
}

func (m *messageUsecase) GetPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*model.Message, error) {
	return m.messageRepo.GetPinnedMessages(ctx, channelID)
}

func (m *messageUsecase) PatchMessage(ctx context.Context, messageID uuid.UUID, req *model.RequestPatchMessage) (*model.Message, error) {
	return m.messageRepo.PatchMessage(ctx, messageID, req)
}

func (m *messageUsecase) PinnMessage(ctx context.Context, messageID uuid.UUID) error {
	return m.messageRepo.PinnMessage(ctx, messageID)
}

func (m *messageUsecase) UnpinnMessage(ctx context.Context, messageID uuid.UUID) error {
	return m.messageRepo.UnpinnMessage(ctx, messageID)
}

func (m *messageUsecase) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	return m.messageRepo.DeleteMessage(ctx, messageID)
}
