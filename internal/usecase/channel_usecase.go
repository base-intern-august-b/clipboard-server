package usecase

import (
	"context"
	"regexp"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
)

var (
	channelNameRegex       = `^[a-zA-Z0-9_-]{4,32}$`
	compiledChannelNameReg = regexp.MustCompile(channelNameRegex)
)

type channelUseCase struct {
	channelRepo repository.ChannelRepository
}

func NewChannelUsecase(channelRepo repository.ChannelRepository) usecase.ChannelUsecase {
	return &channelUseCase{
		channelRepo: channelRepo,
	}
}

func (c *channelUseCase) validateChannelName(channelName string) error {
	if channelName == "" {
		return model.ErrInvalidChannelName
	}
	if !compiledChannelNameReg.MatchString(channelName) {
		return model.ErrBadFormatChannelName
	}
	return nil
}

func (c *channelUseCase) CreateChannel(ctx context.Context, req *model.RequestCreateChannel) (*model.Channel, error) {
	if err := c.validateChannelName(req.ChannelName); err != nil {
		return nil, err
	}
	if req.DisplayName == "" {
		return nil, model.ErrInvalidDisplayName
	}
	return c.channelRepo.CreateChannel(ctx, req)
}

func (c *channelUseCase) GetChannelByName(ctx context.Context, channelName string) (*model.Channel, error) {
	if err := c.validateChannelName(channelName); err != nil {
		return nil, err
	}
	return c.channelRepo.GetChannelByName(ctx, channelName)
}

func (c *channelUseCase) GetChannels(ctx context.Context) ([]*model.Channel, error) {
	return c.channelRepo.GetChannels(ctx)
}

func (c *channelUseCase) PatchChannel(ctx context.Context, channelName string, req *model.RequestPatchChannel) (*model.Channel, error) {
	if err := c.validateChannelName(channelName); err != nil {
		return nil, err
	}
	if req.DisplayName != nil && *req.DisplayName == "" {
		return nil, model.ErrInvalidDisplayName
	}
	return c.channelRepo.PatchChannel(ctx, channelName, req)
}

func (c *channelUseCase) DeleteChannel(ctx context.Context, channelName string) error {
	if err := c.validateChannelName(channelName); err != nil {
		return err
	}
	return c.channelRepo.DeleteChannel(ctx, channelName)
}
