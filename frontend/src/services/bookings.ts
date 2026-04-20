import { API_GATEWAY_BASE_URL, ApiError, type ApiRequestOptions } from '@/services/api'
import type { CarCatalogItem } from '@/services/cars'

export interface BookingRecord {
  id: number
  car_id: number
  start_date: string
  end_date: string
  status: string
  created_at: string
  updated_at: string
  cancelled_at?: string | null
}

export interface BookingHistoryItem extends BookingRecord {
  car: CarCatalogItem
}

export interface BookingsResponse {
  items: BookingHistoryItem[]
}

export interface CancelBookingResponse {
  booking_id: number
  status: string
  message: string
}

type AuthorizedRequest = <TResponse, TBody = unknown>(
  path: string,
  options?: Omit<ApiRequestOptions<TBody>, 'accessToken'>,
) => Promise<TResponse>

interface CreateBookingPayload {
  car_id: number
  start_date: string
  end_date: string
}

export function createBooking(
  request: AuthorizedRequest,
  payload: CreateBookingPayload,
  signal?: AbortSignal,
) {
  return request<BookingRecord, CreateBookingPayload>('/bookings', {
    method: 'POST',
    body: payload,
    signal,
    baseUrl: API_GATEWAY_BASE_URL,
  })
}

export function getBookings(request: AuthorizedRequest, signal?: AbortSignal) {
  return request<BookingsResponse>('/bookings', {
    signal,
    baseUrl: API_GATEWAY_BASE_URL,
  })
}

export function cancelBooking(request: AuthorizedRequest, bookingId: number) {
  return request<CancelBookingResponse>(`/bookings/${bookingId}`, {
    method: 'DELETE',
    baseUrl: API_GATEWAY_BASE_URL,
  })
}

export function getBookingErrorMessage(
  error: unknown,
  fallback = 'Не удалось выполнить операцию с бронированием.',
) {
  if (!(error instanceof ApiError)) {
    return error instanceof Error ? error.message : fallback
  }

  const normalizedMessage = error.message.toLowerCase()

  if (normalizedMessage.includes('booking dates overlap')) {
    return 'Эти даты уже заняты. Выберите другой период бронирования.'
  }

  if (normalizedMessage.includes('start_date must be earlier than end_date')) {
    return 'Дата завершения должна быть позже даты начала аренды.'
  }

  if (normalizedMessage.includes('car not found')) {
    return 'Автомобиль не найден. Попробуйте открыть карточку заново.'
  }

  if (normalizedMessage.includes('completed booking cannot be cancelled')) {
    return 'Завершённое бронирование нельзя отменить.'
  }

  if (normalizedMessage.includes('booking not found')) {
    return 'Бронирование не найдено.'
  }

  return error.message || fallback
}
