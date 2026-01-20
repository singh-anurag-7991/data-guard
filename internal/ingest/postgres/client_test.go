package postgres

import (
	"context"
	"os"
	"testing"
)

func TestClient_FetchRows(t *testing.T) {
	connStr := os.Getenv("TEST_DB_URL")
	if connStr == "" {
		t.Skip("Skipping postgres integration test: TEST_DB_URL not set")
	}

	ctx := context.Background()
	client, err := NewClient(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	// Assumption: Setup data manually or assume 'information_schema' exists
	rows, err := client.FetchRows(ctx, "SELECT table_name FROM information_schema.tables LIMIT 1")
	if err != nil {
		t.Fatalf("failed to fetch rows: %v", err)
	}

	if len(rows) == 0 {
		t.Log("Query ran but returned no rows (expected if DB empty)")
	} else {
		t.Logf("Fetched row: %v", rows[0])
	}
}
