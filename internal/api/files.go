package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LouisHatton/menu-link-up/internal/api/responses"
	"github.com/LouisHatton/menu-link-up/internal/api/routes"
	internalContext "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/db/query"
	"github.com/LouisHatton/menu-link-up/internal/files"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (api *API) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := internalContext.GetUserFromContext(ctx)
	logger := api.l.With(zap.String("userId", user.Id))
	project, ok := internalContext.GetProjectFromContext(ctx)
	if !ok {
		logger.Error("unable to get project from context")
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}
	logger = logger.With(zap.String("projectId", project.Id))

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(zap.String("fileId", id))

	file, err := api.fileStore.Get(id)
	if err != nil {
		logger.Error("error getting document", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	if file.ProjectId != project.Id {
		logger.Warn("file is not a member of the project")
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	render.JSON(w, r, &file)
}

func (api *API) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := internalContext.GetUserFromContext(ctx)
	project, ok := internalContext.GetProjectFromContext(ctx)
	if !ok {
		api.l.Error("unable to get project from context")
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}
	logger := api.l.With(zap.String("userId", user.Id), zap.String("projectId", project.Id))

	data := files.NewFile{}
	if err := render.Decode(r, &data); err != nil {
		logger.Error("error parsing provided file data", zap.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	id := uuid.New().String()
	urlId := uuid.New().String()
	logger = logger.With(zap.String("fileId", id))
	newFile := files.File{
		Id:        id,
		UrlId:     urlId,
		ProjectId: project.Id,
		Name:      data.Name,
		Tags:      data.Tags,
		Metadata: files.Metadata{
			CreatedBy: user.Id,
			CreatedAt: time.Now(),
		},
		Status: files.Deploying,
	}

	if err := api.fileStore.Set(id, &newFile); err != nil {
		logger.Error("failed to store new file", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Info("new file created")
	render.JSON(w, r, &newFile)
}

func (api *API) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := internalContext.GetUserFromContext(ctx)
	project, ok := internalContext.GetProjectFromContext(ctx)
	if !ok {
		api.l.Error("unable to get project from context")
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}
	logger := api.l.With(zap.String("userId", user.Id), zap.String("projectId", project.Id))

	docs, err := api.fileStore.Many(query.Options{}, query.Where{
		Key:     "projectId",
		Matcher: query.EqualTo,
		Value:   project.Id,
	})
	if err != nil {
		logger.Fatal("failed to fetch files", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Debug("number of docs returned", zap.Int("count", len(*docs)))

	render.JSON(w, r, docs)
}

func (api *API) EditFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := internalContext.GetUserFromContext(ctx)
	project, ok := internalContext.GetProjectFromContext(ctx)
	if !ok {
		api.l.Error("unable to get project from context")
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}
	logger := api.l.With(zap.String("userId", user.Id), zap.String("projectId", project.Id))

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

	file, err := api.fileStore.Get(id)
	if err != nil {
		logger.Error("error getting document", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	if file.ProjectId != project.Id {
		logger.Warn("file is not a member of the project")
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	file.Name = newFile.Name
	file.Tags = newFile.Tags

	if err := api.fileStore.Set(id, file); err != nil {
		logger.Error("failed to store file", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Info("file updated")
	render.JSON(w, r, file)
}

func (api *API) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := internalContext.GetUserFromContext(ctx)
	project, ok := internalContext.GetProjectFromContext(ctx)
	if !ok {
		api.l.Error("unable to get project from context")
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}
	logger := api.l.With(zap.String("userId", user.Id), zap.String("projectId", project.Id))

	id, err := getFileIdFromUrl(r)
	if err != nil {
		logger.Error("unable to get file id from url", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}
	logger = logger.With(zap.String("fileId", id))

	file, err := api.fileStore.Get(id)
	if err != nil {
		logger.Error("error getting document", zap.Error(err))
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	if file.ProjectId != project.Id {
		logger.Warn("file is not a member of the project")
		render.Render(w, r, responses.NotFoundResponse("file"))
		return
	}

	err = api.fileStore.Delete(file.Id)
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
