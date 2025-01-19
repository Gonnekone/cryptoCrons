package parser

import (
	"context"
	"errors"
	"github.com/Gonnekone/cryptoCrons/internal/lib/coin_api"
	"github.com/Gonnekone/cryptoCrons/internal/lib/logger/sl"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type CoinPriceSaver interface {
	SaveCoinPrice(ctx context.Context, coinId string, price int, timestamp int64) error
}

type CoinGetter interface {
	GetCoin(ctx context.Context, coin string) (string, error)
}

type CoinSaver interface {
	SaveCoin(ctx context.Context, coin string) (string, error)
}

type CoinProvider interface {
	CoinPriceSaver
	CoinGetter
	CoinSaver
}

type Parser struct {
	logger *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc

	coinsToParse map[string]struct{}
	interval     time.Duration
	coinProvider CoinProvider
}

func New(logger *slog.Logger, interval time.Duration, coinProvider CoinProvider) *Parser {
	return &Parser{
		logger:       logger,
		coinsToParse: make(map[string]struct{}),
		coinProvider: coinProvider,
		interval:     interval,
	}
}

func (p *Parser) AddCoin(coin string) {
	p.coinsToParse[coin] = struct{}{}
}

func (p *Parser) DeleteCoin(coin string) {
	delete(p.coinsToParse, coin)
}

func (p *Parser) Start(ctx context.Context) {
	p.ctx, p.cancel = context.WithCancel(ctx)

	ticker := time.NewTicker(p.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				p.logger.Info("coins", slog.Any("coins", p.coinsToParse))
				for coin := range p.coinsToParse {
					price, timestamp, err := coin_api.GetCoins(p.ctx, coin)
					if err != nil {
						p.logger.Warn("failed to get coin price", sl.Err(err))

						continue
					}

					coinId, err := p.coinProvider.GetCoin(p.ctx, coin)

					if err != nil && !errors.Is(err, pgx.ErrNoRows) {
						p.logger.Warn("failed to get coin id", sl.Err(err))

						continue

					} else if errors.Is(err, pgx.ErrNoRows) {
						coinId, err = p.coinProvider.SaveCoin(p.ctx, coin)
						if err != nil {
							p.logger.Warn("failed to save coin", sl.Err(err))

							continue
						}
					}

					err = p.coinProvider.SaveCoinPrice(p.ctx, coinId, price, timestamp)
					if err != nil {
						p.logger.Warn("failed to save coin price", sl.Err(err))

						continue
					}
				}

			case <-p.ctx.Done():
				return
			}
		}
	}()
}

func (p *Parser) Stop() {
	p.cancel()
}
