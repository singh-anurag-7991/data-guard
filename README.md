# DataGuard - B2B Data Quality & Validation SaaS

## Vision
DataGuard is a lightweight, strictly-typed data validation engine designed for B2B internal tools and data pipelines. It catches silent data failures in APIs and Databases before they break logic downstream. A "Firewall for Bad Data".

## Key Features (MVP)
- **Source Agnostic**: Works validation logic across API payloads and SQL Database rows.
- **Declarative Rule DSL**: JSON-based rules that are portable and readable.
- **Deterministic Engine**: Stateless Go-based engine ensures fast, predictable validation.
- **Deep Observability**: Structured logs and detailed failure reasons.

## Architecture
The system is composed of the following layers:
1. **Ingestion**: API webhooks or DB pollers normalize data into `domain.Record`.
2. **Validation Engine**: The core "Brain". Takes `Record` + `Rule` -> returns `ValidationResult`.
3. **Storage**: Stores Rules and historical Results (Postgres).
4. **Alerting**: State-machine based alerting to reduce noise (Slack/Email).

## Component Design
### Core Domain (`internal/domain`)
- **Rule**: Defines `Field`, `Check` (op + value), and `Condition` (when to run).
- **ValidationResult**: The output containing status and error details.

### Rule DSL
Example Rule:
```json
{
  "id": "amount_positive",
  "field": "amount",
  "when": {
    "field": "status",
    "op": "eq",
    "value": "PAID"
  },
  "checks": [
    { "op": "gt", "value": 0 }
  ]
}
```

## Getting Started
### Prerequisites
- Go 1.22+
- Docker (optional for db)

### Running Tests
```bash
go test ./... -cover
```

### Running the Server
```bash
go run cmd/server/main.go
# Server listens on :8080
```

### Ingesting Data (Example)
```bash
curl -X POST localhost:8080/ingest/api -d '{
  "source_id": "test",
  "schema": {"amount": "number"},
  "rules": [{"id": "r1", "field": "amount", "checks": [{"op": "gt", "value": 0}]}],
  "data": [{"amount": 10}, {"amount": -5}]
}'
```

## Roadmap
- [ ] Phase 1: Core Engine (Done)
- [ ] Phase 2: Ingestion Layers (Next)
- [ ] Phase 3: DB Optimization (Query Builder)
- [ ] Phase 4: Alerting & Storage
