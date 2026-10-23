package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/vald3mare/Rudstock/auth/internal/platform/errs"
	"github.com/vald3mare/Rudstock/auth/internal/platform/httpx"
	"github.com/vald3mare/Rudstock/auth/internal/tokens"
)

// verifier что middleware нужно от пакета tokens
type verifier interface {
	Verify(access string) (tokens.Identity, error)
}

// identityKey закрытый тип ключа: чужой пакет не сможет ни прочитать, ни подменить значение в контексте
type identityKey struct{}

// Middleware пропускает запрос только с валидным Authorization: Bearer <access>
// и кладёт в контекст, кто пришёл; хендлеры достают это через IdentityFrom.
func Middleware(v verifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access, ok := bearer(r)
		if !ok {
			httpx.WriteError(w, r, errs.Unauthorized("unauthorized", "Authorization required", nil))
			return
		}

		identity, err := v.Verify(access)
		if err != nil {
			// отдельный слаг для истёкшего: фронт по нему идёт в /auth/refresh, а не на страницу входа
			if errors.Is(err, tokens.ErrTokenExpired) {
				httpx.WriteError(w, r, errs.Unauthorized("token-expired", "Access token expired", err))
				return
			}
			httpx.WriteError(w, r, errs.Unauthorized("invalid-token", "Invalid access token", err))
			return
		}

		ctx := context.WithValue(r.Context(), identityKey{}, identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IdentityFrom возвращает того, кого пропустил Middleware. false значит ручка не под Middleware.
func IdentityFrom(ctx context.Context) (tokens.Identity, bool) {
	identity, ok := ctx.Value(identityKey{}).(tokens.Identity)
	return identity, ok
}

// bearer достаёт токен из заголовка; схема по RFC 6750 регистронезависимая
func bearer(r *http.Request) (string, bool) {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}

	token = strings.TrimSpace(token)
	return token, token != ""
}
