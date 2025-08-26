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
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (r *userRepository) CreateUser(ctx context.Context, req *model.RequestCreateUser) (*model.User, error) {
	userID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	userQuery := `INSERT INTO u_user (user_id, user_name, nickname, status) VALUES (?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, userQuery, userID.String(), req.UserName, req.Nickname, req.Status)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, model.ErrAlreadyExistUserName
		}
		return nil, fmt.Errorf("failed to insert into u_user: %w", err)
	}

	privateQuery := `INSERT INTO u_user_private (user_id, password_hash) VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, privateQuery, userID.String(), string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to insert into u_user_private: %w", err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var createdUser model.User
	selectQuery := `SELECT * FROM u_user WHERE user_id = ?`

	if err := r.db.GetContext(ctx, &createdUser, selectQuery, userID.String()); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found after successful insert: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}

	return &createdUser, nil
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	query := `SELECT * FROM u_user`
	var users []*model.User
	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		if err == sql.ErrNoRows {
			return []*model.User{}, nil
		}
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	query := `SELECT * FROM u_user WHERE user_id = ?`
	var user model.User
	if err := r.db.GetContext(ctx, &user, query, userID.String()); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) PatchUser(ctx context.Context, userID uuid.UUID, req *model.RequestPatchUser) (*model.User, error) {
	setClauses := []string{}
	args := []interface{}{}

	if req.UserName != nil {
		setClauses = append(setClauses, "user_name = ?")
		args = append(args, *req.UserName)
	}
	if req.Email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, *req.Email)
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
		return r.GetUserByID(ctx, userID)
	}

	args = append(args, userID.String())
	query := "UPDATE u_user SET " + strings.Join(setClauses, ", ") + " WHERE user_id = ?"

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, model.ErrAlreadyExistUserName
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, model.ErrUserNotFound
	}

	return r.GetUserByID(ctx, userID)
}

func (r *userRepository) ChangePassword(ctx context.Context, userID uuid.UUID, req *model.RequestChangePassword) error {
	var storedHash string
	query := `SELECT password_hash FROM u_user_private WHERE user_id = ?`
	err := r.db.GetContext(ctx, &storedHash, query, userID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return model.ErrUserNotFound
		}
		return fmt.Errorf("failed to fetch user private data: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("old password does not match: %w", err)
	}

	// Hash new password
	newHashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE u_user_private SET password_hash = ? WHERE user_id = ?`
	result, err := r.db.ExecContext(ctx, updateQuery, newHashedPassword, userID.String())
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM u_user WHERE user_id = ?`
	result, err := r.db.ExecContext(ctx, query, userID.String())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return err
}
