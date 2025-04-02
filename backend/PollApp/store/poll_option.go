package store

import (
	"context"
	"database/sql"
	"fmt"
)

type PollOption struct {
	ID         int    `json:"id"`
	PollID     int    `json:"poll_id"`
	OptionText string `json:"option_text"`
	VoteCount  int    `json:"vote_count"`
}

func CreatePollOption(ctx context.Context, db *sql.DB, pollID int, optionText string) (int, error) {
	query := "INSERT INTO poll_options (poll_id, option_text) VALUES ($1, $2) RETURNING id"
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var optionID int
	err := db.QueryRowContext(ctx, query, pollID, optionText).Scan(&optionID)
	if err != nil {
		return 0, fmt.Errorf("error creating poll option: %v", err)
	}
	return optionID, nil
}

func GetPollOptions(ctx context.Context, db *sql.DB, pollID int) ([]PollOption, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT id, poll_id, option_text, vote_count FROM poll_options WHERE poll_id = $1", pollID)
	if err != nil {
		return nil, fmt.Errorf("error fetching poll options: %v", err)
	}
	defer rows.Close()

	var options []PollOption
	for rows.Next() {
		var option PollOption
		if err := rows.Scan(&option.ID, &option.PollID, &option.OptionText, &option.VoteCount); err != nil {
			return nil, fmt.Errorf("error reading option row: %v", err)
		}
		options = append(options, option)
	}

	return options, nil
}
