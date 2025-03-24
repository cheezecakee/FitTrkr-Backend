package handler

import (
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/cheezecakee/FitLogr/internal/auth"
	"github.com/cheezecakee/FitLogr/internal/context"
	"github.com/cheezecakee/FitLogr/internal/services"
	"github.com/cheezecakee/FitLogr/pkg/helper"
)

type Config struct {
	DB         *services.Queries
	JWTManager *auth.JWTManager
	Redis      *redis.Client
	ContextKey map[string]*context.ContextKey
	Helper     *helper.Helper
	Logger     *Logger
	APIRoute   string
}

func NewConfig(dbQueries *services.Queries, jwtSecret []byte, redisAddr string) *Config {
	logger := NewLogger()
	jwtManager := auth.NewJWTManager(jwtSecret, time.Hour)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	helper := helper.NewHelper(logger.InfoLog, logger.ErrorLog)

	return &Config{
		DB:         dbQueries,
		JWTManager: jwtManager,
		Redis:      redisClient,
		ContextKey: map[string]*context.ContextKey{
			"userIDKey": context.UserIDKey,
		},
		Helper: helper,
		Logger: logger,
	}
}

type Logger struct {
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	RequestLog *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		InfoLog:    log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:   log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		RequestLog: log.New(os.Stdout, "REQUEST\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
