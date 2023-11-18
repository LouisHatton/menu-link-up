package main

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/LouisHatton/menu-link-up/internal/api"
	api_middleware "github.com/LouisHatton/menu-link-up/internal/api/middleware"
	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	"github.com/LouisHatton/menu-link-up/internal/db/connection"
	files_repository "github.com/LouisHatton/menu-link-up/internal/files/repository"
	files_service "github.com/LouisHatton/menu-link-up/internal/files/service"
	"github.com/LouisHatton/menu-link-up/internal/log"
	s3_service "github.com/LouisHatton/menu-link-up/internal/objectstore/s3/service"
	users_repository "github.com/LouisHatton/menu-link-up/internal/users/repository"
	users_service "github.com/LouisHatton/menu-link-up/internal/users/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/caarlos0/env/v10"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	appconfig.Environment
	appconfig.Server
	appconfig.Database
	appconfig.AWS_S3
}

func main() {

	// --- ENV & Logging
	ctx := context.Background()
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic("failed to parse server config env: " + err.Error())
	}

	logger := log.NewLogger()

	// --- GCloud
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "insight-wave-dev",
	})
	if err != nil {
		logger.Fatal("error initializing app", log.Error(err))
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Fatal("error initializing app auth", log.Error(err))
	}

	// --- Middleware
	authMiddleware, err := api_middleware.NewAuth(logger, authClient)
	if err != nil {
		logger.Fatal("error initializing auth middleware", log.Error(err))
	}

	db, err := connection.Connect(logger, &cfg.Database)
	if err != nil {
		logger.Fatal("error connecting to database", log.Error(err))
	}
	defer db.Close()

	// --- AWS
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(cfg.AWS_S3.Region)})
	if err != nil {
		logger.Fatal("error creating aws session", log.Error(err))
	}

	s3_objectstore, err := s3_service.New(logger, awsSession, cfg.AWS_S3.DefaultBucket)
	if err != nil {
		logger.Fatal("error creating s3 objectstore service", log.Error(err))
	}

	// --- File SVC
	fileSvc, err := files_service.New(logger, files_repository.New(db), s3_objectstore)
	if err != nil {
		logger.Fatal("error initializing files service", log.Error(err))
	}

	// --- User SVC
	userSvc, err := users_service.New(logger, authClient, users_repository.New(db))
	if err != nil {
		logger.Fatal("error initializing users service", log.Error(err))
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	api, err := api.New(logger, cfg.Environment.CurrentEnv, authMiddleware, fileSvc, userSvc)
	if err != nil {
		logger.Fatal("error initializing api", log.Error(err))
	}

	err = api.Register(r)
	if err != nil {
		logger.Fatal("error registering api routes", log.Error(err))
	}

	logger.Info("Webserver started", log.String("port", cfg.Port), log.String("env", string(cfg.Environment.CurrentEnv)))
	http.ListenAndServe(":"+cfg.Port, r)
}
