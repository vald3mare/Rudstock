import type { CatalogCard } from '@/features/catalog/types/catalog-card'
import { api } from '@/lib/api-client'

// ответ GET /cards
export type CatalogCardList = {
  items: CatalogCard[]
  total: number // сколько карточек под фильтром всего, для пагинации
}

// query-параметры GET /cards; все необязательные, дефолты ставит бэкенд
export type GetCardsParams = {
  categoryId?: number
  page?: number
  limit?: number
}

// путь ровно как в контракте; адрес сервиса, статус и JSON берёт на себя api-client
export function getCards(params: GetCardsParams = {}): Promise<CatalogCardList> {
  return api.get<CatalogCardList>('/cards', {
    category_id: params.categoryId,
    page: params.page,
    limit: params.limit,
  })
}
