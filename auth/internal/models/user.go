package models

import (
	"time"

	"github.com/google/uuid"
)

// Role роль пользователя. При регистрации всегда RoleCustomer, из запроса не берётся.
type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

// User учётные данные пользователя: покупатели и сотрудники в одной таблице, различаются Role.
// Профиль (имя, телефон) сюда не входит: это данные будущего сервиса пользователей, связь по ID.
type User struct {
	ID            uuid.UUID
	Email         string // в нижнем регистре, без пробелов по краям
	EmailVerified bool   // ставит POST /auth/email/verify
	PasswordHash  string // argon2id; наружу не отдаётся никогда
	Role          Role
	IsActive      bool // false = заблокирован: login отвечает 403
	CreatedAt     time.Time
}
