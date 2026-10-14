package auth

import "context"

//   "auth": {
//     "POST /auth/register": { "status": "todo", "request": ["email", "password"], "response": ["access", "refresh", "user: {id, email, role}"], "errors": [400, 409, 429], "note": "роль всегда customer; name и phone не спрашиваем, заполняются в PATCH /me или при первом заказе" },
//     "POST /auth/login": { "status": "todo", "request": ["email", "password"], "response": ["access", "refresh", "user: {id, email, role}"], "errors": [401, 403, 429] },
//     "POST /auth/refresh": { "status": "todo", "request": ["refresh"], "response": ["access", "refresh"], "errors": [401] },
//     "POST /auth/logout": { "status": "todo", "request": ["refresh"] },
//     "POST /auth/email/verify": { "status": "todo", "request": ["token"], "errors": [400] },
//     "POST /auth/password/forgot": { "status": "todo", "request": ["email"], "errors": [429], "note": "204 всегда, даже для несуществующего email" },
//     "POST /auth/password/reset": { "status": "todo", "request": ["token", "password"], "errors": [400] }
//   },

func (s *authService) Login(ctx context.Context, in LoginUserInput) (AuthResult, error) {

	return AuthResult{}, nil
}
