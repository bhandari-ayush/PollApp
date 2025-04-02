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

func (p *PollStore) CreatePollOption(ctx context.Context, pollID int, optionText string) error {
	return withTx(p.db, ctx, func(tx *sql.Tx) error {
		_, err := p.createPollOption(ctx, tx, pollID, optionText)
		if err != nil {
			return err
		}
		return nil
	})
}

func (p *PollStore) createPollOption(ctx context.Context, tx *sql.Tx, pollID int, optionText string) (int, error) {
	query := "INSERT INTO poll_options (poll_id, option_text) VALUES ($1, $2) RETURNING id"
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var optionID int
	err := tx.QueryRowContext(ctx, query, pollID, optionText).Scan(&optionID)
	if err != nil {
		return 0, fmt.Errorf("error creating poll option: %v", err)
	}
	return optionID, nil
}

func (p *PollStore) GetPollOptions(ctx context.Context, pollID int) ([]PollOption, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, "SELECT id, poll_id, option_text, vote_count FROM poll_options WHERE poll_id = $1", pollID)
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
