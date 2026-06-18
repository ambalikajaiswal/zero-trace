import { useLocation, Navigate } from 'react-router-dom'
import { useState } from 'react'

interface ShareState {
  id: string
  token: string
  expiresAt: string
}

function ShareLink() {
  const location = useLocation()
  const state = location.state as ShareState | null
  const [copied, setCopied] = useState(false)

  if (!state) {
    return <Navigate to="/" replace />
  }

  const shareUrl = `${window.location.origin}/secret/${state.id}#${state.token}`
  const expiresAt = new Date(state.expiresAt)

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(shareUrl)
      setCopied(true)
      setTimeout(() => setCopied(false), 3000)
    } catch {
      // Fallback for older browsers
      const textArea = document.createElement('textarea')
      textArea.value = shareUrl
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      setCopied(true)
      setTimeout(() => setCopied(false), 3000)
    }
  }

  return (
    <div className="page">
      <div className="card success-card">
        <div className="success-icon">✅</div>
        <h1>Secret Created!</h1>
        <p className="description">
          Share this link with the intended recipient. The secret will be destroyed after viewing or expiration.
        </p>

        <div className="share-link-container">
          <input
            type="text"
            value={shareUrl}
            readOnly
            className="share-link-input"
            aria-label="Secret share link"
          />
          <button onClick={handleCopy} className="btn btn-copy">
            {copied ? '✓ Copied!' : 'Copy Link'}
          </button>
        </div>

        <div className="meta-info">
          <p>
            <strong>Expires:</strong> {expiresAt.toLocaleString()}
          </p>
          <p className="warning">
            ⚠️ This link will only work once. Save it now — you won't see it again.
          </p>
        </div>
      </div>
    </div>
  )
}

export default ShareLink
