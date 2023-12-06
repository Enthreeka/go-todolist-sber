package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-todolist-sber/internal/apperror"
	pgxError "go-todolist-sber/internal/apperror/pgx_errors"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/session"
	"go-todolist-sber/pkg/postgres"
)

type sessionRepo struct {
	*postgres.Postgres
}

func NewSessionRepository(postgres *postgres.Postgres) session.SessionRepository {
	return &sessionRepo{
		postgres,
	}
}

func (s *sessionRepo) collectRow(row pgx.Row) (*entity.Session, error) {
	var session entity.Session
	err := row.Scan(&session.Token, &session.UserID, &session.ExpiresAt)
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

	return &session, err
}

func (s *sessionRepo) collectRows(rows pgx.Rows) ([]entity.Session, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Session, error) {
		session, err := s.collectRow(row)
		return *session, err
	})
}

func (s *sessionRepo) Create(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	query := `insert into session (token,user_id,expires_at) values ($1,$2,$3) returning *`

	row := s.Pool.QueryRow(ctx, query, session.Token, session.UserID, session.ExpiresAt)
	return s.collectRow(row)
}

func (s *sessionRepo) GetByToken(ctx context.Context, token string) (*entity.Session, error) {
	query := `select token, user_id, expires_at from session where token = $1`

	row := s.Pool.QueryRow(ctx, query, token)
	return s.collectRow(row)
}

func (s *sessionRepo) Update(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	query := `update session set token = $1, expires_at = $2 where user_id = $3 returning *`

	row := s.Pool.QueryRow(ctx, query, session.Token, session.ExpiresAt, session.UserID)
	return s.collectRow(row)
}
