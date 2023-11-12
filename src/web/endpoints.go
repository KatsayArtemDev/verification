package web

import (
	"fmt"
	"github.com/KatsayArtemDev/verification/src/processing"
	"github.com/KatsayArtemDev/verification/src/sending"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (app app) getEmail(c *gin.Context) {
	// Struct of request data
	var emailBody struct {
		UserId uuid.UUID `json:"user_id"`
		Email  string    `json:"email"`
	}

	// Bind struct values
	err := c.Bind(&emailBody)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to bind email body: %w", err))
		return
	}

	var (
		userId = emailBody.UserId
		email  = emailBody.Email
	)

	// Check email length
	if len(email) == 0 {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("email is empty"))
		return
	}

	// Validate email
	err = processing.EmailValidation(email)
	if err != nil {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("email is not valid: %w", err))
		return
	}

	// Generate new pin
	pin, err := processing.PinGenerating()
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to create new pin: %w", err))
		return
	}

	// Hash generated pin
	hash, err := processing.PinHashing(pin)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to hash pin: %w", err))
		return
	}

	// Add new user to table
	err = app.pins.AddNewUser(userId, hash)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to add new user: %w", err))
		return
	}

	// Send pin to user
	err = sending.PinToUser(email, pin)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to send email: %w", err))
		return
	}

	// Success end
	app.success(c, nil)
}

func (app app) getPin(c *gin.Context) {
	// Struct of request data
	var pinBody struct {
		UserId uuid.UUID `json:"user_id"`
		Pin    string    `json:"pin"`
	}

	// Bind struct values
	err := c.Bind(&pinBody)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to bind pin body: %w", err))
		return
	}

	var (
		userId = pinBody.UserId
		pin    = pinBody.Pin
	)

	// TODO: think about name @InBlocksTable@
	isUserExistInBlocksTable, err := app.blocks.CheckIfUserExist(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to check if user in blocks table: %w", err))
		return
	}

	if isUserExistInBlocksTable != 0 {
		dbBlockedAt, err := app.blocks.GetUserDbData(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to check if user in blocks table: %w", err))
			return
		}

		temporaryBlockingChange := processing.TimeChecking(dbBlockedAt)

		if temporaryBlockingChange.Minutes() < 2 {
			app.fail(c, http.StatusUnauthorized, fmt.Errorf("user is blocked"))
			return
		}

		err = app.blocks.DeleteUser(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
			return
		}
	}

	// Check pin length
	if len(pin) != 6 {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin is not valid"))
		return
	}

	// Get user db data
	userDbData, err := app.pins.GetUserDbData(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to bind pin body: %w", err))
		return
	}

	var (
		dbPin      = userDbData.DbPin
		dbCreateAt = userDbData.DbCreatedAt
	)

	temporaryCreatingChange := processing.TimeChecking(dbCreateAt)

	// Check difference between hours to remove expired
	if temporaryCreatingChange.Hours() > 2 {
		isUserExistInAttemptsTable, err := app.Attempts.CheckIfUserExist(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to check if user in attempts table: %w", err))
			return
		}

		if isUserExistInAttemptsTable == 0 {
			err = app.Attempts.InitNewUser(userId)
			if err != nil {
				app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to initialize new user: %w", err))
				return
			}
		} else {
			err = app.Attempts.IncrementUserAttempts(userId)
			if err != nil {
				app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to increment user attempts: %w", err))
				return
			}

			attempts, err := app.Attempts.GetUserAttempts(userId)
			if err != nil {
				app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to select user attempts: %w", err))
				return
			}

			if attempts > 5 {
				err = app.blocks.AddNewUser(userId)
				if err != nil {
					app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to add new user: %w", err))
					return
				}

				err = app.Attempts.DeleteUser(userId)
				if err != nil {
					app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
					return
				}
			}
		}

		err = app.pins.DeleteUser(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
		}

		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin has expired"))
		return
	}

	// Compare pin
	err = processing.PinComparing(dbPin, pin)
	if err != nil {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin is not correct"))
		return
	}

	// Delete verified user
	err = app.pins.DeleteUser(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
	}

	// Success end
	app.success(c, "confirmed")
}
