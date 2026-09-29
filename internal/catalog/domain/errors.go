package domain

import "errors"

var (
	ErrCardNotFound     = errors.New("card not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidPrice     = errors.New("invalid price")
)
