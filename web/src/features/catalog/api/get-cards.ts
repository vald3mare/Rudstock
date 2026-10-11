import type { CatalogCard } from '@/features/catalog/types/catalog-card'

// возвращаемый тип для GET /cards и GET /cards/{id}
export type CatalogCardList = {
  items: CatalogCard[]
  total: number // сколько карточек под фильтром всего, для пагинации
}
