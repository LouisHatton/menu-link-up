package main

import (
	"context"

	"net/http"

	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"

	"github.com/LouisHatton/menu-link-up/internal/api"
	apiMiddleware "github.com/LouisHatton/menu-link-up/internal/api/middleware"
	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	"github.com/LouisHatton/menu-link-up/internal/config/environment"
	filesStore "github.com/LouisHatton/menu-link-up/internal/files/store"
	filesStoreReader "github.com/LouisHatton/menu-link-up/internal/files/store/reader"
	filesStoreWriter "github.com/LouisHatton/menu-link-up/internal/files/store/writer"
	projectsStore "github.com/LouisHatton/menu-link-up/internal/projects/store"
	projectsStoreReader "github.com/LouisHatton/menu-link-up/internal/projects/store/reader"
	projectsStoreWriter "github.com/LouisHatton/menu-link-up/internal/projects/store/writer"
	"github.com/caarlos0/env/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type config struct {
	appconfig.Environment
	appconfig.Server
	appconfig.Database
}

func main() {

	// --- ENV & Logging
	ctx := context.Background()
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		panic("failed to parse server config env: " + err.Error())
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	if cfg.Environment.CurrentEnv == environment.Production {
		logger, err = zap.NewProduction()
		if err != nil {
			panic("failed to create production logger: " + err.Error())
		}
	}
	defer logger.Sync()

	// --- GCloud
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "insight-wave-dev",
	})
	if err != nil {
		logger.Fatal("error initializing app", zap.Error(err))
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Fatal("error initializing app auth", zap.Error(err))
	}

	store, err := app.Firestore(ctx)
	if err != nil {
		logger.Fatal("error initializing firestore", zap.Error(err))
	}

	// --- Middleware
	authMiddleware, err := apiMiddleware.NewAuth(logger, authClient)
	if err != nil {
		logger.Fatal("error initializing auth middleware", zap.Error(err))
	}

	// --- Projects Store
	projectReader, err := projectsStoreReader.New(logger, cfg.ProjectsCollectionName, store)
	if err != nil {
		logger.Fatal("error initializing projectsStoreReader", zap.Error(err))
	}

	projectsWriter, err := projectsStoreWriter.New(logger, cfg.ProjectsCollectionName, store)
	if err != nil {
		logger.Fatal("error initializing projectsStoreReader", zap.Error(err))
	}
	projectStore := projectsStore.New(projectReader, projectsWriter)

	// --- Files Store
	fileReader, err := filesStoreReader.New(logger, cfg.FilesCollectionName, store)
	if err != nil {
		logger.Fatal("error initializing filesStoreReader", zap.Error(err))
	}

	fileWriter, err := filesStoreWriter.New(logger, cfg.FilesCollectionName, store)
	if err != nil {
		logger.Fatal("error initializing filesStoreWriter", zap.Error(err))
	}
	fileStore := filesStore.New(fileReader, fileWriter)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	api, err := api.New(logger, cfg.Environment.CurrentEnv, authMiddleware, projectStore, fileStore)
	if err != nil {
		logger.Fatal("error initializing api", zap.Error(err))
	}

	err = api.Register(r)
	if err != nil {
		logger.Fatal("error registering api routes", zap.Error(err))
	}

	logger.Info("Webserver started", zap.String("port", cfg.Port), zap.String("env", string(cfg.Environment.CurrentEnv)))
	http.ListenAndServe(":"+cfg.Port, r)
}
