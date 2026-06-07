import { Link } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

export default function Navbar() {
  const { isAuthenticated, isAdmin, user, logout } = useAuth()

  return (
    <nav className="navbar">
      <div className="navbar-inner">
        <Link to="/" className="navbar-brand">
          Ceibo
        </Link>

        <div className="navbar-links">
          <Link to="/events" className="nav-link">Eventos</Link>
          {isAdmin && <Link to="/admin" className="nav-link">Admin</Link>}
        </div>

        <div className="navbar-actions">
          {isAuthenticated ? (
            <>
              <span className="navbar-user">{user?.name}</span>
              <Link to="/tickets" className="nav-link">Mis Entradas</Link>
              <button type="button" className="btn btn-outline" onClick={logout}>
                Cerrar Sesión
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="btn btn-outline">Iniciar Sesión</Link>
              <Link to="/register" className="btn btn-primary">Registrarse</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  )
}
