package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"ozon_replic/internal/pkg/auth/delivery/grpc/gen"
	http2 "ozon_replic/internal/pkg/auth/delivery/http"
	http3 "ozon_replic/internal/pkg/cart/delivery/http"
	cartRepo "ozon_replic/internal/pkg/cart/repo"
	"ozon_replic/internal/pkg/cart/usecase"
	"ozon_replic/internal/pkg/config"
	"ozon_replic/internal/pkg/middleware"
	gen2 "ozon_replic/internal/pkg/products/delivery/grpc/gen"
	http4 "ozon_replic/internal/pkg/products/delivery/http"
	"ozon_replic/internal/pkg/products/repo"
	profileHandler "ozon_replic/internal/pkg/profile/delivery/http"
	profileRepo "ozon_replic/internal/pkg/profile/repo"
	profileUsecase "ozon_replic/internal/pkg/profile/usecase"
	"ozon_replic/internal/pkg/utils/logger"
	"ozon_replic/internal/pkg/utils/logger/sl"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (err error) {
	cfg := config.MustLoad() // TODO : dev-config.yaml -> readme.

	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("fail open logFile", sl.Err(err))
		err = fmt.Errorf("fail open logFile: %w", err)

		return err
	}
	defer func() {
		err = errors.Join(err, logFile.Close())
	}()

	log := logger.Set(cfg.Enviroment, logFile)

	log.Info(
		"starting zuzu-main",
		slog.String("env", cfg.Enviroment),
		slog.String("addr", cfg.Address),
		slog.String("log_file_path", cfg.LogFilePath),
		slog.String("photos_file_path", cfg.PhotosFilePath),
	)
	log.Debug("debug messages are enabled")

	//:::::DB:::::

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

	//:::::DB:::::

	// :::: -.-.-.-.-.-.-. MAKE CONNECT FOR GRPC (auth, order, product)

	authConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.GRPC.AuthContainerIP, cfg.GRPC.AuthPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("fail grpc.Dial auth", sl.Err(err))
		err = fmt.Errorf("error happened in grpc.Dial auth: %w", err)

		return err
	}
	defer authConn.Close()

	orderConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.GRPC.OrderContainerIP, cfg.GRPC.OrderPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("fail grpc.Dial order", sl.Err(err))
		err = fmt.Errorf("error happened in grpc.Dial order: %w", err)

		return err
	}
	defer orderConn.Close()

	productConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.GRPC.ProductsContainerIP, cfg.GRPC.ProductsPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("fail grpc.Dial product", sl.Err(err))
		err = fmt.Errorf("error happened in grpc.Dial product: %w", err)

		return err
	}
	defer productConn.Close()
	// :::: -.-.-.-.-.-.-. MAKE CONNECT FOR GRPC (auth, order, product)

	//::::::     REPO : USECASE : HANDLER : grpcCLIENT ::::::::  \\\\\\\\

	//  ----profile
	profileRepo := profileRepo.NewProfileRepo(db)
	profileUsecase := profileUsecase.NewProfileUsecase(profileRepo, cfg)
	profileHandler := profileHandler.NewProfileHandler(log, profileUsecase)

	authClient := gen.NewAuthClient(authConn)
	authHandler := http2.NewAuthHandler(authClient, log)

	cartRepo := cartRepo.NewCartRepo(db)
	cartUsecase := usecase.NewCartUsecase(cartRepo)
	cartHandler := http3.NewCartHandler(log, cartUsecase)

	productsClient := gen2.NewProductsClient(productConn)
	productsRepo := repo.NewProductsRepo(db)
	productsHandler := http4.NewProductsHandler(productsClient, log)

	//searchRepo := searchRepo.NewSearchRepo(db)
	//searchUsecase := searchUsecase.NewSearchUsecase(searchRepo, productsRepo)
	//searchHandler := searchHandler.NewSearchHandler(log, searchUsecase)
	//
	//categoryRepo := categoryRepo.NewCategoryRepo(db)
	//categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepo)
	//categoryHandler := categoryHandler.NewCategoryHandler(log, categoryUsecase)
	//
	//addressRepo := addressRepo.NewAddressRepo(db)
	//addressUsecase := addressUsecase.NewAddressUsecase(addressRepo)
	//addressHandler := addressHandler.NewAddressHandler(log, addressUsecase)
	//
	//promoRepo := promoRepo.NewPromoRepo(db)
	//promoUsecase := promoUsecase.NewPromoUsecase(promoRepo)
	//promoHandler := promoHandler.NewPromoHandler(log, promoUsecase)
	//
	//orderRepo := orderRepo.NewOrderRepo(db)
	//
	//orderUsecase := orderUsecase.NewOrderUsecase(orderRepo, cartRepo, addressRepo, promoRepo)
	//orderClient := orderGrpc.NewOrderClient(orderConn)
	//orderHandler := orderHandler.NewOrderHandler(orderClient, log, orderUsecase)
	//
	//commentsRepo := commentsRepo.NewCommentsRepo(db)
	//commentsUsecase := commentsUsecase.NewCommentsUsecase(commentsRepo)
	//commentsHandler := commentsHandler.NewCommentsHandler(log, commentsUsecase)
	//
	//recRepo := recRepo.NewRecommendationsRepo(db)
	//recUsecase := recUsecase.NewRecommendationsUsecase(recRepo)
	//recHandler := recHandler.NewRecommendationsHandler(log, recUsecase)
	//
	//hub := clientHub.NewHub(orderRepo)
	//notificationsRepo := notificationsRepo.NewNotificationsRepo(db)
	//notificationsUsecase := notificationsUsecase.NewNotificationsUsecase(notificationsRepo)
	//notificationsHandler := notificationsHandler.NewNotificationsHandler(hub, notificationsUsecase, log)

	//::::::     REPO : USECASE : HANDLER : grpcCLIENT   ::::::: \\\\\\\

	// ::::: init ROUTER ::::\\\\\

	r := mux.NewRouter().PathPrefix("/api").Subrouter() // ::::: init ROUTER ::::\\\\\

	r.Use(middleware.Recover(log), middleware.CORSMiddleware)
	//r.Use(middleware.Recover(log), middleware.CORSMiddleware, logmw.New(mt, log))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	//r.PathPrefix("/metrics").Handler(promhttp.Handler())
	//
	//r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
	//	httpSwagger.DeepLinking(true),
	//	httpSwagger.DocExpansion("none"),
	//	httpSwagger.DomID("swagger-ui"),
	//)).Methods(http.MethodGet)

	// ::::::; endPOINTS ;::::::\\\\\\\

	// ::::::; endPOINTS ;::::::\\\\\\\

	// ::::::; make SERVER;::::::\\\\\\\

	http.Handle("/", r)

	srv := http.Server{
		Handler:           r,
		Addr:              cfg.Address,
		ReadTimeout:       cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	quit := make(chan os.Signal, 1)
	// SIGINT = ctrl+c; SIGTERM = kill; Interrupt = аппаратное прерывание, в Windows даст ошибку
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//go hub.Run(context.Background())

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("listen and serve returned err: ", sl.Err(err))
		}
	}()

	log.Info("server started")
	sig := <-quit
	log.Debug("handle quit chanel: ", slog.Any("os.Signal", sig.String()))
	log.Info("server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown returned an err: ", sl.Err(err))
		err = fmt.Errorf("error happened in srv.Shutdown: %w", err)

		return err
	}

	log.Info("server stopped")

	return nil
	// ::::::; make SERVER;::::::\\\\\\\
	return err
}
