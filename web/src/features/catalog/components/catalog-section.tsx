import type { CatalogCard } from '@/features/catalog/types/catalog-card'

import styles from './catalog-section.module.css'
import ProductCard from './product-card'

type CatalogSectionProps = {
  cards: CatalogCard[]
}

// CatalogSection раздел каталога: заголовок и сетка карточек. Саму карточку рисует ProductCard
export default function CatalogSection({ cards }: CatalogSectionProps) {
  return (
    <section>
      <h2>Каталог</h2>
      {/* ul/li: диктор объявит "список, N элементов"; key нужен React, чтобы отличать элементы списка */}
      <ul className={styles.grid}>
        {cards.map((card) => (
          <li key={card.id}>
            <ProductCard card={card} />
          </li>
        ))}
      </ul>
    </section>
  )
}
