package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vald3mare/Rudstock/internal/platform/config"
	"github.com/vald3mare/Rudstock/internal/platform/httpx"
)

func Health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
}

func main() {
	// Конфиг читается первым: без него неизвестен даже уровень логирования,
	// поэтому ошибка выводится через fmt и процесс выходит с ненулевым кодом.
	conf, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	var handler slog.Handler

	var level slog.Level
	// Нераспознанный уровень не повод падать: продолжаем с INFO.
	if err := level.UnmarshalText([]byte(conf.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	// В продакшене JSON для машинного разбора, локально текст для чтения глазами.
	if conf.EnvMode == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	// Логгер ставится глобально: пакеты platform пишут через slog.Default()
	slog.SetDefault(slog.New(handler))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", Health)

	// Таймауты обязательны: без них зависшее соединение держит ресурсы бесконечно
	server := &http.Server{
		Addr:         conf.HTTPAddr, // Порт, который будет слушать сервер
		Handler:      mux,           // Маршрутизатор с нашими эндпоинтами
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}

	// Буфер на 1, чтобы горутина не зависла на отправке, если main уже вышел
	errChan := make(chan error, 1)

	go func() {
		// ErrServerClosed это штатный ответ на Shutdown, а не сбой
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// NotifyContext отменяет ctx по SIGINT или SIGTERM, закрывая ctx.Done().
	// stop снимает подписку, чтобы повторный сигнал убил процесс принудительно,
	// если плавная остановка подвиснет
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Ждём, что наступит раньше: падение сервера или сигнал на остановку
	select {
	case err := <-errChan:
		slog.Error("server failed", "error", err)
		return
	case <-ctx.Done():
		slog.Info("shutting down")
	}

	// Контекст новый, а не ctx: тот уже отменён, и Shutdown оборвал бы
	// соединения мгновенно. Пять секунд это потолок ожидания текущих запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
}
