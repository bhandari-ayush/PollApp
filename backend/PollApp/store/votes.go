package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Vote struct {
	ID        int    `json:"id"`
	PollID    int    `json:"poll_id"`
	OptionID  int    `json:"option_id"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func CreateVote(ctx context.Context, db *sql.DB, pollID, optionID, userID int) error {
	query := "INSERT INTO votes (poll_id, option_id, user_id) VALUES ($1, $2, $3)"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := db.ExecContext(ctx, query, pollID, optionID, userID)
	if err != nil {
		return fmt.Errorf("error creating vote: %v", err)
	}
	return nil
}

func UpdateVoteCount(ctx context.Context, db *sql.DB, optionID int) error {
	query := "UPDATE poll_options SET vote_count = vote_count + 1 WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := db.ExecContext(ctx, query, optionID)
	if err != nil {
		return fmt.Errorf("error updating vote count: %v", err)
	}
	return nil
}
