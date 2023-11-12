package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Pins interface {
	AddNewUser(userId uuid.UUID, pin string) error
	GetUserDbData(userId uuid.UUID) (userDbData, error)
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

type userDbData struct {
	DbPin       string
	DbCreatedAt time.Time
}

func (ex pinsExecutor) GetUserDbData(userId uuid.UUID) (userDbData, error) {
	var userDbData userDbData

	var selectUserPinData = `select pin, created_at from pins where user_id = $1`
	err := ex.db.QueryRow(selectUserPinData, userId).Scan(&userDbData.DbPin, &userDbData.DbCreatedAt)
	if err != nil {
		return userDbData, fmt.Errorf("error when quering pin db data: %w", err)
	}

	return userDbData, nil
}

func (ex pinsExecutor) DeleteUser(userId uuid.UUID) error {
	var deleteUser = `delete from pins where user_id = $1`
	_, err := ex.db.Exec(deleteUser, userId)
	if err != nil {
		return fmt.Errorf("error when deleting user from pins table: %w", err)
	}

	return nil
}
