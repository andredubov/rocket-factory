-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    login VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

INSERT INTO
    providers (name, description)
VALUES (
        'telegram',
        'Уведомления через Telegram бота'
    ),
    (
        'email',
        'Уведомления по электронной почте'
    ),
    (
        'push',
        'Push уведомления в браузере'
    ),
    ('sms', 'SMS уведомления');

CREATE TABLE notification_methods (
    user_uuid UUID NOT NULL REFERENCES users (uuid) ON DELETE CASCADE,
    provider_name VARCHAR(255) NOT NULL REFERENCES providers (name) ON DELETE CASCADE,
    target VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_uuid, provider_name)
);

-- Индексы для улучшения производительности
CREATE INDEX idx_users_login ON users (login);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_notification_methods_user ON notification_methods (user_uuid);

CREATE INDEX idx_providers_name ON providers (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_methods;

DROP TABLE IF EXISTS providers;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd