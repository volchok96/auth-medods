package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
	"github.com/volchok96/auth-medods/internal/database/pgsql"
	"github.com/volchok96/auth-medods/internal/handlers"
)


func loadEnvFile() {
	env := os.Getenv("APP_ENV")
	log.Info().Str("APP_ENV", env).Msg("Environment variable APP_ENV detected") // Отладка переменной

	if env == "" {
		env = "local" // По умолчанию использовать локальную среду
	}

	var err error
	switch env {
	case "docker":
		err = godotenv.Load("../../.env.docker") // Загрузка файла для Docker
	default:
		err = godotenv.Load("../../.env") // Загрузка локального файла
	}

	if err != nil {
		log.Fatal().Msgf("Error loading .env file for environment: %s", env)
	}
}

func main() {
	loadEnvFile()

	// Получение переменных окружения
	ownKey := os.Getenv("OWN_KEY")
	tokenTTL, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid token TTL")
	}

	// Подключение к базе данных
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	log.Info().
		Str("user", dbUser).
		Str("host", dbHost).
		Str("port", dbPort).
		Str("database", dbName).
		Msgf("Connecting to database with: %s", connStr)

	// Пример работы с zerolog
	log.Info().
		Str("ownKey", ownKey).
		Dur("tokenTTL", tokenTTL).
		Msg("Configuration loaded successfully")
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
