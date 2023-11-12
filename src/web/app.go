package web

import (
	"database/sql"
	"fmt"
	"github.com/KatsayArtemDev/verification/src/database"
	"github.com/KatsayArtemDev/verification/src/initializers"
	"github.com/KatsayArtemDev/verification/src/usecase"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

type app struct {
	pins   database.Pins
	blocks database.Blocks
	worker usecase.Worker
	logger *zap.SugaredLogger
}

func initServer() (app, *sql.DB, error) {
	err := initializers.LoadEnvVariables()
	if err != nil {
		return app{}, nil, fmt.Errorf("failed to init env variables: %w", err)
	}

	logger, err := initializers.LogConfig("./logs", "./logs/verification_")
	if err != nil {
		return app{}, nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	db, err := initializers.ConnectToDb()

	var app = app{
		pins:   database.NewPins(db),
		blocks: database.NewBlocks(db),
		worker: usecase.NewWorker(db),
		logger: logger.Sugar(),
	}

	return app, db, nil
}

func RunServer() error {
	app, db, err := initServer()
	if err != nil {
		return fmt.Errorf("error when starting the app: %w", err)
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatal("error when closing database ", err)
		}
	}(db)

	app.logger.Info("Server was started")

	r := serverRouter(app)

	err = r.Run()
	if err != nil {
		err = fmt.Errorf("error when closing server: %w", err)
		app.logger.Errorln(err.Error())
		return err
	}

	return nil
}
