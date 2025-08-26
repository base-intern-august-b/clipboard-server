package repository

import (
	"context"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
)

type ChannelRepository interface {
	CreateChannel(ctx context.Context, req *model.RequestCreateChannel) (*model.Channel, error)
	GetChannelByName(ctx context.Context, channelName string) (*model.Channel, error)
	GetChannels(ctx context.Context) ([]*model.Channel, error)
	PatchChannel(ctx context.Context, channelName string, req *model.RequestPatchChannel) (*model.Channel, error)
	DeleteChannel(ctx context.Context, channelName string) error
}
