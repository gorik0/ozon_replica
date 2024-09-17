package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"log/slog"
	"os"
	"ozon_replic/internal/pkg/config"
	"ozon_replic/internal/pkg/utils/logger"
	"ozon_replic/internal/pkg/utils/logger/sl"
)

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {

	//	::: init CONFIG

	cfg := config.MustLoad()

	//	::: init  LOGGER

	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("error opening log file: %v", err)
		return err
	}
	defer func() {
		err = errors.Join(err, logFile.Close())
	}()

	log := logger.Set(cfg.Enviroment, logFile)
	log.Info(
		"starting "+cfg.GRPC.AuthContainerIP,
		slog.String("env", cfg.Enviroment),
		slog.String("addr", fmt.Sprintf("%s:%d", cfg.GRPC.AuthContainerIP, cfg.GRPC.AuthPort)),
		slog.String("log_file_path", cfg.LogFilePath),
	)
	log.Debug("debug messages are enabled")
	//::::::::

	//::::::::
	//	::: inti  DB
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName))
	if err != nil {
		log.Error("fail open postgres", sl.Err(err))
		err = fmt.Errorf("error happened in sql.Open: %w", err)

		return err
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		log.Error("fail ping postgres", sl.Err(err))
		err = fmt.Errorf("error happened in db.Ping: %w", err)

		return err
	}
	//::::::::

	//::::::::
	//	::: inti  REPO
	//	::: inti  USECASE
	//	::: inti  HANDLER
	//	::: init GRPC server

	//	::: go  metric server
	//	::: go  GRPC server
	//	::: gracefull STOP  GRPC

}
