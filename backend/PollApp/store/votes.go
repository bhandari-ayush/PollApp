package store

import (
	"PollApp/payload"
	"context"
	"database/sql"
	"fmt"
)

type Vote struct {
	Id        int    `json:"id"`
	PollId    int    `json:"poll_id"`
	OptionId  int    `json:"option_id"`
	UserId    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type VoteStore struct {
	db *sql.DB
}

func (v *VoteStore) Create(ctx context.Context, voteRequest *payload.VoteRequest) (id int, err error) {
	err = withTx(v.db, ctx, func(tx *sql.Tx) error {
		id, err = v.createVote(ctx, tx, voteRequest)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (v *VoteStore) createVote(ctx context.Context, tx *sql.Tx, voteRequest *payload.VoteRequest) (int, error) {
	query := "INSERT INTO votes (poll_id, option_id, user_id) VALUES ($1, $2, $3) RETURNING id"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var voteId int
	err := tx.QueryRowContext(ctx, query, voteRequest.PollId, voteRequest.OptionId, voteRequest.UserId).Scan(&voteId)
	if err != nil {
		return 0, fmt.Errorf("error creating vote: %v", err)
	}

	err = v.increaseVoteCount(ctx, tx, voteRequest.OptionId)
	if err != nil {
		return 0, fmt.Errorf("error updating vote count: %v", err)
	}
	return voteId, nil
}

func (v *VoteStore) Update(ctx context.Context, voteRequest *payload.VoteRequest) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {
		err := v.updateVote(ctx, tx, voteRequest)
		if err != nil {
			return err
		}
		return nil
	})
}

func (v *VoteStore) updateVote(ctx context.Context, tx *sql.Tx, voteRequest *payload.VoteRequest) error {

	ctx, cancel := context.WithTimeout(ctx, 2*QueryTimeoutDuration)
	defer cancel()

	previousOptionID, err := v.getPreviousVote(ctx, tx, voteRequest.UserId, voteRequest.PollId)
	if err != nil {
		return fmt.Errorf("error checking previous vote: %v", err)
	}

	if previousOptionID != 0 && previousOptionID != voteRequest.OptionId {
		deleteVoteQuery := "DELETE FROM votes WHERE user_id = $1 AND poll_id = $2"
		_, err := tx.ExecContext(ctx, deleteVoteQuery, voteRequest.UserId, voteRequest.PollId)
		if err != nil {
			return fmt.Errorf("error deleting previous vote: %v", err)
		}

		err = v.decreaseVoteCount(ctx, tx, previousOptionID)
		if err != nil {
			return fmt.Errorf("error decreasing previous vote count: %v", err)
		}
	}

	insertVoteQuery := "INSERT INTO votes (poll_id, option_id, user_id) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, insertVoteQuery, voteRequest.PollId, voteRequest.OptionId, voteRequest.UserId)
	if err != nil {
		return fmt.Errorf("error inserting new vote: %v", err)
	}

	err = v.increaseVoteCount(ctx, tx, voteRequest.OptionId)
	if err != nil {
		return fmt.Errorf("error increasing vote count: %v", err)
	}

	return nil
}

func (v *VoteStore) Delete(ctx context.Context, voteRequest *payload.VoteRequest) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {
		err := v.deleteVote(ctx, tx, voteRequest)
		if err != nil {
			return err
		}
		return nil
	})
}

func (v *VoteStore) deleteVote(ctx context.Context, tx *sql.Tx, voteRequest *payload.VoteRequest) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	previousOptionID, err := v.getPreviousVote(ctx, tx, voteRequest.UserId, voteRequest.PollId)
	if err != nil {
		return fmt.Errorf("error checking previous vote: %v", err)
	}

	if previousOptionID != 0 && previousOptionID != voteRequest.OptionId {
		deleteVoteQuery := "DELETE FROM votes WHERE user_id = $1 AND poll_id = $2"
		_, err := tx.ExecContext(ctx, deleteVoteQuery, voteRequest.UserId, voteRequest.PollId)
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

func (v *VoteStore) DeleteByPollID(ctx context.Context, pollId int) error {
	return withTx(v.db, ctx, func(tx *sql.Tx) error {
		err := v.deleteVotesByPollID(ctx, tx, pollId)
		if err != nil {
			return err
		}
		return nil
	})
}

func (v *VoteStore) deleteVotesByPollID(ctx context.Context, tx *sql.Tx, pollId int) error {
	query := "DELETE FROM votes WHERE poll_id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, pollId)
	if err != nil {
		return fmt.Errorf("error deleting votes for poll_id %d: %v", pollId, err)
	}
	return nil
}

func (v *VoteStore) GetUsersForOption(ctx context.Context, optionId int) ([]int, error) {
	query := `SELECT user_id FROM votes WHERE option_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := v.db.QueryContext(ctx, query, optionId)
	if err != nil {
		fmt.Errorf("Error querying users for option %d: %v", optionId, err)
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			fmt.Errorf("Error reading user_id: %v", err)
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}

func (v *VoteStore) increaseVoteCount(ctx context.Context, tx *sql.Tx, optionId int) error {
	query := "UPDATE poll_options SET vote_count = vote_count + 1 WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, optionId)
	if err != nil {
		return fmt.Errorf("error updating vote count: %v", err)
	}
	return nil
}

func (v *VoteStore) decreaseVoteCount(ctx context.Context, tx *sql.Tx, optionId int) error {
	query := "UPDATE poll_options SET vote_count = vote_count - 1 WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, optionId)
	if err != nil {
		return fmt.Errorf("error decreasing vote count: %v", err)
	}
	return nil
}

func (v *VoteStore) getPreviousVote(ctx context.Context, tx *sql.Tx, userId, pollId int) (int, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := "SELECT option_id FROM votes WHERE user_id = $1 AND poll_id = $2 LIMIT 1"
	var optionID int
	err := tx.QueryRowContext(ctx, query, userId, pollId).Scan(&optionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error retrieving previous vote: %v", err)
	}
	return optionID, nil
}
