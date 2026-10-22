package config

import (
	"errors"
	"io/fs"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config собирается из переменных окружения на старте приложения.
// У всего, кроме секретов, есть значение по умолчанию; секреты обязательны,
// поэтому приложение падает, если они не заданы.
type Config struct {
	// HTTPAddr это адрес, на котором слушает HTTP-сервер; catalog занимает :8081.
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8082"`
	// DBDSN это строка подключения к Postgres. Содержит пароль, значения по умолчанию нет.
	DBDSN string `env:"DB_DSN,required"`
	// LogLevel это уровень логирования: DEBUG, INFO, WARN, ERROR.
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
	// EnvMode различает окружения: development, staging, production.
	EnvMode string `env:"ENV_MODE" envDefault:"development"`

	// JWTSecret ключ подписи access-токенов (HS256), не короче 32 байт. Секрет, без значения по умолчанию.
	JWTSecret string `env:"JWT_SECRET,required"`
	// JWTIssuer пишется в iss и проверяется в Verify.
	JWTIssuer string `env:"JWT_ISSUER" envDefault:"rudstock-auth"`
	// AccessTTL срок жизни access-токена: после logout он проживёт ещё столько же.
	AccessTTL time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
	// RefreshTTL скользящий срок refresh: каждый refresh продлевает вход на столько.
	RefreshTTL time.Duration `env:"REFRESH_TTL" envDefault:"720h"`
	// RefreshMaxTTL абсолютный предел от ввода пароля, дальше только новый login.
	RefreshMaxTTL time.Duration `env:"REFRESH_MAX_TTL" envDefault:"2160h"`
}

// Load читает .env (если он есть), затем конфигурацию из окружения
// и проверяет обязательные поля.
func Load() (*Config, error) {
	// .env нужен только локально, в проде переменные задаются окружением.
	// Отсутствие файла не ошибка, а вот битый файл это ошибка.
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
