package api

import (
	"net/http"
	"time"

	"github.com/LouisHatton/menu-link-up/internal/api/responses"
	internalContext "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/db/query"
	"github.com/LouisHatton/menu-link-up/internal/projects"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Returns the project for the given projectId in the url if the user is a member of
// the project.
//
// Expects the user to be in the request context
func (api API) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := internalContext.GetUserFromContext(ctx)
	logger := api.l.With(zap.Any("userId", user.Id))

	project, ok := internalContext.GetProjectFromContext(ctx)
	if !ok {
		logger.Error("failed to read project from context")
		render.Render(w, r, responses.NotFoundResponse("project"))
		return
	}

	render.JSON(w, r, &project)
}

func (api API) CreateProject(w http.ResponseWriter, r *http.Request) {
	user := internalContext.GetUserFromContext(r.Context())
	logger := api.l.With(zap.String("userId", user.Id))

	data := projects.NewProject{}
	if err := render.Decode(r, &data); err != nil {
		logger.Error("error parsing provided project data", zap.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	id := uuid.New().String()
	logger = logger.With(zap.String("projectId", id))
	newProject := projects.Project{
		Id:   id,
		Name: data.Name,
		Slug: data.Slug,
		Metadata: projects.Metadata{
			CreatedBy: user.Id,
			CreatedAt: time.Now(),
		},
		Users: []string{user.Id},
	}

	if err := api.projectStore.Set(id, &newProject); err != nil {
		logger.Error("failed to store new project", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	logger.Info("created new project")
	render.Respond(w, r, &newProject)
}

func (api API) CheckProjectSlug(w http.ResponseWriter, r *http.Request) {
	user := internalContext.GetUserFromContext(r.Context())
	logger := api.l.With(zap.String("userId", user.Id))

	data := projects.NewProject{}
	if err := render.Decode(r, &data); err != nil {
		logger.Error("error parsing provided project data", zap.Error(err))
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	_, err := api.projectStore.One(query.Options{}, query.Where{
		Key:     "slug",
		Matcher: query.EqualTo,
		Value:   data.Slug,
	})
	if err != nil {
		render.Respond(w, r, true)
		return
	}

	render.Respond(w, r, false)
}

func (api API) ListProjects(w http.ResponseWriter, r *http.Request) {
	user := internalContext.GetUserFromContext(r.Context())
	logger := api.l.With(zap.Any("userId", user.Id))

	docs, err := api.projectStore.Many(query.Options{}, query.Where{
		Key:     "users",
		Matcher: query.Contains,
		Value:   user.Id,
	})
	if err != nil {
		logger.Fatal("failed to fetch documents", zap.Error(err))
		render.Render(w, r, responses.ErrInternalServerError())
		return
	}

	render.JSON(w, r, docs)
}
