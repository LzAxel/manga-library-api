package main

import (
	"context"
	"manga-library/internal/config"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/internal/storage/mongodb"
	"manga-library/pkg/hash"
	"manga-library/pkg/logger"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	command := os.Args[1]
	cfg := config.GetEnvConfig()
	log := logger.NewLogrusLogger(cfg.AppConfig.LogLevel, cfg.IsDebug, cfg.IsProd)

	switch command {
	case "createsuperuser":
		mongodb := mongodb.NewMongoDB(cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.Username, cfg.DBConfig.Password, cfg.DBConfig.DBName)
		storage := storage.NewStorage(mongodb, log)

		if len(os.Args) != 3+1 {
			log.Fatalln("Need to pass username and password\nExample: createsuperuser username supesafetypass")
		}

		username := os.Args[2]
		password := os.Args[3]

		if len(password) < 8 {
			log.Fatalln("Password must be longer than 8 characters")
		}

		user := domain.User{
			ID:           uuid.NewString(),
			Username:     username,
			PasswordHash: hash.HashPassword(hash.GenerateSalt(), password),
			IsEditor:     true,
			IsAdmin:      true,
			CreatedAt:    time.Now(),
		}

		if err := storage.Authorization.SignUp(context.Background(), user); err != nil {
			log.Panicln(err)
			return
		}
	default:
		log.Fatalln("Invalid command")
	}

}
