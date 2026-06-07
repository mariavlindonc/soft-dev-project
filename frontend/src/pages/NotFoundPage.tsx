import { Link } from 'react-router-dom'

export default function NotFoundPage() {
  return (
    <div className="page not-found">
      <h1>404</h1>
      <p>Página no encontrada</p>
      <Link to="/" className="btn btn-primary">Volver al inicio</Link>
    </div>
  )
}
