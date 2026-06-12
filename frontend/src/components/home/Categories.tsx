import { Link } from 'react-router-dom'

const CATEGORIES = [
  { name: 'Aire Libre', icon: '🌳', slug: 'aire libre' },
  { name: 'En Salón', icon: '🏛️', slug: 'en salon' },
  { name: 'Música', icon: '🎵', slug: 'musica' },
  { name: 'Teatro', icon: '🎭', slug: 'teatro' },
  { name: 'Gastronomía', icon: '🍽️', slug: 'gastronomia' },
]

export default function Categories() {
  return (
    <section className="categories-section">
      <div className="section-header">
        <h2>Categorías</h2>
      </div>
      <div className="categories-grid">
        {CATEGORIES.map((cat) => (
          <Link
            key={cat.slug}
            to={`/events?category=${cat.slug}`}
            className="category-card"
          >
            <span className="category-icon">{cat.icon}</span>
            <span className="category-name">{cat.name}</span>
          </Link>
        ))}
      </div>
    </section>
  )
}
