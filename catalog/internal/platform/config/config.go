package config

import (
	"errors"
	"io/fs"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config собирается из переменных окружения на старте приложения.
// У всего, кроме секретов, есть значение по умолчанию; секреты обязательны,
// поэтому приложение падает, если они не заданы.
type Config struct {
	// HTTPAddr это адрес, на котором слушает HTTP-сервер, например ":8080".
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8081"`
	// DBDSN это строка подключения к Postgres. Содержит пароль,
	// поэтому значения по умолчанию у неё нет: локально берётся из .env.
	DBDSN string `env:"DB_DSN,required"`
	// LogLevel это уровень логирования: DEBUG, INFO, WARN, ERROR.
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
	// EnvMode различает окружения: development, staging, production.
	EnvMode string `env:"ENV_MODE" envDefault:"development"`
	// CORSOrigins это адреса фронта, которым браузер разрешит читать ответы, через запятую.
	CORSOrigins []string `env:"CORS_ORIGINS" envSeparator:"," envDefault:"http://localhost:5173"`
	// DebugToken это значение X-Debug-Token для /admin/* до появления сервиса auth.
	// Секрет, поэтому без значения по умолчанию.
	DebugToken string `env:"DEBUG_TOKEN,required"`
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
