package main

import (
	"manga-library/internal/config"
	"manga-library/internal/domain"
	"manga-library/internal/handler"
	"manga-library/internal/service"
	"manga-library/internal/storage"
	"manga-library/internal/storage/mongodb"
	"manga-library/pkg/logger"

	"manga-library/internal/server"
	"manga-library/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const configPath = "configs/config.yaml"

// @Title Manga Library API
// @Version 1.0
// @Description API Server for Manga Library App

// @Host localhost:8080
// @BasePath /

// @SecurityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization

func main() {
	cfg := config.GetYAMLConfig(configPath)
	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	adminUser := domain.AdminUser{
		Username: cfg.AppConfig.AdminUser.Login,
		Password: cfg.AppConfig.AdminUser.Password,
	}

	l := logger.NewLogrusLogger(cfg.AppConfig.LogLevel, cfg.IsDebug, cfg.IsProd)
	l.Infoln("logger initializated")

	mongodb := mongodb.NewMongoDB(cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.Username, cfg.DBConfig.Password, cfg.DBConfig.DBName)
	l.Infoln("db connected successfully")

	// TODO: set logger as first argument for layers

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.TokenTTL)
	// storage := storage.NewInMemoryStorage(l)
	storage := storage.NewStorage(mongodb, l)
	service := service.NewService(storage, jwtManager, l, adminUser)
	handler := handler.NewHandler(service, l)

	handlers := handler.InitRoutes()
	server := server.NewServer(l)
	if err := server.Run(cfg.Listen.Port, cfg.Listen.BindIP, handlers); err != nil {
		panic(err)
	}

}
