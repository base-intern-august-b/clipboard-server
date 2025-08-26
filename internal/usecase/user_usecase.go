package usecase

import (
	"context"
	"regexp"
	"unicode"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/model"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/repository"
	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/gofrs/uuid"
)

var (
	RegexpUserName      = `^[a-zA-Z0-9][a-zA-Z0-9_-]{2,30}[a-zA-Z0-9]$`
	RegexpEmail         = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	compiledUserNameReg = regexp.MustCompile(RegexpUserName)
	compiledEmailReg    = regexp.MustCompile(RegexpEmail)
)

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) usecase.UserUsecase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func (u *userUseCase) validateUserName(userName string) error {
	if userName == "" {
		return model.ErrInvalidUserName
	}
	if !compiledUserNameReg.MatchString(userName) {
		return model.ErrBadFormatUserName
	}
	return nil
}

func (u *userUseCase) CreateUser(ctx context.Context, req *model.RequestCreateUser) (*model.User, error) {
	if err := u.validateUserName(req.UserName); err != nil {
		return nil, err
	}
	if req.Nickname == "" {
		return nil, model.ErrInvalidNickname
	}
	if !ValidatePassword(req.Password) {
		return nil, model.ErrWeakPassword
	}
	return u.userRepo.CreateUser(ctx, req)
}

func (u *userUseCase) GetUsers(ctx context.Context) ([]*model.User, error) {
	return u.userRepo.GetUsers(ctx)
}

func (u *userUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return u.userRepo.GetUserByID(ctx, userID)
}

func (u *userUseCase) PatchUser(ctx context.Context, userID uuid.UUID, req *model.RequestPatchUser) (*model.User, error) {
	if req.UserName != nil {
		if err := u.validateUserName(*req.UserName); err != nil {
			return nil, err
		}
	}
	if req.Nickname != nil && *req.Nickname == "" {
		return nil, model.ErrInvalidUserName
	}
	if req.Email != nil {
		if !compiledEmailReg.MatchString(*req.Email) {
			return nil, model.ErrBadFormatEmail
		}
	}
	return u.userRepo.PatchUser(ctx, userID, req)
}

func (u *userUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *model.RequestChangePassword) error {
	if req.OldPassword == req.NewPassword {
		return model.ErrNothingChanged
	}
	if !ValidatePassword(req.NewPassword) {
		return model.ErrWeakPassword
	}
	return u.userRepo.ChangePassword(ctx, userID, req)
}

func (u *userUseCase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return u.userRepo.DeleteUser(ctx, userID)
}
