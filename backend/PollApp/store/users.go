package store

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func NewUser(username, password, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) (id int, err error) {
	err = withTx(u.db, ctx, func(tx *sql.Tx) error {
		id, err = u.createUser(ctx, tx, user)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (u *UserStore) createUser(ctx context.Context, tx *sql.Tx, user *User) (int, error) {
	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id"
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var userId int
	err := tx.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %v", err)
	}
	return userId, nil
}

func (u *UserStore) GetByID(ctx context.Context, userId int) (*User, error) {
	user := &User{}
	query := "SELECT id, username, password, email, created_at FROM users WHERE id = $1"
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := u.db.QueryRowContext(ctx, query, userId).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &User{}, fmt.Errorf("user not found")
		}
		return &User{}, fmt.Errorf("error fetching user: %v", err)
	}
	return user, nil
}

func (u *UserStore) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := u.db.ExecContext(ctx, query, userId)
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
