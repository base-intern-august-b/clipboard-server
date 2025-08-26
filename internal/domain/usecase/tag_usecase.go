package usecase

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type TagUsecase interface {
	// ModifyMessageTags はメッセージについているタグを編集します
	ModifyMessageTags(ctx context.Context, messageID uuid.UUID, rawTags []string) (*model.MessageDetail, error)
	// FindMessagesByTag はタグでメッセージを検索します
	FindMessagesByTag(ctx context.Context, rawTag string, limit, offset int) ([]*model.MessageDetail, error)
	// FindMessagesByTagInChannel はチャンネル内でタグでメッセージを検索します
	FindMessagesByTagInChannel(ctx context.Context, channelID uuid.UUID, rawTag string, limit, offset int) ([]*model.MessageDetail, error)
}
