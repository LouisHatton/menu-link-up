package main

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/LouisHatton/menu-link-up/internal/api"
	api_middleware "github.com/LouisHatton/menu-link-up/internal/api/middleware"
	"github.com/LouisHatton/menu-link-up/internal/config/appconfig"
	"github.com/LouisHatton/menu-link-up/internal/config/environment"
	"github.com/LouisHatton/menu-link-up/internal/db/connection"
	files_repository "github.com/LouisHatton/menu-link-up/internal/files/repository"
	"github.com/caarlos0/env/v10"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// func main() {
// 	ctx := context.Background()
// 	db, err := db_connection.Connect()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	fileRepo := files_repository.New(db)

// 	id := uuid.NewString()
// 	ogFile := files.File{
// 		ID:        id,
// 		Name:      "Testing Menu",
// 		UserId:    "123",
// 		Slug:      "testing-menu",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 		S3Bucket:  "files",
// 		S3Key:     "testing-menu",
// 	}

// 	err = fileRepo.Create(ctx, &ogFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("added file")

// 	file, err := fileRepo.GetById(ctx, id)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("fetched file: ", file)

// 	file.Name = "Not Testing"
// 	fmt.Println("setting file name: ", file.Name)
// 	err = fileRepo.Update(ctx, file)
// 	if err != nil {
// 		panic(err)
// 	}

// 	newFile, err := fileRepo.GetById(ctx, id)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("fetched file: ", newFile.Name)

// 	err = fileRepo.DeleteById(ctx, id)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("deleted file")

// 	count, err := fileRepo.Count(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("the count of files is: ", count)

// 	files, err := fileRepo.GetByUserId(ctx, "123")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// fmt.Println("fetched all user files: ", files)
// 	fmt.Println("count of user files: ", len(*files)+1)

// }

type config struct {
	appconfig.Environment
	appconfig.Server
	appconfig.Database
}

func main() {

	// --- ENV & Logging
	ctx := context.Background()
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
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

	// --- Middleware
	authMiddleware, err := api_middleware.NewAuth(logger, authClient)
	if err != nil {
		logger.Fatal("error initializing auth middleware", zap.Error(err))
	}

	db, err := connection.Connect(&cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fileRepo := files_repository.New(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	api, err := api.New(logger, cfg.Environment.CurrentEnv, authMiddleware, fileRepo)
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
