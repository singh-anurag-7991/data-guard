#!/bin/bash

BASE_URL="http://localhost:8080/ingest/api"

# 1. Source: ecommerce_orders (Mixed)
echo "Seeding ecommerce_orders (PASS)..."
curl -s -X POST $BASE_URL -d '{
  "source_id": "ecommerce_orders",
  "schema": { "total": "number" },
  "rules": [{ "id": "pos", "field": "total", "checks": [{ "op": "gt", "value": 0 }] }],
  "data": [{ "total": 100 }, { "total": 50 }]
}' > /dev/null

echo "Seeding ecommerce_orders (FAIL)..."
curl -s -X POST $BASE_URL -d '{
  "source_id": "ecommerce_orders",
  "schema": { "total": "number" },
  "rules": [{ "id": "pos", "field": "total", "checks": [{ "op": "gt", "value": 0 }] }],
  "data": [{ "total": -10 }]
}' > /dev/null

# 2. Source: user_signups (PASS)
echo "Seeding user_signups (PASS)..."
curl -s -X POST $BASE_URL -d '{
  "source_id": "user_signups",
  "schema": { "email": "string" },
  "rules": [{ "id": "email_req", "field": "email", "checks": [{ "op": "not_null" }] }],
  "data": [{ "email": "alice@example.com" }, { "email": "bob@example.com" }]
}' > /dev/null

# 3. Source: iot_sensors (FAIL)
echo "Seeding iot_sensors (FAIL)..."
curl -s -X POST $BASE_URL -d '{
  "source_id": "iot_sensors",
  "schema": { "temp": "number" },
  "rules": [{ "id": "temp_safe", "field": "temp", "checks": [{ "op": "lt", "value": 100 }] }],
  "data": [{ "temp": 150 }, { "temp": 200 }]
}' > /dev/null

echo "Seeding Complete!"
