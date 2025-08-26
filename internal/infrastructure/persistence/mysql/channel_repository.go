package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type channelRepository struct {
	db *sqlx.DB
}

func NewChannelRepository(db *sqlx.DB) repository.ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) CreateChannel(ctx context.Context, req *model.RequestCreateChannel) (*model.Channel, error) {
	query := `INSERT INTO u_channel (channel_name, display_name, description) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, req.ChannelName, req.DisplayName, req.Description)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, model.ErrAlreadyExistChannelName
		}
		return nil, err
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

func (r *channelRepository) GetChannelByName(ctx context.Context, channelName string) (*model.Channel, error) {
	query := `SELECT * FROM u_channel WHERE channel_name = ?`
	var channel model.Channel
	if err := r.db.GetContext(ctx, &channel, query, channelName); err != nil {
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

func (r *channelRepository) PatchChannel(ctx context.Context, channelName string, req *model.RequestPatchChannel) (*model.Channel, error) {
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
		return r.GetChannelByName(ctx, channelName)
	}

	args = append(args, channelName)
	query := fmt.Sprintf("UPDATE u_channel SET %s WHERE channel_name = ?", strings.Join(setClauses, ", "))

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

	return r.GetChannelByName(ctx, channelName)
}

func (r *channelRepository) DeleteChannel(ctx context.Context, channelName string) error {
	query := `DELETE FROM u_channel WHERE channel_name = ?`
	result, err := r.db.ExecContext(ctx, query, channelName)
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
