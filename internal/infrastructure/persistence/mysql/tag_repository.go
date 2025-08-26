package mysql

import (
	"context"
	"strings"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type tagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) CreateTags(ctx context.Context, tags []*model.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	query := `
		INSERT INTO u_tags (hashed_tag, raw_tag)
		VALUES (:hashed_tag, :raw_tag)
		ON DUPLICATE KEY UPDATE raw_tag = VALUES(raw_tag)
	`
	_, err := r.db.NamedExecContext(ctx, query, tags)
	return err
}

func (r *tagRepository) CreateTagRelations(ctx context.Context, messageID uuid.UUID, channelID uuid.UUID, hashedTags []string) error {
	if len(hashedTags) == 0 {
		return nil
	}

	query := `
		INSERT INTO u_tag_relation (message_id, channel_id, hashed_tag)
		VALUES ` + strings.TrimSuffix(strings.Repeat("(?, ?, ?),", len(hashedTags)), ",")

	args := make([]interface{}, 0, len(hashedTags)*3)
	for _, hashedTag := range hashedTags {
		args = append(args, messageID, channelID, hashedTag)
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *tagRepository) GetTagsByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID][]*model.Tag, error) {
	if len(messageIDs) == 0 {
		return make(map[uuid.UUID][]*model.Tag), nil
	}

	query := `
		SELECT
			r.message_id,
			t.hashed_tag,
			t.raw_tag,
			t.created_at
		FROM
			u_tags t
		JOIN
			u_tag_relation r ON t.hashed_tag = r.hashed_tag
		WHERE
			r.message_id IN (?)
	`
	query, args, err := sqlx.In(query, messageIDs)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	type result struct {
		MessageID uuid.UUID `db:"message_id"`
		model.Tag
	}

	var results []result
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, err
	}

	tagsMap := make(map[uuid.UUID][]*model.Tag)
	for _, res := range results {
		tag := res.Tag
		tagsMap[res.MessageID] = append(tagsMap[res.MessageID], &tag)
	}
    // Ensure messages with no tags are present in the map
    for _, id := range messageIDs {
        if _, ok := tagsMap[id]; !ok {
            tagsMap[id] = []*model.Tag{}
        }
    }

	return tagsMap, nil
}

func (r *tagRepository) FindMessagesByHashedTag(ctx context.Context, hashedTag string, limit int, offset int) ([]*model.Message, error) {
	query := `
		SELECT
			m.*
		FROM
			u_message m
		JOIN
			u_tag_relation r ON m.message_id = r.message_id
		WHERE
			r.hashed_tag = ?
		ORDER BY
			m.created_at DESC
		LIMIT ? OFFSET ?
	`
	var messages []*model.Message
	err := r.db.SelectContext(ctx, &messages, query, hashedTag, limit, offset)
	return messages, err
}

func (r *tagRepository) FindMessagesByHashedTagInChannel(ctx context.Context, channelID uuid.UUID, hashedTag string, limit int, offset int) ([]*model.Message, error) {
	query := `
		SELECT
			m.*
		FROM
			u_message m
		JOIN
			u_tag_relation r ON m.message_id = r.message_id
		WHERE
			r.channel_id = ? AND r.hashed_tag = ?
		ORDER BY
			m.created_at DESC
		LIMIT ? OFFSET ?
	`
	var messages []*model.Message
	err := r.db.SelectContext(ctx, &messages, query, channelID, hashedTag, limit, offset)
	return messages, err
}

func (r *tagRepository) DeleteTagRelationsByMessageID(ctx context.Context, messageID uuid.UUID) error {
	query := "DELETE FROM u_tag_relation WHERE message_id = ?"
	_, err := r.db.ExecContext(ctx, query, messageID)
	return err
}
