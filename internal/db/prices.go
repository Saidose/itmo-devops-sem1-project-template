package db

import (
	"context"
	"fmt"
	"net/url"
	"project_sem/internal/domain"

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
	return &PricesDB{pool: pool}
}

func (r *PricesDB) Migrate(ctx context.Context) error {
	const query = `CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    create_date DATE NOT NULL,
)`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *PricesDB) Create(ctx context.Context, price domain.Price) error {
	const query = `
	INSERT INTO prices (
						id,
						name,
						category,
						price,
						create_date
						)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		price.ID,
		price.Name,
		price.Category,
		price.Price,
		price.CreateDate)
	if err != nil {
		return err
	}
	return nil
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
