package repository

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/gofrs/uuid"
)

type ChannelRepository interface {
	CreateChannel(ctx context.Context, req *model.RequestCreateChannel) (*model.Channel, error)
	GetChannel(ctx context.Context, channelID uuid.UUID) (*model.Channel, error)
	GetChannels(ctx context.Context) ([]*model.Channel, error)
	PatchChannel(ctx context.Context, channelID uuid.UUID, req *model.RequestPatchChannel) (*model.Channel, error)
	DeleteChannel(ctx context.Context, channelID uuid.UUID) error
}
