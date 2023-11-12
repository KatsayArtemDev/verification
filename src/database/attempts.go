package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type Attempts interface {
	CheckIfUserExist(userId uuid.UUID) (int, error)
	InitNewUser(userId uuid.UUID) error
	IncrementUserAttempts(userId uuid.UUID) error
	GetUserAttempts(userId uuid.UUID) (int, error)
	DeleteUser(userId uuid.UUID) error
}

type attemptsExecutor struct {
	db *sql.DB
}

func NewAttempts(db *sql.DB) Attempts {
	return attemptsExecutor{db: db}
}

func (ex attemptsExecutor) CheckIfUserExist(userId uuid.UUID) (int, error) {
	var isUserExistInAttemptsTable int

	var selectUser = `select count(*) from attempts where user_id = $1`
	err := ex.db.QueryRow(selectUser, userId).Scan(&isUserExistInAttemptsTable)
	if err != nil {
		return 0, fmt.Errorf("error when selecting user from attemps table: %w", err)
	}

	return isUserExistInAttemptsTable, nil
}

func (ex attemptsExecutor) InitNewUser(userId uuid.UUID) error {
	var insertNewUser = `insert into attempts (user_id) values ($1)`
	_, err := ex.db.Exec(insertNewUser, userId)
	if err != nil {
		return fmt.Errorf("error when inserting new user to attempts table: %w", err)
	}

	return nil
}

func (ex attemptsExecutor) IncrementUserAttempts(userId uuid.UUID) error {
	var updateUserAttempts = `update attempts set attempts = attempts + 1 where user_id = $1`
	_, err := ex.db.Exec(updateUserAttempts, userId)
	if err != nil {
		return fmt.Errorf("error when updating user attempts: %w", err)
	}

	return nil
}

func (ex attemptsExecutor) GetUserAttempts(userId uuid.UUID) (int, error) {
	var attempts int

	var selectUserAttempts = `select attempts from attempts where user_id = $1`
	err := ex.db.QueryRow(selectUserAttempts, userId).Scan(&attempts)
	if err != nil {
		return 0, fmt.Errorf("error when selecting user attempts: %w", err)
	}

	return attempts, nil
}

func (ex attemptsExecutor) DeleteUser(userId uuid.UUID) error {
	var deleteUser = `delete from attempts where user_id = $1`
	_, err := ex.db.Exec(deleteUser, userId)
	if err != nil {
		return fmt.Errorf("error when deleting user from attempts table: %w", err)
	}

	return nil
}
