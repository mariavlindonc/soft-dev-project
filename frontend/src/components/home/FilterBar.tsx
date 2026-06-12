import { useNavigate } from 'react-router-dom'

const CATEGORIES = [
  { label: 'Aire Libre', slug: 'aire libre' },
  { label: 'En Salón', slug: 'en salon' },
  { label: 'Grupos Emergentes', slug: 'grupos emergentes' },
]

export default function FilterBar() {
  const navigate = useNavigate()

  return (
    <div className="filter-bar">
      <button type="button" className="filter-bar__location-btn" aria-label="Ubicación">
        &#9906;
      </button>
      <div className="filter-bar__dropdown">
        <span>Hoy</span>
        <span>&#9660;</span>
      </div>
      <div className="filter-bar__divider" />
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
