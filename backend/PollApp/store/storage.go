package store

import (
	"PollApp/payload"
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		GetByID(context.Context, int) (*User, error)
		GetByUsername(context.Context, string) (*User, error)
		Create(context.Context, *User) error
		Delete(context.Context, int) error
	}
	Polls interface {
		GetByID(ctx context.Context, pollID int) (*Poll, error)
		ListPolls(ctx context.Context) ([]*Poll, error)
		Create(context.Context, *payload.PollRequest) error
		Delete(ctx context.Context, pollID int) error
	}
	Votes interface {
		Create(ctx context.Context, db *sql.DB, voteRequest *payload.VoteRequest) error
		Update(ctx context.Context, db *sql.DB, voteRequest *payload.VoteRequest) error
		Delete(ctx context.Context, db *sql.DB, voteRequest *payload.VoteRequest) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
		Polls: &PollStore{db},
		Votes: &VoteStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
