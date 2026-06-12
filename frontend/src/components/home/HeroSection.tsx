import { useState, useEffect, useCallback } from 'react'
import { Link } from 'react-router-dom'
import banner1 from '../../assets/banner1.jpg'
import banner2 from '../../assets/banner2.jpg'
import banner3 from '../../assets/banner3.jpg'

const BANNERS = [banner1, banner2, banner3]

export default function HeroSection() {
  const [activeIndex, setActiveIndex] = useState(0)

  const next = useCallback(() => {
    setActiveIndex((prev) => (prev + 1) % BANNERS.length)
  }, [])

  useEffect(() => {
    const timer = setInterval(next, 5000)
    return () => clearInterval(timer)
  }, [next])

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
            {BANNERS.map((_, i) => (
              <button
                key={i}
                type="button"
                className={`hero__dot${i === activeIndex ? ' hero__dot--active' : ''}`}
                aria-label={`Slide ${i + 1}`}
                onClick={() => setActiveIndex(i)}
              />
            ))}
          </div>
        </div>
        <div className="hero__right">
          {BANNERS.map((src, i) => (
            <div
              key={i}
              className="hero__carousel-image"
              style={{
                backgroundImage: `url(${src})`,
                opacity: i === activeIndex ? 1 : 0,
              }}
            />
          ))}
        </div>
      </div>
    </section>
  )
}
