package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-todolist-sber/internal/apperror"
	pgxError "go-todolist-sber/internal/apperror/pgx_errors"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/task"
	"go-todolist-sber/pkg/postgres"
	"strings"
	"time"
)

type taskRepository struct {
	*postgres.Postgres
}

func NewTaskRepository(postgres *postgres.Postgres) task.TaskRepository {
	return &taskRepository{
		postgres,
	}
}

func (t *taskRepository) collectRow(row pgx.Row) (*entity.Task, error) {
	var task entity.Task
	err := row.Scan(&task.ID, &task.UserID, &task.Header, &task.Description, &task.CreatedAt, &task.StartDate, &task.Done)
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

	return &task, err
}

func (t *taskRepository) collectRows(rows pgx.Rows) ([]entity.Task, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Task, error) {
		task, err := t.collectRow(row)
		return *task, err
	})
}

func (t *taskRepository) Update(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	var builder strings.Builder
	var increment int
	attribute := []interface{}{}

	builder.WriteString(`update task set `)

	if task.Header != "" {
		increment++
		attribute = append(attribute, task.Header)
		builder.WriteString(fmt.Sprintf(`header = $%d `, increment))
	}
	if task.Description != "" {
		increment++
		attribute = append(attribute, task.Description)
		builder.WriteString(fmt.Sprintf(`,description = $%d `, increment))
	}
	if !task.StartDate.IsZero() {
		increment++
		attribute = append(attribute, task.StartDate)
		builder.WriteString(fmt.Sprintf(`,created_at = $%d `, increment))
	}
	increment++
	builder.WriteString(fmt.Sprintf(`where id = $%d returning *`, increment))
	attribute = append(attribute, task.ID)
	row := t.Pool.QueryRow(ctx, builder.String(), attribute...)

	return t.collectRow(row)
}

func (t *taskRepository) Create(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	query := `insert into task (id_user,header,description,start_date) values ($1,$2,$3,$4) returning *`

	row := t.Pool.QueryRow(ctx, query, task.UserID, task.Header, task.Description, task.StartDate)
	return t.collectRow(row)
}

func (t *taskRepository) GetByUserID(ctx context.Context, id string) ([]entity.Task, error) {
	query := `select id, id_user, header, description, created_at, start_date, done from task
				where id_user = $1`

	rows, err := t.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return t.collectRows(rows)
}

func (t *taskRepository) GetAll(ctx context.Context) ([]entity.Task, error) {
	query := `select id, id_user, header, description, created_at, start_date, done from task`

	rows, err := t.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return t.collectRows(rows)
}

func (t *taskRepository) DeleteByID(ctx context.Context, userID string, id int) error {
	query := `delete from task where id = $1 and id_user = $2`

	_, err := t.Pool.Exec(ctx, query, id, userID)
	return err
}

func (t *taskRepository) GetPageByStatusAndUserID(ctx context.Context, userID string, status bool, offset int) ([]entity.Task, error) {
	query := `select id, id_user, header, description, created_at, start_date, done
				from task
				where id_user = $1 and done = $2
				order by id desc
				offset $3
				limit 3`

	rows, err := t.Pool.Query(ctx, query, userID, status, offset)
	if err != nil {
		return nil, err
	}

	return t.collectRows(rows)
}

func (t *taskRepository) GetByDateAndStatus(ctx context.Context, userID string, date time.Time, status bool) ([]entity.Task, error) {
	query := `select id, id_user, header, description, created_at, start_date, done
				from task
				where id_user = $1 and done = $2 and start_date = $3
				order by id desc`

	rows, err := t.Pool.Query(ctx, query, userID, status, date)
	if err != nil {
		return nil, err
	}

	return t.collectRows(rows)
}
