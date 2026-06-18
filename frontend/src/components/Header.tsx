import { Link } from 'react-router-dom'

function Header() {
  return (
    <header className="header">
      <div className="container header-content">
        <Link to="/" className="logo">
          <span className="logo-icon">🔐</span>
          <span className="logo-text">ZeroTrace</span>
        </Link>
        <p className="tagline">Share secrets securely. Zero trace left behind.</p>
      </div>
    </header>
  )
}

export default Header
