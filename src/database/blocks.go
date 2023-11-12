package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Blocks interface {
	AddNewUser(userId uuid.UUID) error
	CheckIfUserExist(userId uuid.UUID) (int, error)
	GetUserDbData(userId uuid.UUID) (time.Time, error)
	DeleteUser(userId uuid.UUID) error
}

type blocksExecutor struct {
	db *sql.DB
}

func NewBlocks(db *sql.DB) Blocks {
	return blocksExecutor{db: db}
}

func (ex blocksExecutor) AddNewUser(userId uuid.UUID) error {
	var insertNewUser = `insert into blocks (user_id) values ($1)`
	_, err := ex.db.Exec(insertNewUser, userId)
	if err != nil {
		return fmt.Errorf("error when inserting new user to blocks table: %w", err)
	}

	return nil
}

func (ex blocksExecutor) CheckIfUserExist(userId uuid.UUID) (int, error) {
	var isUserExistInBlocksTable int

	var selectUser = `select count(*) from blocks where user_id = $1`
	err := ex.db.QueryRow(selectUser, userId).Scan(&isUserExistInBlocksTable)
	if err != nil {
		return 0, fmt.Errorf("error when selecting user from blocks table: %w", err)
	}

	return isUserExistInBlocksTable, nil
}

func (ex blocksExecutor) GetUserDbData(userId uuid.UUID) (time.Time, error) {
	var dbBlockedAt time.Time

	var selectUserPinData = `select blocked_at from blocks where user_id = $1`
	err := ex.db.QueryRow(selectUserPinData, userId).Scan(&dbBlockedAt)
	if err != nil {
		return dbBlockedAt, fmt.Errorf("error when quering blocks db data: %w", err)
	}

	return dbBlockedAt, nil
}

func (ex blocksExecutor) DeleteUser(userId uuid.UUID) error {
	var deleteUser = `delete from blocks where user_id = $1`
	_, err := ex.db.Exec(deleteUser, userId)
	if err != nil {
		return fmt.Errorf("error when deleting user from blocks table: %w", err)
	}

	return nil
}
