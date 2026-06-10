const CATEGORIES = ['Conciertos', 'Teatro', 'Deportes', 'Conferencias', 'Festivales', 'Talleres']

export default function FilterBar() {
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
        <button key={cat} type="button" className="filter-bar__pill">
          {cat}
        </button>
      ))}
    </div>
  )
}
