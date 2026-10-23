-- +goose Up
-- Одна строка на каждую выдачу refresh-токена. Использованные и отозванные строки не удаляются сразу:
-- по ним ловим повторное использование токена. Чистить можно строки с истёкшим expires_at.
CREATE TABLE sessions (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    family_id        UUID        NOT NULL,
    token_hash       BYTEA       NOT NULL UNIQUE CHECK (octet_length(token_hash) = 32), -- sha256
    authenticated_at TIMESTAMPTZ NOT NULL, -- переносится при ротации
    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at       TIMESTAMPTZ,          -- NULL = живая
    CHECK (expires_at > created_at)
);

-- поиск по token_hash покрывает UNIQUE; эти два для массового гашения
CREATE INDEX sessions_family_id_idx ON sessions (family_id); -- повторное использование: гасим семью
CREATE INDEX sessions_user_id_idx ON sessions (user_id);     -- сброс пароля: гасим все сессии пользователя

-- +goose Down
DROP TABLE sessions;
