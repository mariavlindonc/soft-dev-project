import { Outlet, useLocation } from 'react-router-dom'
import Navbar from './Navbar'
import Footer from './Footer'

export default function Layout() {
  const location = useLocation()
  const isHome = location.pathname === '/'

  return (
    <div className="app-layout">
      {!isHome && <Navbar />}
      <main className="main-content">
        <Outlet />
      </main>
      <Footer />
    </div>
  )
}
