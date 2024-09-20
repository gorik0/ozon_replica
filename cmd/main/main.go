package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	_ "ozon_replic/docs"
	http8 "ozon_replic/internal/pkg/address/delivery/http"
	repo4 "ozon_replic/internal/pkg/address/repo"
	usecase5 "ozon_replic/internal/pkg/address/usecase"
	"ozon_replic/internal/pkg/auth/delivery/grpc/gen"
	http2 "ozon_replic/internal/pkg/auth/delivery/http"
	http4 "ozon_replic/internal/pkg/cart/delivery/http"
	cartRepo "ozon_replic/internal/pkg/cart/repo"
	usecase2 "ozon_replic/internal/pkg/cart/usecase"
	http6 "ozon_replic/internal/pkg/category/delivery/http"
	repo2 "ozon_replic/internal/pkg/category/repo"
	usecase3 "ozon_replic/internal/pkg/category/usecase"
	http10 "ozon_replic/internal/pkg/comments/delivery/http"
	repo7 "ozon_replic/internal/pkg/comments/repo"
	usecase7 "ozon_replic/internal/pkg/comments/usecase"
	"ozon_replic/internal/pkg/config"
	"ozon_replic/internal/pkg/middleware"
	"ozon_replic/internal/pkg/middleware/authmw"
	"ozon_replic/internal/pkg/middleware/csrfmw"
	gen3 "ozon_replic/internal/pkg/order/delivery/grpc/gen"
	http7 "ozon_replic/internal/pkg/order/delivery/http"
	repo3 "ozon_replic/internal/pkg/order/repo"
	usecase4 "ozon_replic/internal/pkg/order/usecase"
	gen2 "ozon_replic/internal/pkg/products/delivery/grpc/gen"
	http5 "ozon_replic/internal/pkg/products/delivery/http"
	repo8 "ozon_replic/internal/pkg/products/repo"
	http3 "ozon_replic/internal/pkg/profile/delivery/http"
	"ozon_replic/internal/pkg/profile/repo"
	"ozon_replic/internal/pkg/profile/usecase"
	http12 "ozon_replic/internal/pkg/promo/delivery/http"
	repo5 "ozon_replic/internal/pkg/promo/repo"
	usecase9 "ozon_replic/internal/pkg/promo/usecase"
	http11 "ozon_replic/internal/pkg/recommendations/delivery/http"
	repo9 "ozon_replic/internal/pkg/recommendations/repo"
	usecase8 "ozon_replic/internal/pkg/recommendations/usecase"
	http9 "ozon_replic/internal/pkg/search/delivery/http"
	repo6 "ozon_replic/internal/pkg/search/repo"
	usecase6 "ozon_replic/internal/pkg/search/usecase"
	"ozon_replic/internal/pkg/utils/jwter"
	"ozon_replic/internal/pkg/utils/logger"
	"ozon_replic/internal/pkg/utils/logger/sl"
	"syscall"
	"time"
)

// @title ZuZu Backend API
// @description API server for ZuZu.

// @contact.name Dima
// @contact.url http://t.me/belozerovmsk

// @securityDefinitions	AuthKey
// @in					header
// @name				Authorization
func main() {
	os.Setenv("AUTH_JWT_SECRET_KEY", "a")
	os.Setenv("CSRF_JWT_SECRET_KEY", "a")
	os.Setenv("POSTGRES_DB", "postgres")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("POSTGRES_PASSWORD", "gorik")
	os.Setenv("POSTGRES_USER", "goirk")

	os.Setenv("GRPC_AUTH_CONTAINER_IP", "localhost")
	os.Setenv("GRPC_ORDER_CONTAINER_IP", "localhost")
	os.Setenv("GRPC_PRODUCTS_CONTAINER_IP", "localhost")

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
	fmt.Printf("Starting auth service on port %d\n", cfg.GRPC.AuthPort)
	defer authConn.Close()

	orderConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.GRPC.OrderContainerIP, cfg.GRPC.OrderPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("fail grpc.Dial order", sl.Err(err))
		err = fmt.Errorf("error happened in grpc.Dial order: %w", err)

		return err
	}
	fmt.Printf("Starting order service on port %d\n", cfg.GRPC.OrderPort)

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
	profileRepo := repo.NewProfileRepo(db)
	profileUsecase := usecase.NewProfileUsecase(profileRepo, cfg)
	profileHandler := http3.NewProfileHandler(log, profileUsecase)

	authClient := gen.NewAuthClient(authConn)
	authHandler := http2.NewAuthHandler(authClient, log)
	//
	cartRepo := cartRepo.NewCartRepo(db)
	cartUsecase := usecase2.NewCartUsecase(cartRepo)
	cartHandler := http4.NewCartHandler(log, cartUsecase)

	productsClient := gen2.NewProductsClient(productConn)
	productsRepo := repo8.NewProductsRepo(db)
	productsHandler := http5.NewProductsHandler(productsClient, log)

	searchRepo := repo6.NewSearchRepo(db)
	searchUsecase := usecase6.NewSearchUsecase(searchRepo, productsRepo)
	searchHandler := http9.NewSearchHandler(log, searchUsecase)
	////
	categoryRepo := repo2.NewCategoryRepo(db)
	categoryUsecase := usecase3.NewCategoryUsecase(categoryRepo)
	categoryHandler := http6.NewCategoryHandler(log, categoryUsecase)
	////
	addressRepo := repo4.NewAddressRepo(db)
	addressUsecase := usecase5.NewAddressUsecase(addressRepo)
	addressHandler := http8.NewAddressHandler(log, addressUsecase)
	////
	promoRepo := repo5.NewPromoRepo(db)
	promoUsecase := usecase9.NewPromoUsecase(promoRepo)
	promoHandler := http12.NewPromoHandler(log, promoUsecase)
	////
	////
	orderRepo := repo3.NewOrderRepo(db)
	orderUsecase := usecase4.NewOrderUsecase(orderRepo, cartRepo, addressRepo, promoRepo)
	orderClient := gen3.NewOrderClient(orderConn)
	orderHandler := http7.NewOrderHandler(orderClient, log, orderUsecase)
	////
	commentsRepo := repo7.NewCommentsRepo(db)
	commentsUsecase := usecase7.NewCommentsUsecase(commentsRepo)
	commentsHandler := http10.NewCommentsHandler(log, commentsUsecase)
	//
	recRepo := repo9.NewRecommendationsRepo(db)
	recUsecase := usecase8.NewRecommendationsUsecase(recRepo)
	recHandler := http11.NewRecommendationsHandler(log, recUsecase)
	////
	//hub := NewHub(orderRepo)
	//notificationsRepo := NewNotificationsRepo(db)
	//notificationsUsecase := NewNotificationsUsecase(notificationsRepo)
	//notificationsHandler := NewNotificationsHandler(hub, notificationsUsecase, log)

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
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// ::::::; endPOINTS ;::::::\\\\\\\

	authMW := authmw.New(log, jwter.New(cfg.AuthJWT))
	csrfMW := csrfmw.New(log, jwter.New(cfg.CSRFJWT))

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", csrfMW(http.HandlerFunc(authHandler.SignUp))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		auth.Handle("/signin", csrfMW(http.HandlerFunc(authHandler.SignIn))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		auth.Handle("/logout", authMW(http.HandlerFunc(authHandler.LogOut))).
			Methods(http.MethodGet, http.MethodOptions)

		auth.Handle("/check_auth", authMW(http.HandlerFunc(authHandler.CheckAuth))).
			Methods(http.MethodGet, http.MethodOptions)
	}

	profile := r.PathPrefix("/profile").Subrouter()
	{
		profile.HandleFunc("/{id:[0-9a-fA-F-]+}", profileHandler.GetProfile).
			Methods(http.MethodGet, http.MethodOptions)

		profile.Handle("/update-photo", authMW(csrfMW(http.HandlerFunc(profileHandler.UpdatePhoto)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		profile.Handle("/update-data", authMW(csrfMW(http.HandlerFunc(profileHandler.UpdateProfileData)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
	}

	cart := r.PathPrefix("/cart").Subrouter()
	{
		cart.Handle("/update", authMW(http.HandlerFunc(cartHandler.UpdateCart))).
			Methods(http.MethodPost, http.MethodOptions)

		cart.Handle("/summary", authMW(http.HandlerFunc(cartHandler.GetCart))).
			Methods(http.MethodGet, http.MethodOptions)

		cart.Handle("/add_product", authMW(http.HandlerFunc(cartHandler.AddProduct))).
			Methods(http.MethodPost, http.MethodOptions)

		cart.Handle("/delete_product", authMW(http.HandlerFunc(cartHandler.DeleteProduct))).
			Methods(http.MethodDelete, http.MethodOptions)
	}

	products := r.PathPrefix("/products").Subrouter()
	{
		products.HandleFunc("/{id:[0-9a-fA-F-]+}", productsHandler.Product).
			Methods(http.MethodGet, http.MethodOptions)

		products.HandleFunc("/get_all", productsHandler.Products).
			Methods(http.MethodGet, http.MethodOptions)

		products.HandleFunc("/category", productsHandler.Category).
			Methods(http.MethodGet, http.MethodOptions)
	}
	category := r.PathPrefix("/category").Subrouter()
	{
		category.HandleFunc("/get_all", categoryHandler.Categories).
			Methods(http.MethodGet, http.MethodOptions)
	}

	order := r.PathPrefix("/order").Subrouter()
	{
		order.Handle("/create", authMW(csrfMW(http.HandlerFunc(orderHandler.CreateOrder)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		order.Handle("/get_current", authMW(http.HandlerFunc(orderHandler.GetCurrentOrder))).
			Methods(http.MethodGet, http.MethodOptions)

		order.Handle("/get_all", authMW(http.HandlerFunc(orderHandler.GetOrders))).
			Methods(http.MethodGet, http.MethodOptions)
	}

	address := r.PathPrefix("/address").Subrouter()
	{
		address.Handle("/add", authMW(csrfMW(http.HandlerFunc(addressHandler.AddAddress)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		address.Handle("/update", authMW(csrfMW(http.HandlerFunc(addressHandler.UpdateAddress)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		address.Handle("/delete", authMW(csrfMW(http.HandlerFunc(addressHandler.DeleteAddress)))).
			Methods(http.MethodDelete, http.MethodGet, http.MethodOptions)

		address.Handle("/make_current", authMW(csrfMW(http.HandlerFunc(addressHandler.MakeCurrentAddress)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		address.Handle("/get_current", authMW(http.HandlerFunc(addressHandler.GetCurrentAddress))).
			Methods(http.MethodGet, http.MethodOptions)

		address.Handle("/get_all", authMW(http.HandlerFunc(addressHandler.GetAllAddresses))).
			Methods(http.MethodGet, http.MethodOptions)
	}

	search := r.PathPrefix("/search").Subrouter()
	{
		search.HandleFunc("/", searchHandler.SearchProducts).
			Methods(http.MethodGet, http.MethodOptions)
	}

	comments := r.PathPrefix("/comments").Subrouter()
	{
		comments.Handle("/create", authMW(csrfMW(http.HandlerFunc(commentsHandler.CreateComment)))).
			Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

		comments.HandleFunc("/get_all", commentsHandler.GetProductComments).
			Methods(http.MethodGet, http.MethodOptions)
	}

	promo := r.PathPrefix("/promo").Subrouter()
	{
		promo.Handle("/check", authMW(http.HandlerFunc(promoHandler.CheckPromocode))).
			Methods(http.MethodGet, http.MethodOptions)

		promo.Handle("/use", authMW(http.HandlerFunc(promoHandler.UsePromocode))).
			Methods(http.MethodGet, http.MethodOptions)
	}

	recs := r.PathPrefix("/recommendations").Subrouter()
	{
		recs.Handle("/get_all", authMW(http.HandlerFunc(recHandler.Recommendations))).
			Methods(http.MethodGet, http.MethodOptions)

		recs.Handle("/update", authMW(http.HandlerFunc(recHandler.UpdateUserActivity))).
			Methods(http.MethodPost, http.MethodOptions)

		recs.HandleFunc("/get_anon", recHandler.AnonRecommendations).
			Methods(http.MethodGet, http.MethodOptions)
	}
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
