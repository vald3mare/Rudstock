import type { CatalogCard } from '../types/catalog-card'

export default function CatalogSection({ cards }: { cards: CatalogCard[] }) {
  return (
    <section>
      <h1>Каталог</h1>
      <div className="catalog-grid">
        {cards.map((card) => (
          <div key={card.id} className="catalog-card">
            <img src={card.photo_url} alt={card.title} />
            <h2>{card.title}</h2>
            <p>{card.description}</p>
            <p>{(card.price / 100).toFixed(2)} ₽</p>
          </div>
        ))}
      </div>
    </section>
  )
}
