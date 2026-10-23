import CatalogSection from '@/features/catalog/components/catalog-section'
import type { CatalogCard } from '@/features/catalog/types/catalog-card'

// примеры карточек в формате ответа GET /cards, пока API не подключён
const exampleCards: CatalogCard[] = [
  {
    id: '3f2b8c1e-6a1d-4c2e-9b7a-2d5f0e8a1c44',
    category_id: 1,
    title: 'Michelin Primacy 4 205/55 R16',
    description: 'Шина летняя 205/55 R16',
    price: 14999,
    photo_url: '',
    created_at: '2026-10-05T12:00:00Z',
  },
  {
    id: '7a1c0e52-3b9d-4f6e-8c21-5d4e3f2a1b09',
    category_id: 1,
    title: 'Nokian Hakkapeliitta R5 215/60 R17',
    description: 'Шина зимняя нешипованная',
    price: 1289000,
    photo_url: '',
    created_at: '2026-10-05T12:00:00Z',
  },
  {
    id: 'c4d2b6a8-9e1f-4a3b-b7c5-0f8e6d4c2a10',
    category_id: 2,
    title: 'Диск литой K&K Сочи R16',
    description: '',
    price: 640000,
    photo_url: '',
    created_at: '2026-10-05T12:00:00Z',
  },
  {
    id: 'e9f8a7b6-c5d4-4e3f-a2b1-c0d9e8f7a6b5',
    category_id: 3,
    title: 'Присадка в топливо Liqui Moly',
    description: 'Очиститель инжектора',
    price: 89000,
    photo_url: '',
    created_at: '2026-10-05T12:00:00Z',
  },
]

export default function App() {
  return (
    <main>
      <h1>RudStock</h1>
      <CatalogSection cards={exampleCards} />
    </main>
  )
}
