package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/gofrs/uuid"
)

type messageUsecase struct {
	messageRepo repository.MessageRepository
	tagRepo     repository.TagRepository
	// userRepo is needed for the composer
	userRepo repository.UserRepository
	messageComposer
}

func NewMessageUsecase(messageRepo repository.MessageRepository, userRepo repository.UserRepository, tagRepo repository.TagRepository) usecase.MessageUsecase {
	return &messageUsecase{
		messageRepo: messageRepo,
		tagRepo:     tagRepo,
		userRepo:    userRepo,
		messageComposer: messageComposer{
			userRepo: userRepo,
			tagRepo:  tagRepo,
		},
	}
}

func (uc *messageUsecase) CreateMessage(ctx context.Context, req *model.RequestCreateMessage) (*model.MessageDetail, error) {
	// Create the base message
	msg, err := uc.messageRepo.CreateMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	// Handle tags
	if len(req.Tags) > 0 {
		tags, hashedTags := uc.createTagsFromRaw(req.Tags)

		if err := uc.tagRepo.CreateTags(ctx, tags); err != nil {
			// Log the error but don't fail the message creation
		}

		if err := uc.tagRepo.CreateTagRelations(ctx, msg.MessageID, msg.ChannelID, hashedTags); err != nil {
			// Log the error but don't fail the message creation
		}
	}

	// Return the full MessageDetail
	details, err := uc.composeMessageDetails(ctx, []*model.Message{msg})
	if err != nil {
		return nil, err // Should not happen if message creation succeeded
	}
	return details[0], nil
}

func (uc *messageUsecase) GetMessage(ctx context.Context, messageID uuid.UUID) (*model.MessageDetail, error) {
	msg, err := uc.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	details, err := uc.composeMessageDetails(ctx, []*model.Message{msg})
	if err != nil {
		return nil, err
	}

	return details[0], nil
}

func (uc *messageUsecase) GetMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]*model.MessageDetail, error) {
	if limit < 1 || limit > 1000 {
		return nil, model.ErrInvalidRequestLimit
	}
	messages, err := uc.messageRepo.GetMessages(ctx, channelID, limit, offset)
	if err != nil {
		return nil, err
	}
	return uc.composeMessageDetails(ctx, messages)
}

func (uc *messageUsecase) GetMessagesInDuration(ctx context.Context, channelID uuid.UUID, start, end time.Time) ([]*model.MessageDetail, error) {
	if start.After(end) {
		return nil, model.ErrInvalidTimeRange
	}
	messages, err := uc.messageRepo.GetMessagesInDuration(ctx, channelID, start, end)
	if err != nil {
		return nil, err
	}
	return uc.composeMessageDetails(ctx, messages)
}

func (uc *messageUsecase) GetPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*model.MessageDetail, error) {
	messages, err := uc.messageRepo.GetPinnedMessages(ctx, channelID)
	if err != nil {
		return nil, err
	}
	return uc.composeMessageDetails(ctx, messages)
}

func (uc *messageUsecase) PatchMessage(ctx context.Context, messageID uuid.UUID, req *model.RequestPatchMessage) (*model.MessageDetail, error) {
	// Note: This does not handle changing tags. A separate endpoint is used for that.
	msg, err := uc.messageRepo.PatchMessage(ctx, messageID, req)
	if err != nil {
		return nil, err
	}

	details, err := uc.composeMessageDetails(ctx, []*model.Message{msg})
	if err != nil {
		return nil, err
	}
	return details[0], nil
}

func (uc *messageUsecase) PinnMessage(ctx context.Context, messageID uuid.UUID) error {
	return uc.messageRepo.PinnMessage(ctx, messageID)
}

func (uc *messageUsecase) UnpinnMessage(ctx context.Context, messageID uuid.UUID) error {
	return uc.messageRepo.UnpinnMessage(ctx, messageID)
}

func (uc *messageUsecase) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	return uc.messageRepo.DeleteMessage(ctx, messageID)
}

func (uc *messageUsecase) ToMessageDetails(ctx context.Context, messages []*model.Message) ([]*model.MessageDetail, error) {
	return uc.composeMessageDetails(ctx, messages)
}

// --- Helper functions for tags ---

func (uc *messageUsecase) getHashedTag(rawTag string) string {
	hash := md5.Sum([]byte(rawTag))
	return hex.EncodeToString(hash[:])
}

func (uc *messageUsecase) createTagsFromRaw(rawTags []string) ([]*model.Tag, []string) {
	tags := make([]*model.Tag, len(rawTags))
	hashedTags := make([]string, len(rawTags))
	for i, rawTag := range rawTags {
		hashedTag := uc.getHashedTag(rawTag)
		tags[i] = &model.Tag{
			HashedTag: hashedTag,
			RawTag:    rawTag,
		}
		hashedTags[i] = hashedTag
	}
	return tags, hashedTags
}
