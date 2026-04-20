import type { ApiRequestOptions } from '@/services/api'
import type { CarCatalogItem } from '@/services/cars'

export interface FavoriteItem {
  car_id: number
  added_at: string
  car: CarCatalogItem
}

export interface FavoritesResponse {
  items: FavoriteItem[]
}

export interface FavoriteMutationResponse {
  car_id: number
  message: string
}

type AuthorizedRequest = <TResponse, TBody = unknown>(
  path: string,
  options?: Omit<ApiRequestOptions<TBody>, 'accessToken'>,
) => Promise<TResponse>

const favoritesBaseUrl = '/api/v1'

export function getFavorites(request: AuthorizedRequest, signal?: AbortSignal) {
  return request<FavoritesResponse>('/favorites', {
    signal,
    baseUrl: favoritesBaseUrl,
  })
}

export function addFavorite(request: AuthorizedRequest, carId: number) {
  return request<FavoriteMutationResponse>(`/favorites/${carId}`, {
    method: 'POST',
    baseUrl: favoritesBaseUrl,
  })
}

export function removeFavorite(request: AuthorizedRequest, carId: number) {
  return request<FavoriteMutationResponse>(`/favorites/${carId}`, {
    method: 'DELETE',
    baseUrl: favoritesBaseUrl,
  })
}
