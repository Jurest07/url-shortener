package delete

import (
	"errors"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.6 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrAliasNotFound) {
			log.Info("alias not found", slog.String("alias", alias))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("alias not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		log.Info("url deleted")
		render.JSON(w, r, resp.OK())
	}
}