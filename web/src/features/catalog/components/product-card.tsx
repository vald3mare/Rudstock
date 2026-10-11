import type { CatalogCard } from '../types/catalog-card'
import styles from './product-card.module.css'

// форма props: компонент ждёт объект с полем card типа CatalogCard
type ProductCardProps = {
  card: CatalogCard
}

// цена приходит в копейках; форматтер создаём один раз на модуль, а не на каждый рендер.
// Когда цена понадобится в корзине и заказах, переедет в src/utils
const priceFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
})

export default function ProductCard({ card }: ProductCardProps) {
  return (
    <article className={styles.card}>
      {card.photo_url ? (
        <img src={card.photo_url} alt={card.title} width={240} height={240} />
      ) : (
        <div role="img" aria-label="Фото отсутствует">
          Нет фото
        </div>
      )}

      <h2>{card.title}</h2>
      {card.description && <p>{card.description}</p>}

      {/* копейки -> рубли и формат "149,99 ₽" */}
      <p>{priceFormatter.format(card.price / 100)}</p>
    </article>
  )
}
