# DataGuard - B2B Data Quality & Validation SaaS

[![Go CI](https://github.com/singh-anurag-7991/data-guard/actions/workflows/go.yml/badge.svg)](https://github.com/singh-anurag-7991/data-guard/actions/workflows/go.yml)

## Vision
DataGuard is a lightweight, strictly-typed data validation engine designed for B2B internal tools and data pipelines. It catches silent data failures in APIs and Databases before they break logic downstream. A "Firewall for Bad Data".

## Key Features
- **Source Agnostic**: Works validation logic across API payloads and SQL Database rows.
- **SQL Pushdown Optimization**: Converts validation rules into optimized SQL queries to find failures without fetching all data.
- **Declarative Rule DSL**: JSON-based rules that are portable and readable.
- **Real-time Dashboard**: Next.js UI to visualize validation history and health status.
- **Smart Alerting**: State-machine based alerting (Slack) that only notifies on status changes (PASS <-> FAIL) to look reduce noise.

## Architecture
The system is composed of the following layers:
1.  **Ingestion**: 
    -   API Webhook (`POST /ingest/api`)
    -   Postgres Connector (Pull-based)
2.  **Validation Engine**: 
    -   **Standard**: Stateless Go-based memory execution.
    -   **Optimizer**: Translates rules to SQL WHERE clauses for failure detection.
3.  **Storage**: Postgres persistence for Validation Runs, Errors, and Alert States.
4.  **Dashboard**: Next.js (React) Frontend for monitoring.

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+ (for Dashboard)
- PostgreSQL 14+ (optional, required history/alerting)

### 1. Running Backend (Go)
```bash
# Optional: Set DB URL for persistence
export DATABASE_URL="postgres://user:pass@localhost:5432/dataguard"

go run cmd/server/main.go
# Server listens on :8080
```

### 2. Running Frontend (Next.js)
```bash
cd web
npm install
npm run dev
# Dashboard available at http://localhost:3000
```

### 3. Running with Docker (Recommended)
Build and run the entire backend as a single container:
```bash
docker build -t dataguard .
docker run -p 8080:8080 -e DATABASE_URL=... dataguard
```

## API Usage
**Endpoint**: `POST /ingest/api`

**Payload**:
```json
{
  "source_id": "orders",
  "schema": { "amount": "number" },
  "rules": [
    {
      "id": "positive_amount",
      "field": "amount",
      "checks": [{ "op": "gt", "value": 0 }]
    }
  ],
  "data": [
    { "amount": 100 },
    { "amount": -10 }
  ]
}
```

## Roadmap
- [x] **Phase 1**: Core Engine (Memory)
- [x] **Phase 2**: Ingestion Layers (API & Postgres)
- [x] **Phase 3**: DB Optimization (SQL Query Builder)
- [x] **Phase 4**: Storage & Alerting (State Machine)
- [x] **Phase 5**: Dashboard (Next.js) & API
- [x] **Phase 6**: Polish & Deploy (Docker + CI)
