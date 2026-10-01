package domain

// Category категория товаров (шины, диски, присадки).
type Category struct {
	ID   int64
	Name string
}

// CategoryPatch частичное обновление категории, nil означает "не менять поле".
type CategoryPatch struct {
	Name *string
}
