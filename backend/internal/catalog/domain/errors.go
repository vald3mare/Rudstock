package domain

import "errors"

var (
	ErrCardNotFound        = errors.New("card not found")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidCardTitle    = errors.New("invalid card title")
	ErrInvalidCategoryID   = errors.New("invalid category id")
	ErrInvalidCategoryName = errors.New("invalid category name")
	ErrCategoryExists      = errors.New("category already exists")
	ErrCategoryHasCards    = errors.New("category has cards")
	ErrInvalidPagination   = errors.New("invalid pagination")
	ErrEmptyPatch          = errors.New("nothing to update")
)
