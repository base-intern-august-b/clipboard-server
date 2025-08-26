package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type messageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(ctx context.Context, req *model.RequestCreateMessage) (*model.Message, error) {
	messageID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	query := `INSERT INTO u_message (message_id, channel_id, user_id, content) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, messageID.String(), req.ChannelID.String(), req.UserID.String(), req.Content)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, model.ErrMessageNotFound
	}

	var createdMessage model.Message
	selectQuery := `SELECT * FROM u_message WHERE message_id = ? LIMIT 1`

	if err := r.db.GetContext(ctx, &createdMessage, selectQuery, messageID.String()); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found after successful insert: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch created message: %w", err)
	}

	return &createdMessage, nil
}

func (r *messageRepository) GetMessages(ctx context.Context, channelID uuid.UUID, limit int, offset int) ([]*model.Message, error) {
	query := `SELECT * FROM u_message WHERE channel_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var messages []*model.Message
	if err := r.db.SelectContext(ctx, &messages, query, channelID.String(), limit, offset); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) GetMessagesInDuration(ctx context.Context, channelID uuid.UUID, start, end time.Time) ([]*model.Message, error) {
	query := `SELECT * FROM u_message WHERE channel_id = ? AND created_at BETWEEN ? AND ? ORDER BY created_at DESC`
	var messages []*model.Message
	if err := r.db.SelectContext(ctx, &messages, query, channelID.String(), start, end); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) GetPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*model.Message, error) {
	query := `SELECT m.* FROM u_message m
	JOIN u_pinned_message pm ON m.message_id = pm.message_id
	WHERE m.channel_id = ? ORDER BY pm.created_at DESC`
	var messages []*model.Message
	if err := r.db.SelectContext(ctx, &messages, query, channelID.String()); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) PatchMessage(ctx context.Context, messageID uuid.UUID, req *model.RequestPatchMessage) (*model.Message, error) {
	setClauses := []string{}
	args := []interface{}{}

	if req.Content != nil {
		setClauses = append(setClauses, "content = ?")
		args = append(args, *req.Content)
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	args = append(args, messageID.String())
	query := fmt.Sprintf("UPDATE u_message SET %s WHERE message_id = ?",
		strings.Join(setClauses, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, model.ErrMessageNotFound
	}

	var updatedMessage model.Message
	selectQuery := `SELECT * FROM u_message WHERE message_id = ? LIMIT 1`

	if err := r.db.GetContext(ctx, &updatedMessage, selectQuery, messageID.String()); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found after successful update: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch updated message: %w", err)
	}

	return &updatedMessage, nil
}

func (r *messageRepository) PinnMessage(ctx context.Context, messageID uuid.UUID) error {
	var message struct {
		ChannelID uuid.UUID `db:"channel_id"`
	}
	// First, get the channel_id from the message
	err := r.db.GetContext(ctx, &message, "SELECT channel_id FROM u_message WHERE message_id = ?", messageID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return model.ErrMessageNotFound
		}
		return err
	}

	// Now, insert into u_pinned_message with both message_id and channel_id
	query := `INSERT INTO u_pinned_message (message_id, channel_id) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, messageID.String(), message.ChannelID.String())
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return model.ErrMessageAlreadyPinned
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// This case should ideally not be reached if the insert succeeds without error
		return fmt.Errorf("failed to pin message, no rows affected")
	}
	return nil
}

func (r *messageRepository) UnpinnMessage(ctx context.Context, messageID uuid.UUID) error {
	query := `DELETE FROM u_pinned_message WHERE message_id = ?`
	result, err := r.db.ExecContext(ctx, query, messageID.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrMessageNotPinned
	}
	return nil
}

func (r *messageRepository) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	query := `DELETE FROM u_message WHERE message_id = ?`
	result, err := r.db.ExecContext(ctx, query, messageID.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrMessageNotFound
	}
	return nil
}
