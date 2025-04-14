package store

import (
	"PollApp/payload"
	"context"
	"database/sql"
	"fmt"
)

type Poll struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	CreatorId   int    `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
}

type PollStore struct {
	db *sql.DB
}

func (p *PollStore) Create(ctx context.Context, pollRequest *payload.PollRequest) (id int, err error) {
	err = withTx(p.db, ctx, func(tx *sql.Tx) error {
		id, err = p.createPoll(ctx, tx, pollRequest)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (p *PollStore) createPoll(ctx context.Context, tx *sql.Tx, pollRequest *payload.PollRequest) (int, error) {
	query := "INSERT INTO polls (description, creator_id) VALUES ($1, $2) RETURNING id"
	var pollId int

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, pollRequest.Description, pollRequest.CreatorId).Scan(&pollId)
	if err != nil {
		return 0, fmt.Errorf("error creating poll: %v", err)
	}
	return pollId, nil
}

func (p *PollStore) GetByID(ctx context.Context, pollId int) (*Poll, error) {

	poll := &Poll{}
	query := "SELECT id, description, creator_id, created_at FROM polls WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := p.db.QueryRowContext(ctx, query, pollId).Scan(&poll.Id, &poll.Description, &poll.CreatorId, &poll.CreatedAt)
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

	rows, err := p.db.QueryContext(ctx, "SELECT id, description, creator_id, created_at FROM polls")
	if err != nil {
		return nil, fmt.Errorf("error fetching polls: %v", err)
	}
	defer rows.Close()

	polls := make([]*Poll, 0)
	for rows.Next() {
		poll := &Poll{}
		if err := rows.Scan(&poll.Id, &poll.Description, &poll.CreatorId, &poll.CreatedAt); err != nil {
			return nil, fmt.Errorf("error reading poll row: %v", err)
		}
		polls = append(polls, poll)
	}

	return polls, nil
}

func (p *PollStore) Delete(ctx context.Context, pollId int) error {
	query := `DELETE FROM polls WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := p.db.ExecContext(ctx, query, pollId)
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
