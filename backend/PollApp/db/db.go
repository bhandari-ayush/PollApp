package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	openOnce   sync.Once
	closeOnce  sync.Once
)

func DisconnectDB() error {
	var err error
	closeOnce.Do(func() {
		if dbInstance != nil {
			err = dbInstance.Close()
			if err != nil {
				log.Printf("Error closing the database connection: %v", err)
			} else {
				fmt.Println("Disconnected from the database successfully!")
			}
		}
	})
	return err
}

func ConnectDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	var err error
	openOnce.Do(func() {
		_, err = new(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	})
	return dbInstance, err
}

func new(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	dbInstance = db
	return db, nil
}
