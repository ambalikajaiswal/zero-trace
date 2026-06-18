import { useState, useEffect } from 'react'
import { useParams, useLocation } from 'react-router-dom'
import { getSecret } from '../api'

function ViewSecret() {
  const { id } = useParams<{ id: string }>()
  const location = useLocation()
  const [content, setContent] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [revealed, setRevealed] = useState(false)

  // Extract token from URL hash
  const token = location.hash.slice(1)

  useEffect(() => {
    if (!id || !token) {
      setError('Invalid secret link. The link may be incomplete or corrupted.')
    }
  }, [id, token])

  const handleReveal = async () => {
    if (!id || !token) return

    setLoading(true)
    setError('')

    try {
      const response = await getSecret(id, token)
      setContent(response.content)
      setRevealed(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to retrieve secret')
    } finally {
      setLoading(false)
    }
  }

  if (error && !revealed) {
    return (
      <div className="page">
        <div className="card error-card">
          <div className="error-icon">🚫</div>
          <h1>Unable to Access Secret</h1>
          <p>{error}</p>
          <p className="hint">
            The secret may have already been viewed, expired, or the link may be invalid.
          </p>
        </div>
      </div>
    )
  }

  if (revealed && content !== null) {
    return (
      <div className="page">
        <div className="card">
          <div className="revealed-header">
            <h1>🔓 Secret Revealed</h1>
            <p className="warning">This secret has been destroyed and cannot be viewed again.</p>
          </div>
          <div className="secret-content">
            <pre>{content}</pre>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="page">
      <div className="card">
        <div className="reveal-prompt">
          <span className="lock-icon">🔒</span>
          <h1>You've Received a Secret</h1>
          <p className="description">
            Someone shared a secret with you. Click below to reveal it.
          </p>
          <p className="warning">
            ⚠️ This secret may only be viewable once. Make sure you're ready.
          </p>
          <button
            onClick={handleReveal}
            className="btn btn-primary btn-large"
            disabled={loading}
          >
            {loading ? 'Decrypting...' : 'Reveal Secret'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default ViewSecret
