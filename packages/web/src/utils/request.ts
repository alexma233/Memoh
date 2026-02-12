import router from '@/router'

const BASE_URL = '/api'

export class ApiError extends Error {
  status: number
  statusText: string
  body?: unknown

  constructor(status: number, statusText: string, body?: unknown) {
    super(`API Error ${status}: ${statusText}`)
    this.name = 'ApiError'
    this.status = status
    this.statusText = statusText
    this.body = body
  }
}

export interface FetchOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  body?: unknown
  headers?: Record<string, string>
  /** Do not attach Authorization header */
  noAuth?: boolean
  signal?: AbortSignal
}

/**
 * Type-safe fetch wrapper: JSON body, token injection, 401 redirect.
 * Returns the response JSON directly (no .data unwrap).
 */
export async function fetchApi<T = unknown>(
  url: string,
  options: FetchOptions = {},
): Promise<T> {
  const { method = 'GET', body, headers = {}, noAuth = false, signal } = options

  if (!noAuth) {
    const token = localStorage.getItem('token')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }

  const response = await fetch(`${BASE_URL}${url}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    signal,
  })

  if (response.status === 401) {
    localStorage.removeItem('token')
    router.replace({ name: 'Login' })
    throw new ApiError(response.status, response.statusText)
  }

  if (!response.ok) {
    let errorBody: unknown
    try {
      errorBody = await response.json()
    } catch {
      // response body not JSON
    }
    throw new ApiError(response.status, response.statusText, errorBody)
  }

  if (response.status === 204 || response.headers.get('content-length') === '0') {
    return undefined as T
  }

  return response.json() as Promise<T>
}
