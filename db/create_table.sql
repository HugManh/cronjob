CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    execute VARCHAR(255) NOT NULL,
    message TEXT,
    hash VARCHAR(64) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    code VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_hash ON tasks (hash);
CREATE INDEX IF NOT EXISTS idx_tasks_name ON tasks (name);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks (deleted_at);

CREATE TABLE IF NOT EXISTS slacks (
    id SERIAL PRIMARY KEY,
    bot_token VARCHAR(255) NOT NULL,
    chat_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_slacks_bot_token ON slacks (bot_token);
CREATE INDEX IF NOT EXISTS idx_slacks_deleted_at ON slacks (deleted_at);
