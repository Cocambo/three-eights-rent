import { API_GATEWAY_BASE_URL, apiRequest } from '@/services/api'

export interface PaginationMeta {
  total: number
  limit: number
  offset: number
}

export interface CarsCatalogQuery {
  q?: string
  brand?: string
  model?: string
  year_from?: number
  year_to?: number
  fuel_type?: string
  transmission?: string
  body_type?: string
  seats_min?: number
  price_min?: number
  price_max?: number
  purpose?: string
  sort_by?: 'id' | 'year' | 'price_per_day' | 'created_at'
  sort_order?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export interface CarCatalogItem {
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
}

export interface CarsCatalogResponse {
  items: CarCatalogItem[]
  pagination: PaginationMeta
}

export interface CarImage {
  id: number
  url: string
  is_main: boolean
  sort_order: number
}

export interface CarDetailsResponse {
  id: number
  brand: string
  model: string
  year: number
  fuel_type: string
  transmission: string
  body_type: string
  color: string
  seats_count: number
  price_per_day: number
  purpose: string
  description: string
  images: CarImage[]
}

export interface CarCatalogCardModel {
  id: number
  brand: string
  model: string
  name: string
  year: number
  fuelType: string
  transmission: string
  bodyType: string
  seatsCount: number
  pricePerDay: number
  purpose: string
  imageUrl: string | null
}

export function mapCarCatalogItemToCardModel(item: CarCatalogItem): CarCatalogCardModel {
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

function buildCarsCatalogPath(query: CarsCatalogQuery) {
  const params = new URLSearchParams()

  const entries = Object.entries(query) as Array<[keyof CarsCatalogQuery, string | number | undefined]>
  for (const [key, value] of entries) {
    if (value === undefined || value === null || value === '') {
      continue
    }

    params.set(key, String(value))
  }

  const search = params.toString()
  return search ? `/cars?${search}` : '/cars'
}

export function getCarsCatalog(query: CarsCatalogQuery, signal?: AbortSignal) {
  return apiRequest<CarsCatalogResponse>(buildCarsCatalogPath(query), {
    signal,
    baseUrl: API_GATEWAY_BASE_URL,
  })
}

export function getCarDetails(carId: number, signal?: AbortSignal) {
  return apiRequest<CarDetailsResponse>(`/cars/${carId}`, {
    signal,
    baseUrl: API_GATEWAY_BASE_URL,
  })
}
