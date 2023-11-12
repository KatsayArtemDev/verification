package usecase

import (
	"database/sql"
	"fmt"
	"github.com/KatsayArtemDev/verification/src/database"
	"github.com/google/uuid"
)

type PinWorker interface {
	AttemptsProcessing(userId uuid.UUID) error
}

type pinWorker struct {
	attempts database.Attempts
	blocks   database.Blocks
}

func NewPinWorker(executor *sql.DB) PinWorker {
	return pinWorker{
		attempts: database.NewAttempts(executor),
		blocks:   database.NewBlocks(executor),
	}
}

func (pw pinWorker) AttemptsProcessing(userId uuid.UUID) error {
	isUserExistInAttemptsTable, err := pw.attempts.CheckIfUserExist(userId)
	if err != nil {
		return fmt.Errorf("failed to check if user in attempts table: %w", err)
	}

	if isUserExistInAttemptsTable == 0 {
		err = pw.attempts.InitNewUser(userId)
		if err != nil {
			return fmt.Errorf("failed to initialize new user: %w", err)
		}
	} else {
		err = pw.attempts.IncrementUserAttempts(userId)
		if err != nil {
			return fmt.Errorf("failed to increment user attempts: %w", err)
		}

		attempts, err := pw.attempts.GetUserAttempts(userId)
		if err != nil {
			return fmt.Errorf("failed to select user attempts: %w", err)
		}

		if attempts > 5 {
			err = pw.blocks.AddNewUser(userId)
			if err != nil {
				return fmt.Errorf("failed to add new user: %w", err)
			}

			err = pw.attempts.DeleteUser(userId)
			if err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}
		}
	}

	return nil
}
