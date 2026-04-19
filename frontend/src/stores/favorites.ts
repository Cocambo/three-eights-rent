import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'

import type { CarCatalogItem, CarCatalogCardModel } from '@/services/cars'
import { mapCarCatalogItemToCardModel } from '@/services/cars'
import {
  addFavorite,
  getFavorites,
  removeFavorite,
  type FavoriteItem,
} from '@/services/favorites'
import { useAuthStore } from '@/stores/auth'

function addUniqueCarId(carIds: number[], carId: number) {
  return carIds.includes(carId) ? carIds : [...carIds, carId]
}

function removeCarId(carIds: number[], carId: number) {
  return carIds.filter((id) => id !== carId)
}

export const useFavoritesStore = defineStore('favorites', () => {
  const authStore = useAuthStore()

  const items = ref<FavoriteItem[]>([])
  const favoriteCarIds = ref<number[]>([])
  const pendingCarIds = ref<number[]>([])
  const isLoading = ref(false)
  const errorMessage = ref('')
  const hasLoaded = ref(false)

  const favoriteCards = computed<CarCatalogCardModel[]>(() =>
    items.value.map((item) => mapCarCatalogItemToCardModel(item.car)),
  )

  function reset() {
    items.value = []
    favoriteCarIds.value = []
    pendingCarIds.value = []
    isLoading.value = false
    errorMessage.value = ''
    hasLoaded.value = false
  }

  function setItems(nextItems: FavoriteItem[]) {
    items.value = nextItems
    favoriteCarIds.value = nextItems.map((item) => item.car_id)
    hasLoaded.value = true
  }

  function isFavorite(carId: number) {
    return favoriteCarIds.value.includes(carId)
  }

  function setPending(carId: number, value: boolean) {
    pendingCarIds.value = value
      ? addUniqueCarId(pendingCarIds.value, carId)
      : removeCarId(pendingCarIds.value, carId)
  }

  function isPending(carId: number) {
    return pendingCarIds.value.includes(carId)
  }

  function upsertItemFromCatalog(car: CarCatalogItem) {
    const nextItem: FavoriteItem = {
      car_id: car.id,
      added_at: new Date().toISOString(),
      car,
    }

    const existingIndex = items.value.findIndex((item) => item.car_id === car.id)
    if (existingIndex === -1) {
      items.value = [nextItem, ...items.value]
      return
    }

    const nextItems = [...items.value]
    nextItems[existingIndex] = nextItem
    items.value = nextItems
  }

  async function loadFavorites(options: { force?: boolean; signal?: AbortSignal } = {}) {
    if (!authStore.isAuthenticated) {
      reset()
      return []
    }

    if (hasLoaded.value && !options.force) {
      return items.value
    }

    isLoading.value = true
    errorMessage.value = ''

    try {
      const response = await getFavorites(authStore.authorizedRequest, options.signal)
      setItems(response.items)
      return items.value
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить избранное.'
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function ensureLoaded(signal?: AbortSignal) {
    return loadFavorites({ signal })
  }

  async function addToFavorites(carId: number, car?: CarCatalogItem) {
    setPending(carId, true)
    errorMessage.value = ''

    try {
      await addFavorite(authStore.authorizedRequest, carId)
      favoriteCarIds.value = addUniqueCarId(favoriteCarIds.value, carId)

      if (car) {
        upsertItemFromCatalog(car)
      } else if (hasLoaded.value) {
        await loadFavorites({ force: true })
      }
    } catch (error) {
      errorMessage.value =
        error instanceof Error ? error.message : 'Не удалось добавить автомобиль в избранное.'
      throw error
    } finally {
      setPending(carId, false)
    }
  }

  async function removeFromFavorites(carId: number) {
    setPending(carId, true)
    errorMessage.value = ''

    try {
      await removeFavorite(authStore.authorizedRequest, carId)
      favoriteCarIds.value = removeCarId(favoriteCarIds.value, carId)
      items.value = items.value.filter((item) => item.car_id !== carId)
    } catch (error) {
      errorMessage.value =
        error instanceof Error ? error.message : 'Не удалось удалить автомобиль из избранного.'
      throw error
    } finally {
      setPending(carId, false)
    }
  }

  async function toggleFavorite(carId: number, car?: CarCatalogItem) {
    if (isPending(carId)) {
      return
    }

    if (isFavorite(carId)) {
      await removeFromFavorites(carId)
      return
    }

    await addToFavorites(carId, car)
  }

  watch(
    () => authStore.isAuthenticated,
    (isAuthenticated) => {
      if (!isAuthenticated) {
        reset()
      }
    },
    { immediate: true },
  )

  return {
    items,
    favoriteCarIds,
    pendingCarIds,
    isLoading,
    errorMessage,
    hasLoaded,
    favoriteCards,
    reset,
    isFavorite,
    isPending,
    loadFavorites,
    ensureLoaded,
    addToFavorites,
    removeFromFavorites,
    toggleFavorite,
  }
})
