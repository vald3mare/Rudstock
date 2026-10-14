package category

// CreateCategoryInput входные данные для создания категории.
type CreateCategoryInput struct {
	Name string // обязательно, не пустое, уникальное
}

// UpdateCategoryInput входные данные для частичного обновления категории.
// nil означает "не менять поле"; хотя бы одно поле должно быть задано.
// Поле пока одно, но указатель оставляем как у карточки: новые поля (sort и т.п.) добавятся без смены подхода.
type UpdateCategoryInput struct {
	Name *string // если задано: не пустое, уникальное
}
