package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type channelRepository struct {
	db *sqlx.DB
}

func NewChannelRepository(db *sqlx.DB) repository.ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) CreateChannel(ctx context.Context, req *model.RequestCreateChannel) (*model.Channel, error) {
	channelID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	query := `INSERT INTO u_channel (channel_id, channel_name, display_name, description) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, channelID.String(), req.ChannelName, req.DisplayName, req.Description)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, model.ErrAlreadyExistChannelName
		}
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, model.ErrChannelNotFound
	}

	var createdChannel model.Channel
	selectQuery := `SELECT * FROM u_channel WHERE channel_name = ?`

	if err := r.db.GetContext(ctx, &createdChannel, selectQuery, req.ChannelName); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("channel not found after successful insert: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch created channel: %w", err)
	}

	return &createdChannel, nil
}

func (r *channelRepository) GetChannel(ctx context.Context, channelID uuid.UUID) (*model.Channel, error) {
	query := `SELECT * FROM u_channel WHERE channel_id = ?`
	var channel model.Channel
	if err := r.db.GetContext(ctx, &channel, query, channelID.String()); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &channel, nil
}

func (r *channelRepository) GetChannels(ctx context.Context) ([]*model.Channel, error) {
	query := `SELECT * FROM u_channel`
	var channels []*model.Channel
	if err := r.db.SelectContext(ctx, &channels, query); err != nil {
		if err == sql.ErrNoRows {
			return []*model.Channel{}, nil
		}
		return nil, err
	}
	return channels, nil
}

func (r *channelRepository) PatchChannel(ctx context.Context, channelID uuid.UUID, req *model.RequestPatchChannel) (*model.Channel, error) {
	setClauses := []string{}
	args := []interface{}{}

	if req.DisplayName != nil {
		setClauses = append(setClauses, "display_name = ?")
		args = append(args, *req.DisplayName)
	}
	if req.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *req.Description)
	}

	if len(setClauses) == 0 {
		return r.GetChannel(ctx, channelID)
	}

	args = append(args, channelID.String())
	query := fmt.Sprintf("UPDATE u_channel SET %s WHERE channel_id = ?", strings.Join(setClauses, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, model.ErrChannelNotFound
	}

	return r.GetChannel(ctx, channelID)
}

func (r *channelRepository) DeleteChannel(ctx context.Context, channelID uuid.UUID) error {
	query := `DELETE FROM u_channel WHERE channel_id = ?`
	result, err := r.db.ExecContext(ctx, query, channelID.String())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrChannelNotFound
	}
	return nil
}
