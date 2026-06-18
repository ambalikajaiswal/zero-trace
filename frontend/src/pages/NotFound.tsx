import { Link } from 'react-router-dom'

function NotFound() {
  return (
    <div className="page">
      <div className="card">
        <h1>404 — Not Found</h1>
        <p>The page you're looking for doesn't exist.</p>
        <Link to="/" className="btn btn-primary">
          Go Home
        </Link>
      </div>
    </div>
  )
}

export default NotFound
