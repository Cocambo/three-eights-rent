<template>
  <article class="car-card">
    <div class="car-card__image-wrap">
      <img :src="car.image" :alt="car.name" class="car-card__image" />

      <button class="favorite-button" type="button" aria-label="Добавить в избранное">
        <span class="material-symbols-outlined">
          {{ car.favorite ? 'favorite' : 'favorite' }}
        </span>
      </button>
    </div>

    <div class="car-card__body">
      <div class="car-card__top">
        <div>
          <h3>{{ car.name }}</h3>
          <p>{{ car.category }}</p>
        </div>
        <div class="car-card__meta">
          <div class="car-card__price">
            <strong>{{ formatPrice(car.pricePerDay) }}</strong>
            <span>в сутки</span>
          </div>
          <div v-if="car.tags.length" class="car-card__badges car-card__badges--stacked">
            <span
              v-for="badge in car.tags"
              :key="badge"
              class="badge"
              :class="{
                'badge--premium': badge.toLowerCase() === 'premium',
              }"
            >
              {{ badge }}
            </span>
          </div>
        </div>
      </div>

      <div class="car-card__specs">
        <div class="spec-item">
          <span class="material-symbols-outlined">group</span>
          <span>{{ car.seats }}</span>
        </div>
        <div class="spec-item">
          <span class="material-symbols-outlined">local_gas_station</span>
          <span>{{ car.catalogFuelType }}</span>
        </div>
        <div class="spec-item">
          <span class="material-symbols-outlined">settings_suggest</span>
          <span>{{ car.catalogTransmission }}</span>
        </div>
      </div>

      <RouterLink class="secondary-button" :to="{ name: 'car-details', params: { id: car.id } }">
        Подробнее
      </RouterLink>
    </div>
  </article>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

import type { Car } from '@/data/cars'

type CatalogCar = Car & {
  catalogFuelType: string
  catalogTransmission: string
}

defineProps<{
  car: CatalogCar
}>()

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

.car-card__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.35s ease;
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
  width: 32px;
  height: 32px;
  border: 0;
  padding: 0;
  background: transparent;
  color: #fff;
  text-shadow: 0 8px 18px rgba(0, 0, 0, 0.35);
  cursor: pointer;
  transition: 0.2s ease;
}

.favorite-button:hover,
.secondary-button:hover {
  transform: translateY(-1px);
}

.favorite-button .material-symbols-outlined {
  font-size: 28px;
  font-variation-settings:
    'FILL' 0,
    'wght' 300,
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

.badge--premium {
  background: linear-gradient(135deg, #2f74c4 0%, #194987 100%);
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
