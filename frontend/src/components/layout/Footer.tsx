import { Link } from 'react-router-dom'

export default function Footer() {
  return (
    <footer className="footer">
      <div className="footer-inner">
        <div className="footer-section">
          <h3>Ceibo</h3>
          <p>Tu plataforma de eventos y entradas</p>
        </div>
        <div className="footer-section">
          <h4>Enlaces</h4>
          <Link to="/events">Eventos</Link>
        </div>
        <div className="footer-section">
          <h4>Soporte</h4>
          <Link to="/faq">FAQ</Link>
          <Link to="/terms">Términos y Condiciones</Link>
          <Link to="/privacy">Privacidad</Link>
        </div>
      </div>
      <div className="footer-bottom">
        <p>&copy; {new Date().getFullYear()} Ceibo. Todos los derechos reservados.</p>
      </div>
    </footer>
  )
}
