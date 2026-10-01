-- +goose Up
CREATE TABLE cards (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id BIGINT      NOT NULL REFERENCES categories (id),
    description TEXT        NOT NULL DEFAULT '',
    price       BIGINT      NOT NULL CHECK (price > 0), -- копейки
    photo_url   TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX cards_category_id_idx ON cards (category_id);

-- +goose Down
DROP TABLE cards;