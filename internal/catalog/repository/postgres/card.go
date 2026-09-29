package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/internal/catalog/domain"
)

// репа может быть открытой, тк сервис ничего не знает о репозитории, единственное где мы дергаем его, это в неэкспортируемом поле при вызове конструктора
type CardRepo struct {
	pool *pgxpool.Pool
}

func NewCardRepo(pool *pgxpool.Pool) *CardRepo {
	return &CardRepo{
		pool: pool,
	}
}

func (r *CardRepo) Create(ctx context.Context, card domain.Card) (int64, error) {
	return 0, nil
}
