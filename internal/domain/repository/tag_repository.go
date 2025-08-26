package repository

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type TagRepository interface {
	// CreateTags はタグを複数作成します。既に存在する場合は無視されます。
	CreateTags(ctx context.Context, tags []*model.Tag) error
	// CreateTagRelations はメッセージとタグの関連を複数作成します。
	CreateTagRelations(ctx context.Context, messageID uuid.UUID, channelID uuid.UUID, hashedTags []string) error
	// GetTagsByMessageIDs は複数のメッセージIDに紐づくタグを取得します。
	GetTagsByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID][]*model.Tag, error)
	// FindMessagesByHashedTag はハッシュ化されたタグでメッセージを検索します。
	FindMessagesByHashedTag(ctx context.Context, hashedTag string, limit int, offset int) ([]*model.Message, error)
	// FindMessagesByHashedTagInChannel はチャンネル内でハッシュ化されたタグでメッセージを検索します。
	FindMessagesByHashedTagInChannel(ctx context.Context, channelID uuid.UUID, hashedTag string, limit int, offset int) ([]*model.Message, error)
	// DeleteTagRelationsByMessageID はメッセージに紐づくタグの関連を全て削除します。
	DeleteTagRelationsByMessageID(ctx context.Context, messageID uuid.UUID) error
}
