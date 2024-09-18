package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"ozon_replic/internal/pkg/config"
	"ozon_replic/internal/pkg/utils/logger"
	"ozon_replic/internal/pkg/utils/logger/sl"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg := config.MustLoad() // TODO : dev-config.yaml -> readme.

	logFile, err := os.OpenFile("products.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("fail open logFile", sl.Err(err))
		return fmt.Errorf("fail open logFile: %w", err)
	}
	defer func() {
		err = errors.Join(err, logFile.Close())
	}()

	log := logger.Set(cfg.Enviroment, logFile)

	log.Info(
		"starting "+cfg.GRPC.ProductsContainerIP,
		slog.String("env", cfg.Enviroment),
		slog.String("addr", fmt.Sprintf("%s:%d", cfg.GRPC.ProductsContainerIP, cfg.GRPC.ProductsPort)),
		slog.String("log_file_path", cfg.LogFilePath),
	)
	log.Debug("debug messages are enabled")

	//	::::: DB CREATE

	return err
}
