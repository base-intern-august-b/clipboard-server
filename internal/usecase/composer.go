package usecase

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/gofrs/uuid"
)

// messageComposer is a helper for composing MessageDetail objects.
// It is not a full usecase, but a component to be embedded in other usecases.
type messageComposer struct {
	userRepo repository.UserRepository
	tagRepo  repository.TagRepository
}

func (c *messageComposer) composeMessageDetails(ctx context.Context, messages []*model.Message) ([]*model.MessageDetail, error) {
	if len(messages) == 0 {
		return []*model.MessageDetail{}, nil
	}

	// Get user IDs and message IDs from messages
	userIDs := make([]uuid.UUID, len(messages))
	messageIDs := make([]uuid.UUID, len(messages))
	for i, msg := range messages {
		userIDs[i] = msg.UserID
		messageIDs[i] = msg.MessageID
	}

	// Fetch all users and tags in batch
	users, err := c.userRepo.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	tagsMap, err := c.tagRepo.GetTagsByMessageIDs(ctx, messageIDs)
	if err != nil {
		return nil, err
	}

	// Create a map of user ID to user for easy lookup
	userMap := make(map[uuid.UUID]*model.User, len(users))
	for _, u := range users {
		userMap[u.UserID] = u
	}

	// Compose MessageDetail objects
	details := make([]*model.MessageDetail, len(messages))
	for i, msg := range messages {
		details[i] = &model.MessageDetail{
			Message: *msg,
			User:    userMap[msg.UserID],
			Tags:    tagsMap[msg.MessageID],
		}
	}

	return details, nil
}
