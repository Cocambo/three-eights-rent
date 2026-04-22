import { API_GATEWAY_BASE_URL, type ApiRequestOptions } from '@/services/api'
import type { CarCatalogCardModel } from '@/services/cars'

export interface RecommendationItem {
  id: number
  brand: string
  model: string
  year: number
  fuel_type: string
  transmission: string
  body_type: string
  seats_count: number
  price_per_day: number
  purpose: string
  main_image_url: string | null
  score?: number
}

export interface RecommendationsResponse {
  items: RecommendationItem[]
}

type AuthorizedRequest = <TResponse, TBody = unknown>(
  path: string,
  options?: Omit<ApiRequestOptions<TBody>, 'accessToken'>,
) => Promise<TResponse>

export function mapRecommendationItemToCardModel(
  item: RecommendationItem,
): CarCatalogCardModel {
  return {
    id: item.id,
    brand: item.brand,
    model: item.model,
    name: `${item.brand} ${item.model}`,
    year: item.year,
    fuelType: item.fuel_type,
    transmission: item.transmission,
    bodyType: item.body_type,
    seatsCount: item.seats_count,
    pricePerDay: item.price_per_day,
    purpose: item.purpose,
    imageUrl: item.main_image_url,
  }
}

export function getPersonalRecommendations(
  request: AuthorizedRequest,
  signal?: AbortSignal,
) {
  return request<RecommendationsResponse>('/recommendations/me', {
    signal,
    baseUrl: API_GATEWAY_BASE_URL,
  })
}
