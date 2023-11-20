package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/LouisHatton/menu-link-up/internal/api/responses"
	"github.com/LouisHatton/menu-link-up/internal/bandwidth"
	internalContext "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/files"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func (api *API) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(log.String("userId", userId))

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", log.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(log.String("fileId", id))

	file, err := api.fileSvc.GetById(ctx, id)
	if err != nil {
		logger.Error("error getting file by id", log.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	render.JSON(w, r, &file)
}

func (api *API) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(log.String("userId", userId))

	data := files.NewFile{}
	if err := render.Decode(r, &data); err != nil {
		logger.Error("error parsing provided file data", log.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	newFile, err := api.fileSvc.Create(ctx, userId, data)
	switch err {
	case nil:
	case files.ErrSlugAlreadyInUse:
		logger.Warn("attempting to create file with invalid slug", log.Error(err))
		render.Render(w, r, &responses.HttpResponse{
			StatusCode: http.StatusConflict,
			StatusText: "File with slug already exists",
		})
		return
	case bandwidth.ErrUploadLimitReached:
		logger.Warn("user has reached file upload limit", log.Error(err), zap.Int("fileSize", data.FileSize))
		render.Render(w, r, &responses.HttpResponse{
			StatusCode: http.StatusTooManyRequests,
			StatusText: "File upload limit reached, contact support",
		})
		return
	default:
		logger.Error("attempting to create file", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, &newFile)
}

func (api *API) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(log.String("userId", userId))

	docs, err := api.fileSvc.GetByUserId(ctx, userId)
	if err != nil {
		logger.Error("failed to fetch files", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, &docs)
}

func (api *API) EditFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(log.String("userId", userId))

	newFile := files.NewFile{}
	if err := render.Decode(r, &newFile); err != nil {
		logger.Error("error parsing provided file data", log.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", log.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(log.String("fileId", id))

	err = api.fileSvc.Edit(ctx, id, newFile)
	msg := "attempting to edit file"
	switch err {
	case nil:
	case files.ErrNotUsersFile:
		logger.Warn(msg, log.Error(err))
		render.Render(w, r, responses.ErrForbidden())
	default:
		logger.Error(msg, log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
	}

	render.Status(r, http.StatusOK)
}

func (api *API) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(log.UserId(userId), log.Context(ctx))

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", log.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(log.String("fileId", id))

	err = api.fileSvc.Delete(ctx, id)
	msg := "attempting to delete file"
	switch err {
	case nil:
	case files.ErrNotUsersFile:
		logger.Warn(msg, log.Error(err))
		render.Render(w, r, responses.ErrForbidden())
	default:
		logger.Error(msg, log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
	}

	render.Status(r, http.StatusOK)
}

func (api *API) GetObjectStoreLinkForFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := api.l.With(log.Context(ctx))
	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", log.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	link, err := api.fileSvc.GetLinkFromSlug(ctx, id)
	switch err {
	case nil:
	case files.ErrFileNotFound:
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	case bandwidth.ErrBytesTransferredLimitReached:
		logger.Info("the requested file's user has reached the bandwidth limit", log.Error(err))
		render.Render(w, r, &responses.HttpResponse{
			StatusCode: http.StatusTooManyRequests,
			StatusText: "Too many requests",
		})
		return
	default:
		logger.Error("attempting to get link from slug", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
	}

	render.JSON(w, r, link)
}

func (api *API) CheckFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internalContext.GetUserIdFromContext(ctx)
	logger := api.l.With(log.UserId(userId), log.Context(ctx))

	data := map[string]string{}
	if err := render.Decode(r, &data); err != nil {
		logger.Error("error parsing provided file data", log.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	slug, ok := data["slug"]
	if !ok {
		logger.Error("slug not provided in request")
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("slug not provided")))
		return
	}

	exists := api.fileSvc.ExistsWithSlug(ctx, slug)
	if exists {
		logger.Info("file with slug already exists", log.String("slug", slug), log.Error(files.ErrSlugAlreadyInUse))
		render.Render(w, r, responses.AlreadyExists("file"))
		return
	}

	render.Status(r, http.StatusOK)
}

func getFileIdFromUrl(r *http.Request) (string, error) {
	if id := chi.URLParam(r, FileIdParam); id != "" {
		return id, nil
	} else {
		return "", fmt.Errorf("url does not contain file id: url: %s", r.URL.String())
	}
}
