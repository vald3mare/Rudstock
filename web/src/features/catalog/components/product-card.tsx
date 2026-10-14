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
        // классы из модуля только через styles.*: строка "photo" не совпадёт с уникальным именем после сборки
        <img className={styles.photo} src={card.photo_url} alt={card.title} />
      ) : (
        <div className={styles.photoPlaceholder} role="img" aria-label="Фото отсутствует">
          Нет фото
        </div>
      )}

      {/*Контент карточки*/}
      {/* h3: страница h1 -> секция каталога h2 -> карточка h3, уровни не перескакиваем */}
      <h3>{card.title}</h3>
      {card.description && <p>{card.description}</p>}

      {/* копейки -> рубли и формат "149,99 ₽" */}
      <p className={styles.price}>{priceFormatter.format(card.price / 100)}</p>
    </article>
  )
}
