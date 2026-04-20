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

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

export interface ApiRequestOptions<TBody> {
  method?: HttpMethod
  body?: TBody
  accessToken?: string
  signal?: AbortSignal
  baseUrl?: string
}

const fallbackUserApiBaseUrl = '/api/v1/users'
const fallbackGatewayApiBaseUrl = '/api/v1'

function normalizeBaseUrl(baseUrl: string) {
  return baseUrl.replace(/\/$/, '')
}

export const USER_SERVICE_API_BASE_URL = normalizeBaseUrl(
  import.meta.env.VITE_API_BASE_URL || fallbackUserApiBaseUrl,
)

export const API_GATEWAY_BASE_URL = normalizeBaseUrl(
  import.meta.env.VITE_GATEWAY_API_BASE_URL || fallbackGatewayApiBaseUrl,
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
  return text ? text : null
}

function normalizeErrorPayload(data: unknown): { message: string; code?: string } {
  if (typeof data === 'string' && data.trim()) {
    return { message: data }
  }

  if (data && typeof data === 'object') {
    const payload = data as Record<string, unknown>
    const topLevelCode = typeof payload.code === 'string' ? payload.code : undefined

    if (typeof payload.error === 'string' && payload.error.trim()) {
      return {
        message: payload.error,
        code: topLevelCode,
      }
    }

    if (payload.error && typeof payload.error === 'object') {
      const nestedError = payload.error as Record<string, unknown>
      const nestedMessage =
        typeof nestedError.message === 'string' && nestedError.message.trim()
          ? nestedError.message
          : undefined
      const nestedCode =
        typeof nestedError.code === 'string' ? nestedError.code : topLevelCode

      if (nestedMessage) {
        return {
          message: nestedMessage,
          code: nestedCode,
        }
      }
    }

    if (typeof payload.message === 'string' && payload.message.trim()) {
      return {
        message: payload.message,
        code: topLevelCode,
      }
    }
  }

  return {
    message: 'Не удалось выполнить запрос.',
  }
}

export async function apiRequest<TResponse, TBody = unknown>(
  path: string,
  options: ApiRequestOptions<TBody> = {},
): Promise<TResponse> {
  const baseUrl = options.baseUrl || USER_SERVICE_API_BASE_URL

  const response = await fetch(`${baseUrl}${path}`, {
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
    const payload = normalizeErrorPayload(data)
    throw new ApiError(payload.message, response.status, payload.code)
  }

  return data as TResponse
}
