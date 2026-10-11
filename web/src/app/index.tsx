import ProductCard from '@/features/catalog/components/product-card'
import type { CatalogCard } from '@/features/catalog/types/catalog-card'

// пример карточки в формате ответа GET /cards/{id}, пока API не подключён
const exampleCard: CatalogCard = {
  id: '3f2b8c1e-6a1d-4c2e-9b7a-2d5f0e8a1c44',
  category_id: 1,
  title: 'Michelin Primacy 4 205/55 R16',
  description: 'Шина летняя 205/55 R16',
  price: 14999,
  photo_url: '',
  created_at: '2026-10-05T12:00:00Z',
}

export default function App() {
  return (
    <main>
      <h1>RudStock</h1>
      <ProductCard card={exampleCard} />
    </main>
  )
}
