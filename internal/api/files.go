package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LouisHatton/menu-link-up/internal/api/responses"
	"github.com/LouisHatton/menu-link-up/internal/api/routes"
	internalContext "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/files"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (api *API) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(zap.String("userId", userId))

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(zap.String("fileId", id))

	file, err := api.fileStore.GetById(ctx, id)
	if err != nil {
		logger.Error("error getting document", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	render.JSON(w, r, &file)
}

func (api *API) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(zap.String("userId", userId))

	data := files.NewFile{}
	if err := render.Decode(r, &data); err != nil {
		logger.Error("error parsing provided file data", zap.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	id := uuid.New().String()
	logger = logger.With(zap.String("fileId", id))
	newFile := files.File{
		ID:        id,
		Slug:      data.Slug,
		Name:      data.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    userId,
	}

	if err := api.fileStore.Create(ctx, &newFile); err != nil {
		logger.Error("failed to store new file", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Info("new file created")
	render.JSON(w, r, &newFile)
}

func (api *API) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(zap.String("userId", userId))

	logger.Debug("getting files")

	docs, err := api.fileStore.GetByUserId(ctx, userId)
	if err != nil {
		logger.Error("failed to fetch files", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	render.JSON(w, r, &docs)
}

func (api *API) EditFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)

	logger := api.l.With(zap.String("userId", userId))

	newFile := files.NewFile{}
	if err := render.Decode(r, &newFile); err != nil {
		logger.Error("error parsing provided file data", zap.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(zap.String("fileId", id))

	file, err := api.fileStore.GetById(ctx, id)
	if err != nil {
		logger.Error("error getting document", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	file.Name = newFile.Name

	if err := api.fileStore.Update(ctx, file); err != nil {
		logger.Error("failed to store file", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Info("file updated")
	render.JSON(w, r, file)
}

func (api *API) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(zap.String("userId", userId))

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(zap.String("fileId", id))

	file, err := api.fileStore.GetById(ctx, id)
	if err != nil {
		logger.Error("error getting document", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	err = api.fileStore.DeleteById(ctx, file.ID)
	if err != nil {
		logger.Error("error deleting file", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Info("file deleted")
	render.Status(r, 200)
}

func getFileIdFromUrl(r *http.Request) (string, error) {
	if id := chi.URLParam(r, routes.FileIdParam); id != "" {
		return id, nil
	} else {
		return "", fmt.Errorf("url does not contain file id: url: %s", r.URL.String())
	}
}
