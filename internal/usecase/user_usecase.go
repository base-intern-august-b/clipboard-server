package usecase

import (
	"context"
	"regexp"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
)

var (
	RegexpUserName      = `^[a-zA-Z0-9][a-zA-Z0-9_-]{2,30}[a-zA-Z0-9]$`
	compiledUserNameReg = regexp.MustCompile(RegexpUserName)
)

// userUseCase はユーザーに関するユースケース
type userUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase は新しいUserUsecaseを作成する
func NewUserUsecase(userRepo repository.UserRepository) usecase.UserUsecase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) validateUserID(userName string) error {
	if userName == "" {
		return model.ErrInvalidUserName
	}
	if !compiledUserNameReg.MatchString(userName) {
		return model.ErrBadFormatUserName
	}
	return nil
}

func (u *userUseCase) CreateUser(ctx context.Context, req *model.RequestCreateUser) (*model.User, error) {
	if err := u.validateUserID(req.UserName); err != nil {
		return nil, err
	}
	if req.Nickname == "" {
		return nil, model.ErrInvalidNickname
	}
	return u.userRepo.CreateUser(ctx, req)
}

func (u *userUseCase) GetUsers(ctx context.Context) ([]*model.User, error) {
	return u.userRepo.GetUsers(ctx)
}

func (u *userUseCase) GetUserByName(ctx context.Context, userID string) (*model.User, error) {
	if err := u.validateUserID(userID); err != nil {
		return nil, err
	}
	return u.userRepo.GetUserByName(ctx, userID)
}

func (u *userUseCase) GetUsersByName(ctx context.Context, req []string) ([]*model.User, error) {
	var userIDs []string
	for _, userID := range req {
		if err := u.validateUserID(userID); err != nil {
			continue
		}
		userIDs = append(userIDs, userID)
	}
	return u.userRepo.GetUsersByName(ctx, userIDs)
}

func (u *userUseCase) PatchUser(ctx context.Context, userID string, req *model.RequestPatchUser) (*model.User, error) {
	if err := u.validateUserID(userID); err != nil {
		return nil, err
	}
	return u.userRepo.PatchUser(ctx, userID, req)
}

func (u *userUseCase) DeleteUser(ctx context.Context, userID string) error {
	if err := u.validateUserID(userID); err != nil {
		return err
	}
	return u.userRepo.DeleteUser(ctx, userID)
}
