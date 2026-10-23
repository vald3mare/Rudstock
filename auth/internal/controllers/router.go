package controllers

import (
	"net/http"

	"github.com/vald3mare/Rudstock/auth/internal/controllers/auth"
)

func NewRouter(
	credentials *auth.CredentialsHandler,
	sessions *auth.SessionHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", credentials.Register)
	mux.HandleFunc("POST /auth/login", credentials.Login)
	mux.HandleFunc("POST /auth/refresh", sessions.Refresh)
	mux.HandleFunc("POST /auth/logout", sessions.Logout)

	return mux
}
