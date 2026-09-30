package domain

import "errors"

var (
	ErrCardNotFound        = errors.New("card not found")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidCategoryID   = errors.New("invalid category id")
	ErrInvalidCategoryName = errors.New("invalid category name")
	ErrCategoryExists      = errors.New("category already exists")
	ErrInvalidPagination   = errors.New("invalid pagination")
)
