// Category категория товаров в том виде, в каком её отдаёт сервис каталога
// (GET /categories, ответ PATCH /admin/categories/{id}; CategoryResponse в catalog/internal/controllers/category/category_dto.go).
// Имена полей совпадают с JSON один в один, поэтому ответ можно использовать без преобразований.
export type Category = {
  id: number // int64 на бэке, не UUID, в отличие от карточки
  name: string
}
