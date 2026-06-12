import { useNavigate } from 'react-router-dom'

const CATEGORIES = [
  { label: 'Aire Libre', slug: 'aire libre' },
  { label: 'En Salón', slug: 'en salon' },
  { label: 'Música', slug: 'musica' },
  { label: 'Teatro', slug: 'teatro' },
  { label: 'Gastronomía', slug: 'gastronomia' },
]

export default function FilterBar() {
  const navigate = useNavigate()

  return (
    <div className="filter-bar">
      <button type="button" className="filter-bar__location-btn" aria-label="Ubicación">
        &#9906;
      </button>

      {CATEGORIES.map((cat) => (
        <button
          key={cat.slug}
          type="button"
          className="filter-bar__pill"
          onClick={() => navigate(`/events?category=${cat.slug}`)}
        >
          {cat.label}
        </button>
      ))}
    </div>
  )
}
