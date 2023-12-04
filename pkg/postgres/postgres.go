package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func New(ctx context.Context, maxAttempts int, url string) (*Postgres, error) {

	db := &Postgres{}

	err := DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err := pgxpool.New(ctx, url)
		if err != nil {
			return err
		}

		err = pool.Ping(ctx)
		if err != nil {
			return err
		}

		db = &Postgres{
			Pool: pool,
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return db, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}
