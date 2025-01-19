package price

import (
	"context"
	"errors"
	resp "github.com/Gonnekone/cryptoCrons/internal/lib/api/response"
	"github.com/Gonnekone/cryptoCrons/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

var (
	ErrCoinNotProvided      = errors.New("coin is not provided")
	ErrTimestampNotProvided = errors.New("timestamp is not provided")
)

type CoinPriceGetter interface {
	GetCoinPrice(ctx context.Context, coin string, timestamp int64) (int, int, error)
}

type CoinGetter interface {
	GetCoin(ctx context.Context, coin string) (string, error)
}

type CoinProvider interface {
	CoinPriceGetter
	CoinGetter
}

type Request struct {
	Coin      string `json:"coin"`
	Timestamp int64  `json:"timestamp"`
}

type Response struct {
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
	Timestamp int     `json:"timestamp"`
}

func New(log *slog.Logger, coinProvider CoinProvider) http.HandlerFunc {
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

		if req.Timestamp == 0 {
			log.Error("invalid request", sl.Err(ErrTimestampNotProvided))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(resp.InvalidRequest))

			return
		}

		coinId, err := coinProvider.GetCoin(r.Context(), req.Coin)
		if err != nil {
			log.Error("failed to get coin id", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get coin id"))

			return
		}

		price, ts, err := coinProvider.GetCoinPrice(r.Context(), coinId, req.Timestamp)
		if err != nil {
			log.Error("failed to get coin price", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get coin price"))

			return
		}

		render.JSON(w, r, Response{
			Coin:      req.Coin,
			Price:     float64(price) / 100,
			Timestamp: ts,
		})
	}
}
