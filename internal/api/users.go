package api

import (
	"net/http"

	"github.com/LouisHatton/menu-link-up/internal/api/responses"
	internal_context "github.com/LouisHatton/menu-link-up/internal/context"
	"github.com/LouisHatton/menu-link-up/internal/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (api *API) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internal_context.GetUserIdFromContext(ctx)
	logger := api.l.With(log.UserId(userId), log.Context(ctx))

	requestedId := chi.URLParam(r, UserIdParam)
	if requestedId != userId {
		logger.Warn("user attempting to request someone else", log.String("requestedUserId", requestedId))
		render.Render(w, r, responses.ErrForbidden())
		return
	}

	user, err := api.userSvc.GetById(ctx, requestedId)
	if err != nil {
		logger.Error("attempting to fetch user", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, user)
}

func (api *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internal_context.GetUserIdFromContext(ctx)
	logger := api.l.With(log.UserId(userId), log.Context(ctx))

	requestedId := chi.URLParam(r, UserIdParam)
	logger = logger.With(log.RequestedId(requestedId))
	if requestedId != userId {
		logger.Warn("user attempting to delete someone else")
		render.Render(w, r, responses.ErrForbidden())
		return
	}

	err := api.fileSvc.DeleteByUserId(ctx, requestedId)
	if err != nil {
		logger.Error("attempting to delete users files", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	err = api.userSvc.DeleteById(ctx, requestedId)
	if err != nil {
		logger.Error("attempting to delete user", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	logger.Info("user deleted")
	render.Status(r, http.StatusOK)
}

func (api *API) GetUserBilling(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internal_context.GetUserIdFromContext(ctx)
	logger := api.l.With(log.UserId(userId), log.Context(ctx))

	requestedId := chi.URLParam(r, UserIdParam)
	logger = logger.With(log.RequestedId(requestedId))
	if requestedId != userId {
		logger.Warn("user attempting to get someone else")
		render.Render(w, r, responses.ErrForbidden())
		return
	}

	billing, err := api.userSvc.GetBilling(ctx, userId)
	if err != nil {
		logger.Error("attempting get user billing", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, billing)
}

func (api *API) UpdateUserBilling(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := internal_context.GetUserIdFromContext(ctx)
	logger := api.l.With(log.UserId(userId), log.Context(ctx))

	requestedId := chi.URLParam(r, UserIdParam)
	logger = logger.With(log.RequestedId(requestedId))
	if requestedId != userId {
		logger.Warn("user attempting to update billing for someone else")
		render.Render(w, r, responses.ErrForbidden())
		return
	}

	url, err := api.userSvc.UpdateBillingLink(ctx, userId)
	if err != nil {
		logger.Error("attempting update user billing", log.Error(err))
		render.Render(w, r, responses.ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, url)
}
