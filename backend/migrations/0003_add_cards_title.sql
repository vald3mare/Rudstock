-- +goose Up
-- DEFAULT '' нужен только чтобы заполнить уже существующие строки,
-- сразу после этого его снимаем: новые карточки обязаны приходить с title.
ALTER TABLE cards ADD COLUMN title TEXT NOT NULL DEFAULT '';
ALTER TABLE cards ALTER COLUMN title DROP DEFAULT;

-- +goose Down
ALTER TABLE cards DROP COLUMN title;
