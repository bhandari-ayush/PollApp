package store

import (
	"PollApp/payload"
	"context"
	"database/sql"
	"errors"
	"log"
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
		Create(ctx context.Context, voteRequest *payload.VoteRequest) error
		Update(ctx context.Context, voteRequest *payload.VoteRequest) error
		Delete(ctx context.Context, voteRequest *payload.VoteRequest) error
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

var (
	CREATE_POLLS_TABLE = `CREATE TABLE IF NOT EXISTS polls (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    creator_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	CREATE_USERS_TABLE = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

	CREATE_VOTES_TABLE = `CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL,
    option_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (poll_id) REFERENCES polls(id),
    FOREIGN KEY (option_id) REFERENCES poll_options(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
`

	CREATES_POLLS_OPTION_TABLE = `CREATE TABLE IF NOT EXISTS poll_options (
    id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL,
    option_text VARCHAR(255) NOT NULL,
    vote_count INT DEFAULT 0,
    FOREIGN KEY (poll_id) REFERENCES polls(id)
);
`
)

func CreateTables(db *sql.DB, ctx context.Context) error {
	return withTx(db, ctx, func(tx *sql.Tx) error {
		err := createTables(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})
}

func createTables(ctx context.Context, tx *sql.Tx) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	tables := [4]string{
		CREATE_USERS_TABLE,
		CREATE_POLLS_TABLE,
		CREATES_POLLS_OPTION_TABLE,
		CREATE_VOTES_TABLE,
	}

	for _, table := range tables {
		log.Printf("executing cmd %s\n", table)
		_, err := tx.ExecContext(ctx, table)
		if err != nil {
			return err
		}
	}
	return nil
}
