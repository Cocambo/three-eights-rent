<template>
  <section class="recommendations-section" aria-labelledby="home-recommendations-title">
    <div class="container">
      <div class="recommendations-shell">
        <div class="recommendations-heading">
          <div>
            <span class="recommendations-heading__eyebrow" id="home-recommendations-title">
              Специально для вас
            </span>
            <p>
              Подобрали автомобили, с учётом ваших поездок и предпочтений.
            </p>
          </div>

          <RouterLink class="recommendations-heading__link" :to="{ name: 'cars' }">
            <span>Смотреть весь каталог</span>
            <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span>
          </RouterLink>
        </div>

        <div v-if="isLoading" class="recommendations-state recommendations-state--loading">
          <article v-for="index in loadingCardCount" :key="index" class="recommendations-skeleton">
            <div class="recommendations-skeleton__image"></div>
            <div class="recommendations-skeleton__line recommendations-skeleton__line--title"></div>
            <div class="recommendations-skeleton__line"></div>
            <div class="recommendations-skeleton__line recommendations-skeleton__line--short"></div>
          </article>
        </div>

        <div v-else-if="errorMessage" class="recommendations-state recommendations-state--error">
          <div class="recommendations-state__icon">
            <span class="material-symbols-outlined" aria-hidden="true">wifi_off</span>
          </div>
          <div>
            <h3>Не удалось загрузить рекомендации</h3>
            <p>{{ errorMessage }}</p>
          </div>
          <button class="recommendations-state__button" type="button" @click="loadRecommendations">
            Повторить
          </button>
        </div>

        <div v-else-if="!recommendationCards.length" class="recommendations-state">
          <div class="recommendations-state__icon">
            <span class="material-symbols-outlined" aria-hidden="true">directions_car</span>
          </div>
          <div>
            <h3>Рекомендации пока недоступны</h3>
            <p>Сервис пока вернул пустую подборку. Попробуйте обновить страницу чуть позже.</p>
          </div>
        </div>

        <div v-else class="recommendations-content">
          <div class="recommendations-content__meta">
          </div>

          <CarCatalogList :cars="recommendationCards" />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import CarCatalogList from '@/components/CarCatalogList.vue'
import {
  getPersonalRecommendations,
  mapRecommendationItemToCardModel,
  type RecommendationItem,
} from '@/services/recommendations'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const loadingCardCount = 3
const items = ref<RecommendationItem[]>([])
const isLoading = ref(true)
const errorMessage = ref('')

let currentController: AbortController | null = null

const recommendationCards = computed(() =>
  items.value.map((item) => mapRecommendationItemToCardModel(item)),
)

function formatRecommendationCount(count: number) {
  const lastTwoDigits = count % 100
  const lastDigit = count % 10

  if (lastTwoDigits >= 11 && lastTwoDigits <= 14) {
    return `${count} автомобилей`
  }

  if (lastDigit === 1) {
    return `${count} автомобиль`
  }

  if (lastDigit >= 2 && lastDigit <= 4) {
    return `${count} автомобиля`
  }

  return `${count} автомобилей`
}

async function loadRecommendations() {
  currentController?.abort()
  currentController = new AbortController()
  isLoading.value = true
  errorMessage.value = ''

  try {
    const response = await getPersonalRecommendations(
      authStore.authorizedRequest,
      currentController.signal,
    )
    items.value = response.items
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      return
    }

    items.value = []
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось получить персональные рекомендации.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  void loadRecommendations()
})

onBeforeUnmount(() => {
  currentController?.abort()
})
</script>

<style scoped>
.recommendations-section {
  padding: 28px 0 8px;
}

.container {
  width: min(1200px, calc(100% - 32px));
  margin: 0 auto;
}

.recommendations-shell {
  position: relative;
  padding: 28px;
  border-radius: 32px;
  background:
    radial-gradient(circle at top right, rgba(105, 255, 135, 0.18), transparent 30%),
    radial-gradient(circle at left center, rgba(176, 198, 255, 0.34), transparent 34%),
    linear-gradient(135deg, rgba(255, 255, 255, 0.92), rgba(236, 244, 249, 0.9));
  border: 1px solid rgba(0, 25, 68, 0.08);
  box-shadow: 0 30px 60px rgba(14, 40, 64, 0.08);
  overflow: hidden;
}

.recommendations-shell::before {
  content: '';
  position: absolute;
  inset: 14px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.56);
  pointer-events: none;
}

.recommendations-heading,
.recommendations-state {
  position: relative;
  z-index: 1;
}

.recommendations-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
}

.recommendations-heading__eyebrow {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  margin-bottom: 10px;
  padding: 0 14px;
  border-radius: 999px;
  background: rgba(0, 88, 202, 0.08);
  color: #0058ca;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  white-space: nowrap;
}

.recommendations-state h3 {
  margin: 0;
  color: #001944;
  font-family: 'Plus Jakarta Sans', sans-serif;
}

.recommendations-heading p,
.recommendations-state p,
.recommendations-content__meta p {
  max-width: 640px;
  margin: 0;
  color: #556372;
  font-size: 15px;
  line-height: 1.7;
}

.recommendations-heading__link,
.recommendations-state__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
  gap: 10px;
  padding: 0 18px 0 20px;
  border: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
  color: #001944;
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease,
    color 0.2s ease;
  box-shadow:
    inset 0 0 0 1px rgba(0, 25, 68, 0.08),
    0 10px 24px rgba(14, 40, 64, 0.08);
}

.recommendations-heading__link:hover,
.recommendations-state__button:hover {
  transform: translateY(-1px);
  background: #001944;
  color: #ffffff;
  box-shadow: 0 16px 30px rgba(0, 25, 68, 0.16);
}

.recommendations-heading__link .material-symbols-outlined {
  font-size: 18px;
}

.recommendations-content {
  position: relative;
  z-index: 1;
}

.recommendations-content__meta {
  margin-bottom: 18px;
}

.recommendations-content__meta p {
  margin-top: 0;
  font-size: 14px;
}

.recommendations-content :deep(.cars-grid) {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.recommendations-content :deep(.car-card) {
  border-radius: 20px;
  box-shadow: 0 14px 30px rgba(14, 40, 64, 0.05);
}

.recommendations-content :deep(.car-card__image-wrap) {
  height: 180px;
}

.recommendations-content :deep(.favorite-button) {
  top: 12px;
  right: 12px;
  width: 36px;
  height: 36px;
}

.recommendations-content :deep(.favorite-button .material-symbols-outlined) {
  font-size: 22px;
}

.recommendations-content :deep(.car-card__body) {
  padding: 16px;
}

.recommendations-content :deep(.car-card__top) {
  gap: 12px;
}

.recommendations-content :deep(.car-card__top h3) {
  font-size: 20px;
  line-height: 1.1;
}

.recommendations-content :deep(.car-card__top p) {
  margin-top: 6px;
  font-size: 13px;
}

.recommendations-content :deep(.car-card__price strong) {
  font-size: 20px;
}

.recommendations-content :deep(.car-card__price span) {
  font-size: 11px;
}

.recommendations-content :deep(.badge) {
  min-height: 28px;
  padding: 0 12px;
  font-size: 10px;
}

.recommendations-content :deep(.car-card__specs) {
  gap: 12px;
  margin: 14px 0;
  padding: 12px 0;
}

.recommendations-content :deep(.spec-item) {
  font-size: 13px;
}

.recommendations-content :deep(.secondary-button) {
  height: 42px;
  border-radius: 12px;
  font-size: 14px;
}

.recommendations-state {
  display: grid;
  gap: 18px;
}

.recommendations-state--loading {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.recommendations-state--error,
.recommendations-state:not(.recommendations-state--loading) {
  align-items: center;
  justify-items: center;
  padding: 34px 24px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.68);
  border: 1px solid rgba(0, 25, 68, 0.08);
  text-align: center;
}

.recommendations-state--error {
  border-color: rgba(186, 26, 26, 0.16);
}

.recommendations-state__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 20px;
  background: rgba(0, 25, 68, 0.06);
  color: #001944;
}

.recommendations-state__icon .material-symbols-outlined {
  font-size: 30px;
}

.recommendations-skeleton {
  min-height: 360px;
  padding: 18px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid rgba(0, 25, 68, 0.06);
}

.recommendations-skeleton__image,
.recommendations-skeleton__line {
  border-radius: 18px;
  background: linear-gradient(
    90deg,
    rgba(219, 229, 239, 0.7) 0%,
    rgba(243, 247, 251, 0.95) 50%,
    rgba(219, 229, 239, 0.7) 100%
  );
  background-size: 200% 100%;
  animation: recommendations-shimmer 1.4s linear infinite;
}

.recommendations-skeleton__image {
  height: 210px;
  margin-bottom: 18px;
}

.recommendations-skeleton__line {
  height: 16px;
  margin-bottom: 12px;
}

.recommendations-skeleton__line--title {
  height: 24px;
  width: 72%;
}

.recommendations-skeleton__line--short {
  width: 44%;
  margin-bottom: 0;
}

@keyframes recommendations-shimmer {
  from {
    background-position: 200% 0;
  }

  to {
    background-position: -200% 0;
  }
}

@media (max-width: 1100px) {
  .recommendations-state--loading {
    grid-template-columns: 1fr;
  }

  .recommendations-content :deep(.cars-grid) {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .container {
    width: min(100% - 24px, 1200px);
  }

  .recommendations-section {
    padding-top: 20px;
  }

  .recommendations-shell {
    padding: 22px;
    border-radius: 26px;
  }

  .recommendations-heading {
    align-items: stretch;
    flex-direction: column;
  }

  .recommendations-heading__link,
  .recommendations-state__button {
    width: 100%;
  }

  .recommendations-content :deep(.cars-grid) {
    grid-template-columns: 1fr;
  }

  .recommendations-content :deep(.car-card__image-wrap) {
    height: 210px;
  }
}
</style>
