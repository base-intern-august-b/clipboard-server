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

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, req *model.RequestCreateUser) (*model.User, error) {
	query := `INSERT INTO u_user (user_name, nickname, status) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, req.UserName, req.Nickname, req.Status)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, model.ErrAlreadyExistUserName
		}
		return nil, err
	}

	var createdUser model.User
	selectQuery := `SELECT * FROM u_user WHERE user_name = ?`

	if err := r.db.GetContext(ctx, &createdUser, selectQuery, req.UserName); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found after successful insert: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}

	return &createdUser, nil
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	query := `SELECT user_name, nickname FROM u_user`
	var users []*model.User
	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		if err == sql.ErrNoRows {
			return []*model.User{}, nil
		}
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByName(ctx context.Context, userName string) (*model.User, error) {
	query := `SELECT user_name, nickname FROM u_user WHERE user_name = ?`
	var user model.User
	if err := r.db.GetContext(ctx, &user, query, userName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsersByName(ctx context.Context, userNames []string) ([]*model.User, error) {
	if len(userNames) == 0 {
		return []*model.User{}, nil
	}

	placeholders := make([]string, len(userNames))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	placeholderStr := strings.Join(placeholders, ",")

	query := fmt.Sprintf("SELECT user_name, nickname FROM u_user WHERE user_name IN (%s)", placeholderStr)

	args := make([]interface{}, len(userNames))
	for i, id := range userNames {
		args[i] = id
	}

	var users []*model.User
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return []*model.User{}, nil
		}
		return nil, err
	}

	return users, nil
}

func (r *userRepository) PatchUser(ctx context.Context, userName string, req *model.RequestPatchUser) (*model.User, error) {
	setClauses := []string{}
	args := []interface{}{}

	if req.UserName != nil {
		setClauses = append(setClauses, "user_name = ?")
		args = append(args, *req.UserName)
	}
	if req.Nickname != nil {
		setClauses = append(setClauses, "nickname = ?")
		args = append(args, *req.Nickname)
	}
	if req.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *req.Status)
	}

	if len(setClauses) == 0 {
		return r.GetUserByName(ctx, userName)
	}

	args = append(args, userName)
	query := "UPDATE u_user SET " + strings.Join(setClauses, ", ") + " WHERE user_name = ?"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, model.ErrAlreadyExistUserName
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	updatedUserID := userName
	if req.UserName != nil {
		updatedUserID = *req.UserName
	}

	return r.GetUserByName(ctx, updatedUserID)
}

func (r *userRepository) DeleteUser(ctx context.Context, userName string) error {
	query := `DELETE FROM u_user WHERE user_name = ?`
	_, err := r.db.ExecContext(ctx, query, userName)
	return err
}
