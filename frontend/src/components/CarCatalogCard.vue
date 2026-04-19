<template>
  <article class="car-card">
    <div class="car-card__image-wrap">
      <img
        v-if="car.imageUrl"
        :src="car.imageUrl"
        :alt="car.name"
        class="car-card__image"
      />
      <div v-else class="car-card__image-placeholder">
        <span class="material-symbols-outlined">image</span>
        <span>Изображение недоступно</span>
      </div>

      <button
        class="favorite-button"
        type="button"
        :aria-label="favoriteButtonLabel"
        :aria-pressed="isFavorite"
        :disabled="favoritesStore.isPending(car.id)"
        @click="handleFavoriteClick"
      >
        <span
          class="material-symbols-outlined"
          :class="{ 'favorite-button__icon--active': isFavorite }"
        >
          favorite
        </span>
      </button>
    </div>

    <div class="car-card__body">
      <div class="car-card__top">
        <div>
          <h3>{{ car.name }}</h3>
          <p>{{ car.bodyType }}</p>
        </div>
        <div class="car-card__meta">
          <div class="car-card__price">
            <strong>{{ formatPrice(car.pricePerDay) }}</strong>
            <span>в сутки</span>
          </div>
          <div v-if="car.purpose" class="car-card__badges car-card__badges--stacked">
            <span class="badge">
              {{ car.purpose }}
            </span>
          </div>
        </div>
      </div>

      <div class="car-card__specs">
        <div class="spec-item">
          <span class="material-symbols-outlined">group</span>
          <span>{{ car.seatsCount }}</span>
        </div>
        <div class="spec-item">
          <span class="material-symbols-outlined">local_gas_station</span>
          <span>{{ car.fuelType }}</span>
        </div>
        <div class="spec-item">
          <span class="material-symbols-outlined">settings_suggest</span>
          <span>{{ car.transmission }}</span>
        </div>
        <div class="spec-item">
          <span class="material-symbols-outlined">calendar_today</span>
          <span>{{ car.year }}</span>
        </div>
      </div>

      <RouterLink class="secondary-button" :to="{ name: 'car-details', params: { id: car.id } }">
        Подробнее
      </RouterLink>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import type { CarCatalogCardModel, CarCatalogItem } from '@/services/cars'
import { useAuthStore } from '@/stores/auth'
import { useFavoritesStore } from '@/stores/favorites'

const props = defineProps<{
  car: CarCatalogCardModel
}>()

const authStore = useAuthStore()
const favoritesStore = useFavoritesStore()
const router = useRouter()
const route = useRoute()

const isFavorite = computed(() => favoritesStore.isFavorite(props.car.id))
const favoriteButtonLabel = computed(() =>
  isFavorite.value ? 'Удалить из избранного' : 'Добавить в избранное',
)

function toCatalogItem(car: CarCatalogCardModel): CarCatalogItem {
  return {
    id: car.id,
    brand: car.brand,
    model: car.model,
    year: car.year,
    fuel_type: car.fuelType,
    transmission: car.transmission,
    body_type: car.bodyType,
    seats_count: car.seatsCount,
    price_per_day: car.pricePerDay,
    purpose: car.purpose,
    main_image_url: car.imageUrl,
  }
}

async function handleFavoriteClick() {
  if (!authStore.isAuthenticated) {
    await router.push({
      name: 'login',
      query: {
        redirect: route.fullPath,
      },
    })
    return
  }

  try {
    await favoritesStore.toggleFavorite(props.car.id, toCatalogItem(props.car))
  } catch {
    // The store keeps the last backend error message for the views that need it.
  }
}

function formatPrice(price: number) {
  return `${new Intl.NumberFormat('ru-RU').format(price)} ₽`
}
</script>

<style scoped>
.car-card {
  overflow: hidden;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(17, 29, 35, 0.08);
  box-shadow: 0 16px 40px rgba(14, 40, 64, 0.06);
}

.car-card__image-wrap {
  position: relative;
  height: 250px;
  overflow: hidden;
}

.car-card__image,
.car-card__image-placeholder {
  width: 100%;
  height: 100%;
}

.car-card__image {
  object-fit: cover;
  transition: transform 0.35s ease;
}

.car-card__image-placeholder {
  display: grid;
  place-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #eef4f8 0%, #d8e2eb 100%);
  color: #526170;
  text-align: center;
  padding: 24px;
}

.car-card__image-placeholder .material-symbols-outlined {
  font-size: 36px;
}

.car-card:hover .car-card__image {
  transform: scale(1.04);
}

.favorite-button {
  position: absolute;
  top: 16px;
  right: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 0;
  padding: 0;
  border-radius: 999px;
  background: rgba(0, 25, 68, 0.34);
  color: #fff;
  backdrop-filter: blur(10px);
  cursor: pointer;
  transition: 0.2s ease;
}

.favorite-button:hover,
.secondary-button:hover {
  transform: translateY(-1px);
}

.favorite-button:disabled {
  opacity: 0.72;
  cursor: not-allowed;
}

.favorite-button .material-symbols-outlined {
  font-size: 26px;
  font-variation-settings:
    'FILL' 0,
    'wght' 300,
    'GRAD' 0,
    'opsz' 24;
}

.favorite-button__icon--active {
  color: #ff7e8a;
  font-variation-settings:
    'FILL' 1,
    'wght' 500,
    'GRAD' 0,
    'opsz' 24;
}

.car-card__body {
  padding: 20px;
}

.car-card__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.car-card__top h3 {
  margin: 0;
}

.car-card__top p {
  margin: 8px 0 0;
  color: #617080;
}

.car-card__meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 12px;
  flex-shrink: 0;
}

.car-card__price {
  text-align: right;
}

.car-card__price strong {
  display: block;
  font-size: 22px;
  color: #11284b;
}

.car-card__price span {
  font-size: 12px;
  color: #617080;
}

.car-card__badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.car-card__badges--stacked {
  justify-content: flex-end;
  max-width: 220px;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 32px;
  padding: 0 14px;
  border-radius: 9999px;
  border: 1px solid rgba(22, 63, 119, 0.08);
  background: linear-gradient(135deg, #245ca2 0%, #163f77 100%);
  box-shadow: 0 10px 24px rgba(22, 63, 119, 0.16);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
  line-height: 1;
  text-transform: uppercase;
  white-space: nowrap;
  overflow: hidden;
}

.car-card__specs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin: 18px 0;
  padding: 14px 0;
  border-top: 1px solid #ebf0f4;
  border-bottom: 1px solid #ebf0f4;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  color: #526170;
  font-size: 14px;
}

.secondary-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 48px;
  border: 0;
  border-radius: 14px;
  background: #edf3f7;
  color: #163f77;
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
  transition: 0.2s ease;
}

@media (max-width: 768px) {
  .car-card__top,
  .car-card__specs {
    flex-direction: column;
    align-items: flex-start;
  }

  .car-card__meta,
  .car-card__price,
  .car-card__badges--stacked {
    align-items: flex-start;
    text-align: left;
    justify-content: flex-start;
  }
}
</style>
