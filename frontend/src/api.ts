const API_BASE = '/api'

export interface CreateSecretRequest {
  content: string
  ttl_seconds: number
  view_once: boolean
}

export interface CreateSecretResponse {
  id: string
  token: string
  expires_at: string
  share_url: string
}

export interface GetSecretResponse {
  content: string
  created_at: string
  expires_at: string
}

export interface ApiError {
  error: string
  message: string
}

export async function createSecret(data: CreateSecretRequest): Promise<CreateSecretResponse> {
  const response = await fetch(`${API_BASE}/secrets`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    const err: ApiError = await response.json()
    throw new Error(err.message)
  }

  return response.json()
}

export async function getSecret(id: string, token: string): Promise<GetSecretResponse> {
  const response = await fetch(`${API_BASE}/secrets/${id}`, {
    headers: { 'X-Access-Token': token },
  })

  if (!response.ok) {
    const err: ApiError = await response.json()
    throw new Error(err.message)
  }

  return response.json()
}
