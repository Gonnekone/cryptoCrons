package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) SaveCoinPrice(ctx context.Context, coinId string, price int, timestamp int64) error {
	const op = "storage.postgres.SaveCoinPrice"

	_, err := s.db.Exec(ctx, "INSERT INTO prices (id, price, ts) VALUES ($1, $2, $3)", coinId, price, timestamp)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetCoinPrice(ctx context.Context, coin string, timestamp int64) (int, int, error) {
	const op = "storage.postgres.GetCoinPrice"

	query := `
		SELECT price, ts
		FROM prices
		WHERE id = $1
		ORDER BY ABS(ts - $2)
		LIMIT 1;
	`

	var price, ts int
	err := s.db.QueryRow(ctx, query, coin, timestamp).Scan(&price, &ts)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return price, ts, nil
}

func (s *Storage) GetCoin(ctx context.Context, coin string) (string, error) {
	const op = "storage.postgres.GetCoin"

	row := s.db.QueryRow(ctx, "SELECT id FROM coins WHERE name = $1", coin)

	var coinId string
	err := row.Scan(&coinId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return coinId, nil
}

func (s *Storage) SaveCoin(ctx context.Context, coin string) (string, error) {
	const op = "storage.postgres.SaveCoin"

	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(ctx, "INSERT INTO coins(id, name) VALUES ($1, $2)", id.String(), coin)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (s *Storage) Close() {
	s.db.Close()
}
