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

// User пользователь: покупатели и сотрудники в одной таблице, различаются Role.
type User struct {
	ID            uuid.UUID
	Email         string // в нижнем регистре, без пробелов по краям
	EmailVerified bool
	PasswordHash  string // argon2id; наружу не отдаётся никогда
	Name          string
	Phone         string // необязательный, пустая строка если не указан
	Role          Role
	IsActive      bool // false = заблокирован: login отвечает 403
	CreatedAt     time.Time
}
