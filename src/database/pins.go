package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Pins interface {
	AddNewUser(userId uuid.UUID, hash string) error
	CheckIfUserExist(userId uuid.UUID) (int, error)
	GetUserDbData(userId uuid.UUID) (userDbData, error)
	UpdateUserPin(userId uuid.UUID, hash string) error
	UpdatePinSentAt(userId uuid.UUID) error
	DeleteUser(userId uuid.UUID) error
}

type pinsExecutor struct {
	db *sql.DB
}

func NewPins(db *sql.DB) Pins {
	return pinsExecutor{db: db}
}

func (ex pinsExecutor) AddNewUser(userId uuid.UUID, hash string) error {
	var insertNewUser = `insert into pins (user_id, pin) values($1, $2)`
	_, err := ex.db.Exec(insertNewUser, userId, hash)
	if err != nil {
		return fmt.Errorf("error when saving pin hash: %w", err)
	}

	return nil
}

func (ex pinsExecutor) CheckIfUserExist(userId uuid.UUID) (int, error) {
	var isUserExistInPinsTable int

	var selectUser = `select count(*) from pins where user_id = $1`
	err := ex.db.QueryRow(selectUser, userId).Scan(&isUserExistInPinsTable)
	if err != nil {
		return 0, fmt.Errorf("error when selecting user from pins table: %w", err)
	}

	return isUserExistInPinsTable, nil
}

type userDbData struct {
	DbPin    string
	DbSentAt time.Time
}

func (ex pinsExecutor) GetUserDbData(userId uuid.UUID) (userDbData, error) {
	var userDbData userDbData

	var selectUserPinData = `select pin, sent_at from pins where user_id = $1`
	err := ex.db.QueryRow(selectUserPinData, userId).Scan(&userDbData.DbPin, &userDbData.DbSentAt)
	if err != nil {
		return userDbData, fmt.Errorf("error when quering pin db data: %w", err)
	}

	return userDbData, nil
}

func (ex pinsExecutor) UpdateUserPin(userId uuid.UUID, hash string) error {
	var updateUserPin = `update pins set pin = $1 where user_id = $2`
	_, err := ex.db.Exec(updateUserPin, hash, userId)
	if err != nil {
		return fmt.Errorf("error when updating user pin in pins table: %w", err)
	}

	return nil
}

func (ex pinsExecutor) UpdatePinSentAt(userId uuid.UUID) error {
	var sentAtNow = time.Now().UTC()

	var updateUserPin = `update pins set sent_at = $1 where user_id = $2`
	_, err := ex.db.Exec(updateUserPin, sentAtNow, userId)
	if err != nil {
		return fmt.Errorf("error when updating pin sent at in pins table: %w", err)
	}

	return nil
}

func (ex pinsExecutor) DeleteUser(userId uuid.UUID) error {
	var deleteUser = `delete from pins where user_id = $1`
	_, err := ex.db.Exec(deleteUser, userId)
	if err != nil {
		return fmt.Errorf("error when deleting user from pins table: %w", err)
	}

	return nil
}
