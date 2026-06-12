import { Link } from 'react-router-dom'

export default function HeroSection() {
  return (
    <section className="hero-section">
      <div className="hero-content">
        <h1>Descubre los mejores eventos</h1>
        <p>
          Encuentra y compra entradas para los eventos más destacados de tu ciudad.
          Conciertos, teatro, deportes y más.
        </p>
        <div className="hero-actions">
          <Link to="/events" className="btn btn-primary btn-lg">
            Explorar Eventos
          </Link>
          <Link to="/register" className="btn btn-outline btn-lg">
            Crear Cuenta
          </Link>
        </div>
      </div>
      <div className="hero-image" />
    </section>
  )
}
