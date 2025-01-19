package add

import (
	"errors"
	resp "github.com/Gonnekone/cryptoCrons/internal/lib/api/response"
	"github.com/Gonnekone/cryptoCrons/internal/lib/logger/sl"
	"github.com/Gonnekone/cryptoCrons/internal/parser"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

var (
	ErrCoinNotProvided = errors.New("coin is not provided")
)

type Request struct {
	Coin string `json:"coin"`
}

func New(log *slog.Logger, parser *parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.add.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		if req.Coin == "" {
			log.Error("invalid request", sl.Err(ErrCoinNotProvided))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(resp.InvalidRequest))

			return
		}

		parser.AddCoin(req.Coin)

		render.JSON(w, r, resp.OK())
	}
}
