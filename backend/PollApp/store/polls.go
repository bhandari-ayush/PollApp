package store

import (
	"PollApp/payload"
	"context"
	"database/sql"
	"fmt"
)

type Poll struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorID   int    `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
}

type PollStore struct {
	db *sql.DB
}

func (p *PollStore) Create(ctx context.Context, pollRequest *payload.PollRequest) error {
	return withTx(p.db, ctx, func(tx *sql.Tx) error {
		_, err := p.createPoll(ctx, tx, pollRequest)
		if err != nil {
			return err
		}
		return nil
	})
}

func (p *PollStore) createPoll(ctx context.Context, tx *sql.Tx, pollRequest *payload.PollRequest) (int, error) {
	query := "INSERT INTO polls (title, description, creator_id) VALUES ($1, $2, $3) RETURNING id"
	var pollID int

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, pollRequest.Title, pollRequest.Description, pollRequest.CreatorID).Scan(&pollID)
	if err != nil {
		return 0, fmt.Errorf("error creating poll: %v", err)
	}
	return pollID, nil
}

func (p *PollStore) GetByID(ctx context.Context, pollID int) (*Poll, error) {

	poll := &Poll{}
	query := "SELECT id, title, description, creator_id, created_at FROM polls WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := p.db.QueryRowContext(ctx, query, pollID).Scan(&poll.ID, &poll.Title, &poll.Description, &poll.CreatorID, &poll.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Poll{}, fmt.Errorf("poll not found")
		}
		return &Poll{}, fmt.Errorf("error fetching poll: %v", err)
	}
	return poll, nil
}

func (p *PollStore) ListPolls(ctx context.Context) ([]*Poll, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, "SELECT id, title, description, creator_id, created_at FROM polls")
	if err != nil {
		return nil, fmt.Errorf("error fetching polls: %v", err)
	}
	defer rows.Close()

	polls := make([]*Poll, 0)
	for rows.Next() {
		poll := &Poll{}
		if err := rows.Scan(&poll.ID, &poll.Title, &poll.Description, &poll.CreatorID, &poll.CreatedAt); err != nil {
			return nil, fmt.Errorf("error reading poll row: %v", err)
		}
		polls = append(polls, poll)
	}

	return polls, nil
}

func (p *PollStore) Delete(ctx context.Context, pollID int) error {
	query := `DELETE FROM polls WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := p.db.ExecContext(ctx, query, pollID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
