package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-todolist-sber/internal/apperror"
	pgxError "go-todolist-sber/internal/apperror/pgx_errors"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/user"
	"go-todolist-sber/pkg/postgres"
)

type userRepository struct {
	*postgres.Postgres
}

func NewUserRepository(postgres *postgres.Postgres) user.UserRepository {
	return &userRepository{
		postgres,
	}
}

func (u *userRepository) collectRow(row pgx.Row) (*entity.User, error) {
	var user entity.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, apperror.ErrNoRows
	}
	errCode := pgxError.ErrorCode(err)
	if errCode == pgxError.ForeignKeyViolation {
		return nil, apperror.ErrForeignKeyViolation
	}
	if errCode == pgxError.UniqueViolation {
		return nil, apperror.ErrUniqueViolation
	}

	return &user, err
}

func (u *userRepository) collectRows(rows pgx.Rows) ([]entity.User, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.User, error) {
		user, err := u.collectRow(row)
		return *user, err
	})
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `insert into "user" (login,password) values ($1,$2) returning *`

	row := u.Pool.QueryRow(ctx, query, user.Login, user.Password)
	return u.collectRow(row)
}

func (u *userRepository) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	query := `select id,login,password from "user" where login = $1`

	row := u.Pool.QueryRow(ctx, query, login)
	return u.collectRow(row)
}
