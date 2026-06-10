import { Link } from 'react-router-dom'

export default function HeroSection() {
  return (
    <section className="hero">
      <div className="hero__card">
        <div className="hero__left">
          <h1 className="hero__title">
            Descubre los<br />mejores eventos
          </h1>
          <Link to="/events" className="hero__cta">
            Explorar Eventos
          </Link>
          <div className="hero__dots">
            <button type="button" className="hero__dot hero__dot--active" aria-label="Slide 1" />
            <button type="button" className="hero__dot" aria-label="Slide 2" />
            <button type="button" className="hero__dot" aria-label="Slide 3" />
          </div>
        </div>
        <div className="hero__right">
          <div className="hero__event-label">Events</div>
          <div className="hero__favorites">
            &#9829; 12.8k
          </div>
        </div>
      </div>
    </section>
  )
}
