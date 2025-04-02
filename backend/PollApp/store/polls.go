package store

import (
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

// CreatePoll creates a new poll
func CreatePoll(ctx context.Context, db *sql.DB, title, description string, creatorID int) (int, error) {
	query := "INSERT INTO polls (title, description, creator_id) VALUES ($1, $2, $3) RETURNING id"
	var pollID int

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := db.QueryRowContext(ctx, query, title, description, creatorID).Scan(&pollID)
	if err != nil {
		return 0, fmt.Errorf("error creating poll: %v", err)
	}
	return pollID, nil
}

// GetPoll fetches a specific poll by its ID
func GetPoll(ctx context.Context, db *sql.DB, pollID int) (Poll, error) {

	var poll Poll
	query := "SELECT id, title, description, creator_id, created_at FROM polls WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := db.QueryRowContext(ctx, query, pollID).Scan(&poll.ID, &poll.Title, &poll.Description, &poll.CreatorID, &poll.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Poll{}, fmt.Errorf("poll not found")
		}
		return Poll{}, fmt.Errorf("error fetching poll: %v", err)
	}
	return poll, nil
}

// ListPolls fetches all polls
func ListPolls(ctx context.Context, db *sql.DB) ([]Poll, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT id, title, description, creator_id, created_at FROM polls")
	if err != nil {
		return nil, fmt.Errorf("error fetching polls: %v", err)
	}
	defer rows.Close()

	var polls []Poll
	for rows.Next() {
		var poll Poll
		if err := rows.Scan(&poll.ID, &poll.Title, &poll.Description, &poll.CreatorID, &poll.CreatedAt); err != nil {
			return nil, fmt.Errorf("error reading poll row: %v", err)
		}
		polls = append(polls, poll)
	}

	return polls, nil
}
