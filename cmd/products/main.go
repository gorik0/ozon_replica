package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	grp "google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"ozon_replic/internal/pkg/config"
	"ozon_replic/internal/pkg/products/delivery/grpc"
	"ozon_replic/internal/pkg/products/delivery/grpc/gen"
	"ozon_replic/internal/pkg/products/repo"
	"ozon_replic/internal/pkg/products/usecase"
	"ozon_replic/internal/pkg/utils/logger"
	"ozon_replic/internal/pkg/utils/logger/sl"
	"syscall"
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
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName))
	if err != nil {
		log.Error("fail open postgres", sl.Err(err))
		return fmt.Errorf("error happened in sql.Open: %w", err)
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		log.Error("fail ping postgres", sl.Err(err))
		return fmt.Errorf("error happened in db.Ping: %w", err)
	}
	//	::::: DB CREATE

	prodRepo := repo.NewProductsRepo(db)
	prodUsecase := usecase.NewProductsUsecase(prodRepo)
	grpcProdHandler := grpc.NewGrpcProductsHandler(prodUsecase, log)

	//grpcMetrics := metrics.NewMetricGRPC(metrics.ServiceAuthName)
	//metricsMw := metricsmw.NewGrpcMiddleware(grpcMetrics)
	gRPCServer := grp.NewServer()

	gen.RegisterProductsServer(gRPCServer, grpcProdHandler)

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8013))
		if err != nil {
			log.Error("listen returned err: ", sl.Err(err))
		}
		log.Info("grpc server started", slog.String("addr", listener.Addr().String()))
		if err := gRPCServer.Serve(listener); err != nil {
			log.Error("serve returned err: ", sl.Err(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	gRPCServer.GracefulStop()
	log.Info("Gracefully stopped")
	return nil
}
