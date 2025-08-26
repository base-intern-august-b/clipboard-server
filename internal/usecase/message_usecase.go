package usecase

import (
	"context"

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

func (m *messageUsecase) GetMessages(ctx context.Context, req *model.RequestGetMessages) ([]*model.Message, error) {
	if req.Limit < 1 || req.Limit > 1000 {
		return nil, model.ErrInvalidRequestLimit
	}
	return m.messageRepo.GetMessages(ctx, req)
}

func (m *messageUsecase) GetMessagesInDuration(ctx context.Context, req *model.RequestGetMessagesInDuration) ([]*model.Message, error) {
	if req.StartTime.After(req.EndTime) {
		return nil, model.ErrInvalidTimeRange
	}
	return m.messageRepo.GetMessagesInDuration(ctx, req)
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
