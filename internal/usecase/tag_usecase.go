package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/gofrs/uuid"
)

type tagUsecase struct {
	tagRepo     repository.TagRepository
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
	messageComposer
}

func NewTagUsecase(tagRepo repository.TagRepository, messageRepo repository.MessageRepository, userRepo repository.UserRepository) usecase.TagUsecase {
	return &tagUsecase{
		tagRepo:     tagRepo,
		messageRepo: messageRepo,
		userRepo:    userRepo,
		messageComposer: messageComposer{
			userRepo: userRepo,
			tagRepo:  tagRepo,
		},
	}
}

func (uc *tagUsecase) ModifyMessageTags(ctx context.Context, messageID uuid.UUID, rawTags []string) (*model.MessageDetail, error) {
	// Get the message to ensure it exists and to get the channelID
	msg, err := uc.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// Delete all existing tag relations for the message
	if err := uc.tagRepo.DeleteTagRelationsByMessageID(ctx, messageID); err != nil {
		return nil, err
	}

	// Create the new tags and relations
	if len(rawTags) > 0 {
		tags, hashedTags := uc.createTagsFromRaw(rawTags)

		if err := uc.tagRepo.CreateTags(ctx, tags); err != nil {
			return nil, err
		}

		if err := uc.tagRepo.CreateTagRelations(ctx, msg.MessageID, msg.ChannelID, hashedTags); err != nil {
			return nil, err
		}
	}

	// Fetch the updated message and compose its details
	updatedMsg, err := uc.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}
	details, err := uc.composeMessageDetails(ctx, []*model.Message{updatedMsg})
	if err != nil {
		return nil, err
	}
	return details[0], nil
}

func (uc *tagUsecase) FindMessagesByTag(ctx context.Context, rawTag string, limit, offset int) ([]*model.MessageDetail, error) {
	hashedTag := uc.getHashedTag(rawTag)
	messages, err := uc.tagRepo.FindMessagesByHashedTag(ctx, hashedTag, limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.composeMessageDetails(ctx, messages)
}

func (uc *tagUsecase) FindMessagesByTagInChannel(ctx context.Context, channelID uuid.UUID, rawTag string, limit, offset int) ([]*model.MessageDetail, error) {
	hashedTag := uc.getHashedTag(rawTag)
	messages, err := uc.tagRepo.FindMessagesByHashedTagInChannel(ctx, channelID, hashedTag, limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.composeMessageDetails(ctx, messages)
}

func (uc *tagUsecase) getHashedTag(rawTag string) string {
	hash := md5.Sum([]byte(rawTag))
	return hex.EncodeToString(hash[:])
}

func (uc *tagUsecase) createTagsFromRaw(rawTags []string) ([]*model.Tag, []string) {
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
