import { Link } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

export default function Navbar() {
  const { isAuthenticated, isAdmin, user, logout } = useAuth()

  return (
    <nav className="navbar">
      <Link to="/" className="navbar__logo">Ceibo</Link>

      <div className="navbar__primary-links">
        <Link to="/events">Eventos</Link>
        {isAdmin && <Link to="/admin">Admin</Link>}
      </div>

      <div className="navbar__secondary-links">
        {isAuthenticated ? (
          <>
            <Link to="/tickets">Mis Entradas</Link>
            <span className="navbar__user-name">{user?.name}</span>
            <button type="button" className="navbar__logout" onClick={logout}>
              Cerrar Sesión
            </button>
            <div className="navbar__avatar">
              {user?.name?.charAt(0).toUpperCase() ?? 'U'}
            </div>
          </>
        ) : (
          <>
            <Link to="/login">Iniciar Sesión</Link>
            <Link to="/register">Registrarse</Link>
          </>
        )}
      </div>
    </nav>
  )
}
