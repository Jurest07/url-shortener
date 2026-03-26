package update

import (
	"errors"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL 	string `json:"url" validate:"required,url"`
	Alias 	string `json:"alias" validate:"required"`
}

type Response struct {
	resp.Response
	Alias  string `json:"alias"`
}

const aliasLength = 6


//go:generate go run github.com/vektra/mockery/v2@v2.53.6 --name=URLUpdater
type URLUpdater interface {
	UpdateURL(urlToSave string, alias string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.6 --name=AliasUpdater
type AliasUpdater interface {
	UpdateAlias(urlToSave string, alias string) error
}

func NewURL(log *slog.Logger, aliasUpdater AliasUpdater) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.NewURL"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validatorErr))
			return
		}

		alias := req.Alias
		err = aliasUpdater.UpdateAlias(req.URL, alias)
		if errors.Is(err, storage.ErrAliasNotFound) {
			log.Info("not found alias to update", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found alias to update"))
			return
		}
		if err != nil {
			log.Error("failed to update alias", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update alias"))
			return
		}

		log.Info("alias updated")
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias: alias,
		})
	}
}

func NewAlias(log *slog.Logger, URLUpdater URLUpdater) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.NewAlias"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validatorErr))
			return
		}

		alias := req.Alias
		err = URLUpdater.UpdateURL(req.URL, alias)
		if errors.Is(err, storage.ErrAliasExists) {
			log.Info("alias already exists", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("alias already exists"))
			return
		}
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("not found url to update", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("not found url to update"))
			return
		}
		if err != nil {
			log.Error("failed to update url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update url"))
			return
		}

		log.Info("url updated")
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias: alias,
		})
	}
}