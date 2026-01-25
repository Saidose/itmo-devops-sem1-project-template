package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"project_sem/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PricesDB struct {
	pool *pgxpool.Pool
}

func NewPricesDb(ctx context.Context, user, pass, name, host string) *PricesDB {
	connString := makePostgresURL(user, pass, name, host)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}
	db := &PricesDB{pool: pool}
	err = db.Migrate(ctx)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func (r *PricesDB) Migrate(ctx context.Context) error {
	const query = `CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    create_date DATE NOT NULL)`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *PricesDB) Close() {
	r.pool.Close()
}

func (r *PricesDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *PricesDB) createTx(ctx context.Context, tx pgx.Tx, p domain.Price) error {
	const q = `
	INSERT INTO prices (name, category, price, create_date)
	VALUES ($1, $2, $3, $4)
	`
	_, err := tx.Exec(ctx, q,
		p.Name,
		p.Category,
		p.Price,
		p.CreateDate,
	)
	return err
}

func (r *PricesDB) GetAll(ctx context.Context) ([]domain.Price, error) {
	const query = `
SELECT id, name, category, price, create_date FROM prices`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prices []domain.Price
	for rows.Next() {
		price := domain.Price{}
		if err := rows.Scan(&price.ID, &price.Name, &price.Category, &price.Price, &price.CreateDate); err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return prices, nil
}

// InsertPrices - вставляет данные в БД транзакцией. Перед коммитом считает статистику.
func (r *PricesDB) InsertPrices(ctx context.Context, prices []domain.Price) (Stats, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Stats{}, err
	}
	defer tx.Rollback(ctx)

	for _, p := range prices {
		if err := r.createTx(ctx, tx, p); err != nil {
			return Stats{}, err
		}
	}

	stats, err := r.getAggregatesTx(ctx, tx, len(prices))
	if err != nil {
		return Stats{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Stats{}, err
	}

	return stats, nil
}

func (r *PricesDB) getAggregatesTx(ctx context.Context, tx pgx.Tx, inserted int) (Stats, error) {
	const q = `
	SELECT 
		COUNT(DISTINCT category),
		COALESCE(SUM(price),0)
	FROM prices
	`

	var s Stats
	s.TotalItems = inserted

	err := tx.QueryRow(ctx, q).
		Scan(&s.TotalCategories, &s.TotalPrice)

	return s, err
}

type Stats struct {
	TotalItems      int     `json:"total_items"`
	TotalCategories int     `json:"total_categories"`
	TotalPrice      float64 `json:"total_price"`
}

func makePostgresURL(dbuser, dbpass, dbname, dbhost string) string {
	port := 5432
	sslmode := "disable"

	user := url.QueryEscape(dbuser)
	pass := url.QueryEscape(dbpass)
	host := url.QueryEscape(dbhost)
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, pass, host, port, dbname, sslmode,
	)
}
