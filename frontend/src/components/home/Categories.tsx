import { Link } from 'react-router-dom'

const CATEGORIES = [
  { name: 'Conciertos', icon: '🎵', slug: 'conciertos' },
  { name: 'Teatro', icon: '🎭', slug: 'teatro' },
  { name: 'Deportes', icon: '⚽', slug: 'deportes' },
  { name: 'Conferencias', icon: '🎤', slug: 'conferencias' },
  { name: 'Festivales', icon: '🎪', slug: 'festivales' },
  { name: 'Talleres', icon: '🎨', slug: 'talleres' },
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
