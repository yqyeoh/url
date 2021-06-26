package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/yqyeoh/url/app"
	"github.com/yqyeoh/url/config"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	db, err = sqlx.Open(conf.DB.DRV, conf.ConnectionString())
	if err != nil {
		log.Fatalf("opening database failed", err)
	}
	log.Info("Database connection opened")

	apiRouter := mux.NewRouter()

	appRepo := app.NewRepo(db)
	appService := appService(log, appRepo)
	app.NewHandler(log, appService, conf).AddRoutes(apiRouter)

	corsInstance := cors.New(cors.Options{
		AllowedOrigins:     conf.CORS.AllowOrigins,
		AllowCredentials:   true,
		AllowedMethods:     []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Origin", "Cache-Control", "Access-Control-Allow-Origin"},
		ExposedHeaders:     []string{"Access-Control-Allow-Origin"},
		MaxAge:             3600,
		OptionsPassthrough: false,
	})

	// start API server
	apiServerAddress := ":" + conf.Port
	apiServer := NewServer(apiServerAddress, corsInstance.Handler(apiRouter))
	go func() {
		log.Infof("[PID=%d] starting API server on %s ...", os.Getpid(), apiServerAddress)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("starting API server error: %v", err)
		}
	}()

	// wait for disruption signal
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info("received SIGINT/SIGTERM and shutting down gracefully ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Info("Server Exited Properly")
}

func NewServer(serverAddr string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	return srv
}
