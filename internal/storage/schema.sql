-- Schema for DataGuard

CREATE TABLE IF NOT EXISTS validation_runs (
    id SERIAL PRIMARY KEY,
    source_id TEXT NOT NULL,
    status TEXT NOT NULL, -- "PASS", "FAIL"
    records_checked INT NOT NULL,
    rules_failed INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS validation_errors (
    id SERIAL PRIMARY KEY,
    run_id INT REFERENCES validation_runs(id) ON DELETE CASCADE,
    rule_id TEXT NOT NULL,
    field TEXT NOT NULL,
    fail_value TEXT, -- Stores value as string
    reason TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS alert_states (
    source_id TEXT PRIMARY KEY,
    last_status TEXT NOT NULL, -- "PASS", "FAIL"
    last_alerted_at TIMESTAMP WITH TIME ZONE
);
