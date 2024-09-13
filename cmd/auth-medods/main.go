package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
	"github.com/volchok96/auth-medods/internal/database/pgsql"
	"github.com/volchok96/auth-medods/internal/handlers"
)



const (
	tokenTTL        = 30 * time.Minute
	ownKey   string = "volchok96"
	// 2 строки подключения:
	// + на localhost
	// + для развертывания в Docker-контейнерах <docker-compose up --build>
	// connStr         = "postgres://postgres:mypass@localhost:5432/postgres?sslmode=disable"
	connStr         = "postgres://postgres:mypass@db:5432/postgres?sslmode=disable"

)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	storage, err := pgsql.NewDB(connStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to init storage")
		os.Exit(1)
	}
	defer storage.Close()

	routes := handlers.SetupRoutes(storage, ownKey, tokenTTL)

	// Настройка HTTP-сервера
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Msg("server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	<-done
	log.Info().Msg("stopping server")

	// Грейсфул остановка сервера
	if err := srv.Close(); err != nil {
		log.Error().Err(err).Msg("server shutdown failed")
	} else {
		log.Info().Msg("server stopped gracefully")
	}
}
