package store

import (
	"PollApp/payload"
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

type VoteStore struct {
	db *sql.DB
}

func (v *VoteStore) Create(ctx context.Context, db *sql.DB, voteRequest *payload.VoteRequest) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {
		err := v.createVote(ctx, tx, voteRequest)
		if err != nil {
			return err
		}
		return nil
	})
}

func (v *VoteStore) createVote(ctx context.Context, tx *sql.Tx, voteRequest *payload.VoteRequest) error {
	query := "INSERT INTO votes (poll_id, option_id, user_id) VALUES ($1, $2, $3)"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, voteRequest.PollID, voteRequest.OptionID, voteRequest.UserID)
	if err != nil {
		return fmt.Errorf("error creating vote: %v", err)
	}

	err = v.increaseVoteCount(ctx, tx, voteRequest.OptionID)
	if err != nil {
		return fmt.Errorf("error updating vote count: %v", err)
	}
	return nil
}

func (v *VoteStore) Update(ctx context.Context, db *sql.DB, voteRequest *payload.VoteRequest) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {
		err := v.updateVote(ctx, tx, voteRequest)
		if err != nil {
			return err
		}
		return nil
	})
}

func (v *VoteStore) updateVote(ctx context.Context, tx *sql.Tx, voteRequest *payload.VoteRequest) error {

	previousOptionID, err := v.getPreviousVote(ctx, tx, voteRequest.UserID, voteRequest.PollID)
	if err != nil {
		return fmt.Errorf("error checking previous vote: %v", err)
	}

	if previousOptionID != 0 && previousOptionID != voteRequest.OptionID {
		deleteVoteQuery := "DELETE FROM votes WHERE user_id = $1 AND poll_id = $2"
		_, err := tx.ExecContext(ctx, deleteVoteQuery, voteRequest.UserID, voteRequest.PollID)
		if err != nil {
			return fmt.Errorf("error deleting previous vote: %v", err)
		}

		err = v.decreaseVoteCount(ctx, tx, previousOptionID)
		if err != nil {
			return fmt.Errorf("error decreasing previous vote count: %v", err)
		}
	}

	insertVoteQuery := "INSERT INTO votes (poll_id, option_id, user_id) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, insertVoteQuery, voteRequest.PollID, voteRequest.OptionID, voteRequest.UserID)
	if err != nil {
		return fmt.Errorf("error inserting new vote: %v", err)
	}

	err = v.increaseVoteCount(ctx, tx, voteRequest.OptionID)
	if err != nil {
		return fmt.Errorf("error increasing vote count: %v", err)
	}

	return nil
}

func (v *VoteStore) Delete(ctx context.Context, db *sql.DB, voteRequest *payload.VoteRequest) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {
		err := v.deleteVote(ctx, tx, voteRequest)
		if err != nil {
			return err
		}
		return nil
	})
}

func (v *VoteStore) deleteVote(ctx context.Context, tx *sql.Tx, voteRequest *payload.VoteRequest) error {

	previousOptionID, err := v.getPreviousVote(ctx, tx, voteRequest.UserID, voteRequest.PollID)
	if err != nil {
		return fmt.Errorf("error checking previous vote: %v", err)
	}

	if previousOptionID != 0 && previousOptionID != voteRequest.OptionID {
		deleteVoteQuery := "DELETE FROM votes WHERE user_id = $1 AND poll_id = $2"
		_, err := tx.ExecContext(ctx, deleteVoteQuery, voteRequest.UserID, voteRequest.PollID)
		if err != nil {
			return fmt.Errorf("error deleting previous vote: %v", err)
		}

		err = v.decreaseVoteCount(ctx, tx, previousOptionID)
		if err != nil {
			return fmt.Errorf("error decreasing previous vote count: %v", err)
		}
	}

	return nil
}

func (v *VoteStore) increaseVoteCount(ctx context.Context, tx *sql.Tx, optionID int) error {
	query := "UPDATE poll_options SET vote_count = vote_count + 1 WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, optionID)
	if err != nil {
		return fmt.Errorf("error updating vote count: %v", err)
	}
	return nil
}

func (v *VoteStore) decreaseVoteCount(ctx context.Context, tx *sql.Tx, optionID int) error {
	query := "UPDATE poll_options SET vote_count = vote_count - 1 WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, optionID)
	if err != nil {
		return fmt.Errorf("error decreasing vote count: %v", err)
	}
	return nil
}

func (v *VoteStore) getPreviousVote(ctx context.Context, tx *sql.Tx, userID, pollID int) (int, error) {
	query := "SELECT option_id FROM votes WHERE user_id = $1 AND poll_id = $2 LIMIT 1"
	var optionID int
	err := tx.QueryRowContext(ctx, query, userID, pollID).Scan(&optionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error retrieving previous vote: %v", err)
	}
	return optionID, nil
}

// func (v *VoteStore) UpdateVoteCount(ctx context.Context, optionID int) error {
// 	return withTx(v.db, ctx, func(tx *sql.Tx) error {
// 		err := v.updateVoteCount(ctx, tx, optionID)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
