-- +goose Up
CREATE TABLE users (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    -- сервис приводит email к нижнему регистру; CHECK страхует UNIQUE от дублей вида A@x.ru / a@x.ru
    email          TEXT        NOT NULL UNIQUE CHECK (email = lower(email)),
    email_verified BOOLEAN     NOT NULL DEFAULT false,
    password_hash  TEXT        NOT NULL, -- argon2id
    role           TEXT        NOT NULL DEFAULT 'customer' CHECK (role IN ('customer', 'admin')),
    is_active      BOOLEAN     NOT NULL DEFAULT true,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;
