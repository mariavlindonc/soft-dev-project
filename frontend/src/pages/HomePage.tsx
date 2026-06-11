import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Navbar from '../components/layout/Navbar'
import HeroSection from '../components/home/HeroSection'
import FilterBar from '../components/home/FilterBar'
import FeaturedEvents from '../components/home/FeaturedEvents'
import Categories from '../components/home/Categories'
import HowItWorks from '../components/home/HowItWorks'

export default function HomePage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')

  function handleSearchKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter') {
      const q = search.trim()
      navigate(q ? `/events?q=${encodeURIComponent(q)}` : '/events')
    }
  }

  return (
    <div className="home">
      <Navbar />
      <main className="home__main">
        <HeroSection />
        <div className="search-bar-section">
          <div className="search-bar">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              placeholder="Buscar eventos, artistas, categorías..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              onKeyDown={handleSearchKeyDown}
            />
          </div>
        </div>
        <FilterBar />
        <FeaturedEvents />
        <Categories />
        <HowItWorks />
      </main>
    </div>
  )
}
