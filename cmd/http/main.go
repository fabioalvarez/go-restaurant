package main

import (
	"context"
	"fmt"
	ahttp "go-restaurant/internal/auth/adapter/handler/http"
	"go-restaurant/internal/auth/adapter/paseto"
	"go-restaurant/internal/common/adapter/handler/http"
	"go-restaurant/internal/common/adapter/storage/postgres"

	aservice "go-restaurant/internal/auth/service"

	ohttp "go-restaurant/internal/order/adapter/handler/http"
	orepository "go-restaurant/internal/order/adapter/storage/postgres"
	oservice "go-restaurant/internal/order/service"

	payhttp "go-restaurant/internal/payment/adapter/handler/http"
	payrepository "go-restaurant/internal/payment/adapter/storage/postgres"
	payservice "go-restaurant/internal/payment/service"

	chttp "go-restaurant/internal/category/adapter/handler/http"
	crepository "go-restaurant/internal/category/adapter/storage/postgres"
	cservice "go-restaurant/internal/category/service"

	phttp "go-restaurant/internal/product/adapter/handler/http"
	prepository "go-restaurant/internal/product/adapter/storage/postgres"
	pservice "go-restaurant/internal/product/service"

	"go-restaurant/internal/common/adapter/config"
	"go-restaurant/internal/common/adapter/logger"
	"go-restaurant/internal/common/adapter/storage/redis"
	uhttp "go-restaurant/internal/user/adapter/handler/http"
	urepository "go-restaurant/internal/user/adapter/storage/postgres"
	uservice "go-restaurant/internal/user/service"
	"log/slog"
	"os"

	_ "github.com/bagashiz/go-pos/docs"
)

// @title						Go POS (Point of Sale) API
// @version					1.0
// @description				This is a simple RESTful Point of Sale (POS) Service API written in Go using Gin web framework, PostgreSQL database, and Redis cache.
//
// @contact.name				Bagas Hizbullah
// @contact.url				https://github.com/bagashiz/go-pos
// @contact.email				bagash.office@simplelogin.com
//
// @license.name				MIT
// @license.url				https://github.com/bagashiz/go-pos/blob/main/LICENSE
//
// @host						gopos.bagashiz.me
// @BasePath					/v1
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	// Load environment variables
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger.Set(config.App)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Successfully connected to the database", "db", config.DB.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")

	// Init cache service
	cache, err := redis.New(ctx, config.Redis)
	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	slog.Info("Successfully connected to the cache server")

	// Init token service
	token, err := paseto.New(config.Token)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	// Dependency injection
	// User
	userRepo := urepository.NewUserRepository(db)
	userService := uservice.NewUserService(userRepo, cache)
	userHandler := uhttp.NewUserHandler(userService)

	// Auth
	authService := aservice.NewAuthService(userRepo, token)
	authHandler := ahttp.NewAuthHandler(authService)

	// Payment
	paymentRepo := payrepository.NewPaymentRepository(db)
	paymentService := payservice.NewPaymentService(paymentRepo, cache)
	paymentHandler := payhttp.NewPaymentHandler(paymentService)

	// Category
	categoryRepo := crepository.NewCategoryRepository(db)
	categoryService := cservice.NewCategoryService(categoryRepo, cache)
	categoryHandler := chttp.NewCategoryHandler(categoryService)

	// Product
	productRepo := prepository.NewProductRepository(db)
	productService := pservice.NewProductService(productRepo, categoryRepo, cache)
	productHandler := phttp.NewProductHandler(productService)

	// Order
	orderRepo := orepository.NewOrderRepository(db)
	orderService := oservice.NewOrderService(orderRepo, productRepo, categoryRepo, userRepo, paymentRepo, cache)
	orderHandler := ohttp.NewOrderHandler(orderService)

	// Init router
	router, err := http.NewRouter(
		config.HTTP,
		token,
		*userHandler,
		*authHandler,
		*paymentHandler,
		*categoryHandler,
		*productHandler,
		*orderHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
