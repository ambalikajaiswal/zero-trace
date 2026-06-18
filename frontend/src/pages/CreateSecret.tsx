import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createSecret } from '../api'

function CreateSecret() {
  const navigate = useNavigate()
  const [content, setContent] = useState('')
  const [ttl, setTtl] = useState(3600)
  const [viewOnce, setViewOnce] = useState(true)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const ttlOptions = [
    { label: '5 minutes', value: 300 },
    { label: '30 minutes', value: 1800 },
    { label: '1 hour', value: 3600 },
    { label: '24 hours', value: 86400 },
    { label: '7 days', value: 604800 },
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await createSecret({
        content,
        ttl_seconds: ttl,
        view_once: viewOnce,
      })

      // Navigate to share page with the link info
      navigate('/shared', {
        state: {
          id: response.id,
          token: response.token,
          expiresAt: response.expires_at,
        },
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create secret')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page">
      <div className="card">
        <h1>Create a Secret</h1>
        <p className="description">
          Your secret will be encrypted and destroyed after viewing or expiration.
        </p>

        <form onSubmit={handleSubmit} className="form">
          <div className="form-group">
            <label htmlFor="secret-content">Secret Content</label>
            <textarea
              id="secret-content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Enter your secret (password, API key, private note...)"
              rows={5}
              maxLength={10000}
              required
              aria-describedby="content-hint"
            />
            <small id="content-hint">{content.length}/10,000 characters</small>
          </div>

          <div className="form-group">
            <label htmlFor="ttl-select">Expires After</label>
            <select
              id="ttl-select"
              value={ttl}
              onChange={(e) => setTtl(Number(e.target.value))}
            >
              {ttlOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group checkbox-group">
            <label htmlFor="view-once">
              <input
                id="view-once"
                type="checkbox"
                checked={viewOnce}
                onChange={(e) => setViewOnce(e.target.checked)}
              />
              Destroy after viewing (view once)
            </label>
          </div>

          {error && (
            <div className="error" role="alert">
              {error}
            </div>
          )}

          <button type="submit" className="btn btn-primary" disabled={loading || !content}>
            {loading ? 'Creating...' : 'Create Secret Link'}
          </button>
        </form>
      </div>

      <div className="features">
        <div className="feature">
          <span className="feature-icon">⏱️</span>
          <h3>TTL Expiration</h3>
          <p>Secrets auto-destruct after your chosen time limit.</p>
        </div>
        <div className="feature">
          <span className="feature-icon">👁️</span>
          <h3>View Once</h3>
          <p>Secrets are destroyed immediately after being viewed.</p>
        </div>
        <div className="feature">
          <span className="feature-icon">🔑</span>
          <h3>Zero Trust</h3>
          <p>Access requires a unique token. No stored logs or traces.</p>
        </div>
      </div>
    </div>
  )
}

export default CreateSecret
