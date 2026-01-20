package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/singh-anurag-7991/data-guard/internal/domain"
)

type Client struct {
	pool *pgxpool.Pool
}

// NewClient creates a new Postgres client
func NewClient(ctx context.Context, connString string) (*Client, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Safety settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &Client{pool: pool}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

// FetchRows executes a query and returns normalized records
func (c *Client) FetchRows(ctx context.Context, query string, args ...interface{}) ([]domain.Record, error) {
	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	columnNames := make([]string, len(fields))
	for i, fd := range fields {
		columnNames[i] = string(fd.Name)
	}

	var records []domain.Record
	for rows.Next() {
		// Create a slice of interface{} to hold values
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		// Map to domain.Record
		record := make(domain.Record)
		for i, col := range columnNames {
			val := values[i]
			// Handle []byte for strings/text
			if b, ok := val.([]byte); ok {
				record[col] = string(b)
			} else {
				record[col] = val
			}
		}
		records = append(records, record)
	}

	return records, nil
}

// ValidateViaSQL executes a generated failure query and returns the FAILING records
func (c *Client) ValidateViaSQL(ctx context.Context, query string, args []interface{}) ([]domain.Record, error) {
	return c.FetchRows(ctx, query, args...)
}
