export class ApiError extends Error {
  status: number
  code?: string

  constructor(message: string, status: number, code?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

type HttpMethod = 'GET' | 'POST' | 'PUT'

interface RequestOptions<TBody> {
  method?: HttpMethod
  body?: TBody
  accessToken?: string
  signal?: AbortSignal
}

interface ErrorPayload {
  error?: string
  code?: string
}

const fallbackBaseUrl = '/api/v1/users'

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL || fallbackBaseUrl).replace(
  /\/$/,
  '',
)

function buildHeaders(accessToken?: string): HeadersInit {
  const headers: HeadersInit = {
    Accept: 'application/json',
  }

  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`
  }

  return headers
}

async function parseResponse(response: Response): Promise<unknown> {
  if (response.status === 204) {
    return null
  }

  const contentType = response.headers.get('content-type') || ''

  if (contentType.includes('application/json')) {
    return response.json()
  }

  const text = await response.text()
  return text ? { error: text } : null
}

export async function apiRequest<TResponse, TBody = unknown>(
  path: string,
  options: RequestOptions<TBody> = {},
): Promise<TResponse> {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    method: options.method || 'GET',
    headers: {
      ...buildHeaders(options.accessToken),
      ...(options.body ? { 'Content-Type': 'application/json' } : {}),
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
    signal: options.signal,
  })

  const data = await parseResponse(response)

  if (!response.ok) {
    const payload = (data || {}) as ErrorPayload
    throw new ApiError(
      payload.error || 'Не удалось выполнить запрос.',
      response.status,
      payload.code,
    )
  }

  return data as TResponse
}
