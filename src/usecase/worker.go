package usecase

import (
	"database/sql"
	"fmt"
	"github.com/KatsayArtemDev/verification/src/database"
	"github.com/KatsayArtemDev/verification/src/processing"
	"github.com/KatsayArtemDev/verification/src/sending"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Worker interface {
	AttemptsProcessing(userId uuid.UUID) error
	BlocksProcessing(userId uuid.UUID) (int, error)
	SendingProcessing(userId uuid.UUID, email string) error
	GetUserDataAndCalcTimeDiff(userId uuid.UUID) (string, time.Duration, error)
	DeleteUserFromPinsAndAttemptsTables(userId uuid.UUID) error
}

type worker struct {
	pins     database.Pins
	attempts database.Attempts
	blocks   database.Blocks
}

func NewWorker(executor *sql.DB) Worker {
	return worker{
		pins:     database.NewPins(executor),
		attempts: database.NewAttempts(executor),
		blocks:   database.NewBlocks(executor),
	}
}

func (w worker) AttemptsProcessing(userId uuid.UUID) error {
	isUserExistInAttemptsTable, err := w.attempts.CheckIfUserExist(userId)
	if err != nil {
		return fmt.Errorf("failed to check if user in attempts table: %w", err)
	}

	if isUserExistInAttemptsTable == 0 {
		err = w.attempts.InitNewUser(userId)
		if err != nil {
			return fmt.Errorf("failed to initialize new user: %w", err)
		}

		return nil
	}

	err = w.attempts.IncrementUserAttempts(userId)
	if err != nil {
		return fmt.Errorf("failed to increment user attempts: %w", err)
	}

	attempts, err := w.attempts.GetUserAttempts(userId)
	if err != nil {
		return fmt.Errorf("failed to select user attempts: %w", err)
	}

	if attempts > 5 {
		err = w.blocks.AddNewUser(userId)
		if err != nil {
			return fmt.Errorf("failed to add new user: %w", err)
		}

		err = w.attempts.DeleteUser(userId)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
	}

	return nil
}

func (w worker) BlocksProcessing(userId uuid.UUID) (int, error) {
	isUserExistInBlocksTable, err := w.blocks.CheckIfUserExist(userId)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to check if user in blocks table: %w", err)
	}

	if isUserExistInBlocksTable == 0 {
		return 0, nil
	}

	dbBlockedAt, err := w.blocks.GetUserDbData(userId)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to check if user in blocks table: %w", err)
	}

	temporaryBlockingChange := processing.TimeChecking(dbBlockedAt)

	if temporaryBlockingChange.Minutes() < 2 {
		return http.StatusUnauthorized, fmt.Errorf("user is blocked")
	}

	err = w.blocks.DeleteUser(userId)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err)
	}

	return 0, nil
}

func (w worker) SendingProcessing(userId uuid.UUID, email string) error {
	// Get and hash generated pin
	hash, pin, err := processing.PinProcessing()
	if err != nil {
		return fmt.Errorf("failed to get and hash pin: %w", err)
	}

	// Check if pin was sent to user
	isUserExistInPinsTable, err := w.pins.CheckIfUserExist(userId)
	if err != nil {
		return fmt.Errorf("failed to check if user exist in pins table: %w", err)
	}

	if isUserExistInPinsTable == 0 {
		// Add new user to table
		err = w.pins.AddNewUser(userId, hash)
		if err != nil {
			return fmt.Errorf("failed to add new user: %w", err)
		}
	} else {
		// Set new pin
		err = w.pins.UpdateUserPin(userId, hash)
		if err != nil {
			return fmt.Errorf("failed to update new pin: %w", err)
		}
	}

	err = w.pins.UpdatePinSentAt(userId)
	if err != nil {
		return fmt.Errorf("failed to set sent at: %w", err)
	}

	// Send pin to user
	err = sending.PinToUser(email, pin)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (w worker) GetUserDataAndCalcTimeDiff(userId uuid.UUID) (string, time.Duration, error) {
	// Get user db data
	userDbData, err := w.pins.GetUserDbData(userId)
	if err != nil {
		return "", 0, fmt.Errorf("failed to bind pin body: %w", err)
	}

	var (
		dbPin    = userDbData.DbPin
		dbSentAt = userDbData.DbSentAt
	)

	temporarySendingChange := processing.TimeChecking(dbSentAt)

	return dbPin, temporarySendingChange, nil
}

func (w worker) DeleteUserFromPinsAndAttemptsTables(userId uuid.UUID) error {
	isUserExistInAttemptsTable, err := w.attempts.CheckIfUserExist(userId)
	if err != nil {
		return fmt.Errorf("failed to check if user in attempts table: %w", err)
	}

	if isUserExistInAttemptsTable != 0 {
		err = w.attempts.DeleteUser(userId)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
	}

	// Delete verified user
	err = w.pins.DeleteUser(userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
