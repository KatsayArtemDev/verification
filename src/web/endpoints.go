package web

import (
	"fmt"
	"github.com/KatsayArtemDev/verification/src/processing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Struct of request email data
var emailBody struct {
	UserId uuid.UUID `json:"user_id"`
	//PinSeqNum uuid.UUID `json:"pin_seq_num"`
	Email string `json:"email"`
}

func (app app) getEmail(c *gin.Context) {
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

	// TODO: move checking exists to usecase
	isUserExistInBlocksTable, err := app.blocks.CheckIfUserExist(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to check if user exist in blocks table: %w", err))
	}

	if isUserExistInBlocksTable != 0 {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("user is blocked"))
		return
	}

	// Check if pin was sent to user
	isUserExistInPinsTable, err := app.pins.CheckIfUserExist(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to check if user exist in pins table: %w", err))
	}

	if isUserExistInPinsTable != 0 {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin was sent"))
		return
	}

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

	// Create new pin code and send to user
	err = app.worker.SendingProcessing(userId, email)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to send pin to user: %w", err))
		return
	}

	// Success end
	app.success(c, nil)
}

// Struct of request data
var pinBody struct {
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Pin    string    `json:"pin"`
}

func (app app) getPin(c *gin.Context) {
	// Bind struct values
	err := c.Bind(&pinBody)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to bind pin body: %w", err))
		return
	}

	var (
		userId = pinBody.UserId
		email  = pinBody.Email
		pin    = pinBody.Pin
	)

	requestStatus, err := app.worker.BlocksProcessing(userId)
	if err != nil {
		app.fail(c, requestStatus, fmt.Errorf("failed to execute blocks processing: %w", err))
		return
	}

	// Check pin length
	if len(pin) != 6 {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin is not valid"))
		return
	}

	dbPin, temporarySendingChange, err := app.worker.GetUserDataAndCalcTimeDiff(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to execute getting user data and calculate time difference: %w", err))
		return
	}

	// Check difference between hours to remove expired
	if temporarySendingChange.Hours() > 2 {
		err := app.worker.AttemptsProcessing(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to execute attempts processing: %w", err))
			return
		}

		// Create new pin code and send to user
		err = app.worker.SendingProcessing(userId, email)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to send pin to user: %w", err))
			return
		}

		err = app.pins.DeleteUser(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
			return
		}

		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin has expired"))
		return
	}

	// Compare pin
	err = processing.PinComparing(dbPin, pin)
	if err != nil {
		err := app.worker.AttemptsProcessing(userId)
		if err != nil {
			app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to execute attempts processing: %w", err))
			return
		}

		app.fail(c, http.StatusUnauthorized, fmt.Errorf("pin is not correct"))
		return
	}

	err = app.worker.DeleteUserFromPinsAndAttemptsTables(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to execute deleting user from pins and attempts tables: %w", err))
		return
	}

	// Success end
	app.success(c, "confirmed")
}

// Struct of request data
var resendPinBody struct {
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func (app app) resendPin(c *gin.Context) {
	// Bind struct values
	err := c.Bind(&resendPinBody)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to bind resend pin body: %w", err))
		return
	}

	var (
		userId = resendPinBody.UserId
		email  = resendPinBody.Email
	)

	isUserExistInBlocksTable, err := app.blocks.CheckIfUserExist(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to check if user exist in blocks table: %w", err))
	}

	if isUserExistInBlocksTable != 0 {
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("user is blocked"))
		return
	}

	_, temporarySendingChange, err := app.worker.GetUserDataAndCalcTimeDiff(userId)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to execute getting user data and calculate time difference: %w", err))
		return
	}

	// Check when was sent last email
	if temporarySendingChange.Minutes() < 1 {
		timeUntilNewSend := 60 - temporarySendingChange.Seconds()
		app.fail(c, http.StatusUnauthorized, fmt.Errorf("time has not passed since the last sending: %d", int(timeUntilNewSend)))
		return
	}

	// Create new pin code and send to user
	err = app.worker.SendingProcessing(userId, email)
	if err != nil {
		app.fail(c, http.StatusBadRequest, fmt.Errorf("failed to send pin to user: %w", err))
		return
	}

	// update sent_at

	app.success(c, "new pin was sent")
}
