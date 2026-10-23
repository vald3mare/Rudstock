// CatalogCard карточка товара в том виде, в каком её отдаёт сервис каталога
// (GET /cards, GET /cards/{id}; CardResponse в catalog/internal/controllers/card/card_dto.go).
// Имена полей совпадают с JSON один в один, поэтому ответ можно использовать без преобразований.
export type CatalogCard = {
  id: string // UUID
  category_id: number
  title: string
  description: string // может быть пустой строкой
  price: number // в копейках: 14999 = 149,99 ₽, в рубли переводить только при выводе
  photo_url: string // ссылка на объект в minio, может быть пустой строкой
  created_at: string // дата в ISO 8601, JSON не передаёт Date, только строку
}
