package app

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-app/config"
	"wallet-app/internal/service/wallet"
	"wallet-app/internal/storage/repository"
	"wallet-app/internal/transport/handler"
	"wallet-app/internal/transport/router"
)

func Run() {
	cfg := config.GetConfig()
	dsn := cfg.GetDBConnStr()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	poolConfig.MaxConns = 100
	poolConfig.MinConns = 10
	poolConfig.MaxConnLifetime = 5 * time.Minute
	poolConfig.HealthCheckPeriod = 30 * time.Second

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer dbPool.Close()

	if err = dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	db, err := repository.NewDatabase(dsn)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Conn.Close(context.Background())

	err = repository.RunMigrations(db.Conn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	repo := repository.NewWalletRepository(dbPool)
	walletService := service.NewWalletService(repo)
	walletHandler := handler.NewWalletHandler(walletService)

	app := gin.Default()

	router.SetupRoutes(app, walletHandler)

	port := cfg.AppPort

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server is starting at port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

}
