<template>
  <section class="car-page">
    <AppHeader />

    <main class="page-shell">
      <section v-if="car" class="car-details">
        <div class="breadcrumbs">
          <RouterLink :to="{ name: 'cars' }">Автомобили</RouterLink>
          <span>/</span>
          <span>{{ car.name }}</span>
        </div>

        <div class="content-grid">
          <div class="content-main">
            <section class="gallery">
              <button class="gallery__hero" type="button" @click="openLightbox(selectedImageIndex)">
                <img :src="selectedImage" :alt="car.name" />
              </button>

              <div class="gallery__thumbs">
                <button
                  v-for="(image, index) in car.images"
                  :key="`${car.id}-${index}`"
                  class="gallery__thumb"
                  :class="{ 'gallery__thumb--active': index === selectedImageIndex }"
                  type="button"
                  @click="selectedImageIndex = index"
                >
                  <img :src="image" :alt="`${car.name} фото ${index + 1}`" />
                </button>
              </div>
            </section>

            <header class="hero-copy">
              <h1>{{ car.name }}</h1>
              <p class="hero-copy__text">{{ car.description }}</p>

              <div class="hero-copy__chips">
                <span class="info-chip">{{ car.purpose }}</span>
              </div>
            </header>

            <section class="spec-grid">
              <article v-for="feature in primaryFeatures" :key="feature.label" class="spec-card">
                <span class="material-symbols-outlined spec-card__icon">
                  {{ feature.icon }}
                </span>
                <p>{{ feature.label }}</p>
                <strong>{{ feature.value }}</strong>
              </article>
            </section>

            <section class="details-panel">
              <div class="details-panel__header">
                <h2>Характеристики</h2>
              </div>

              <div class="details-table">
                <div class="details-row">
                  <span>Марка</span>
                  <strong>{{ car.brand }}</strong>
                </div>
                <div class="details-row">
                  <span>Модель</span>
                  <strong>{{ car.name }}</strong>
                </div>
                <div class="details-row">
                  <span>Кузов</span>
                  <strong>{{ car.bodyType }}</strong>
                </div>
                <div class="details-row">
                  <span>Цвет</span>
                  <strong>{{ car.color }}</strong>
                </div>
                <div class="details-row">
                  <span>Год выпуска</span>
                  <strong>{{ car.year }}</strong>
                </div>
                <div class="details-row">
                  <span>Трансмиссия</span>
                  <strong>{{ car.transmission }}</strong>
                </div>
                <div class="details-row">
                  <span>Количество мест</span>
                  <strong>{{ car.seats }}</strong>
                </div>
                <div class="details-row">
                  <span>Тип топлива</span>
                  <strong>{{ car.fuelType }}</strong>
                </div>
              </div>
            </section>
          </div>

          <aside class="booking-card">
            <div class="booking-card__top">
              <p>Стоимость аренды</p>
              <div>
                <strong>{{ formatPrice(car.pricePerDay) }}</strong>
                <span>/ сутки</span>
              </div>
            </div>

            <div class="booking-card__fields">
              <label class="booking-field">
                <span>Дата начала</span>
                <input v-model="bookingForm.startDate" :min="today" type="date" />
              </label>

              <label class="booking-field">
                <span>Дата завершения</span>
                <input v-model="bookingForm.endDate" :min="bookingForm.startDate || today" type="date" />
              </label>
            </div>

            <div class="availability">
              <p class="availability__title">Недоступные даты</p>
              <div class="availability__list">
                <span v-for="range in car.unavailableRanges" :key="`${range.start}-${range.end}`">
                  {{ formatRange(range.start, range.end) }}
                </span>
              </div>
            </div>

            <div class="booking-summary" :class="bookingState.toneClass">
              <template v-if="bookingState.ready">
                <p class="booking-summary__title">{{ bookingState.title }}</p>
                <p class="booking-summary__text">{{ bookingState.description }}</p>
                <div class="booking-summary__price">
                  <span>{{ bookingState.days }} суток</span>
                  <strong>{{ formatPrice(bookingState.totalPrice) }}</strong>
                </div>
              </template>
              <template v-else>
                <p class="booking-summary__title">{{ bookingState.title }}</p>
                <p class="booking-summary__text">{{ bookingState.description }}</p>
              </template>
            </div>

            <button
              class="booking-button"
              type="button"
              :disabled="!bookingState.ready"
              @click="submitBooking"
            >
              Забронировать
            </button>

            <p class="booking-note">
              <span class="material-symbols-outlined">verified_user</span>
              {{ car.bookingNotice }}
            </p>

            <div v-if="bookingResult" class="booking-result">
              <h3>Заглушка бронирования</h3>
              <p>{{ bookingResult }}</p>
            </div>

          </aside>
        </div>
      </section>

      <section v-else class="not-found">
        <h1>Автомобиль не найден</h1>
        <p>Возможно, карточка была удалена или ссылка устарела.</p>
        <RouterLink class="not-found__link" :to="{ name: 'cars' }">Вернуться в каталог</RouterLink>
      </section>
    </main>

    <AppFooter />

    <div v-if="car && lightboxOpen" class="lightbox" @click.self="closeLightbox">
      <button class="lightbox__close" type="button" @click="closeLightbox">
        <span class="material-symbols-outlined">close</span>
      </button>
      <button class="lightbox__nav lightbox__nav--prev" type="button" @click="prevImage">
        <span class="material-symbols-outlined">chevron_left</span>
      </button>
      <img class="lightbox__image" :src="selectedImage" :alt="car.name" />
      <button class="lightbox__nav lightbox__nav--next" type="button" @click="nextImage">
        <span class="material-symbols-outlined">chevron_right</span>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
import { findCarById } from '@/data/cars'

const route = useRoute()

const car = computed(() => {
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  return findCarById(id)
})

const selectedImageIndex = ref(0)
const lightboxOpen = ref(false)
const bookingResult = ref('')

const bookingForm = reactive({
  startDate: '',
  endDate: '',
})

const today = getToday()

const selectedImage = computed(() => car.value?.images[selectedImageIndex.value] ?? '')

const primaryFeatures = computed(() => {
  if (!car.value) return []

  return [
    { label: 'Год выпуска', value: String(car.value.year), icon: 'calendar_today' },
    { label: 'Топливо', value: car.value.fuelType, icon: 'local_gas_station' },
    { label: 'Трансмиссия', value: car.value.transmission, icon: 'settings' },
    { label: 'Мест', value: String(car.value.seats), icon: 'group' },
  ]
})

const bookingState = computed(() => {
  if (!car.value) {
    return {
      ready: false,
      title: 'Карточка недоступна',
      description: 'Автомобиль не найден.',
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--muted',
    }
  }

  if (!bookingForm.startDate || !bookingForm.endDate) {
    return {
      ready: false,
      title: 'Выберите даты аренды',
      description: 'Укажите дату начала и дату завершения, чтобы увидеть стоимость и доступность.',
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--muted',
    }
  }

  if (bookingForm.endDate <= bookingForm.startDate) {
    return {
      ready: false,
      title: 'Проверьте диапазон',
      description: 'Дата завершения должна быть позже даты начала аренды.',
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--error',
    }
  }

  const overlap = car.value.unavailableRanges.find((range) =>
    rangesOverlap(bookingForm.startDate, bookingForm.endDate, range.start, range.end),
  )

  if (overlap) {
    return {
      ready: false,
      title: 'Диапазон занят',
      description: `Выбранные даты пересекаются с периодом ${formatRange(overlap.start, overlap.end)}.`,
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--error',
    }
  }

  const days = getDaysBetween(bookingForm.startDate, bookingForm.endDate)

  return {
    ready: true,
    title: 'Автомобиль свободен',
    description: `Предварительный расчет для аренды с ${formatDate(bookingForm.startDate)} по ${formatDate(bookingForm.endDate)}.`,
    days,
    totalPrice: days * car.value.pricePerDay,
    toneClass: 'booking-summary--success',
  }
})

watch(
  () => route.params.id,
  () => {
    selectedImageIndex.value = 0
    lightboxOpen.value = false
    bookingResult.value = ''
    bookingForm.startDate = ''
    bookingForm.endDate = ''
  },
)

function openLightbox(index: number) {
  selectedImageIndex.value = index
  lightboxOpen.value = true
}

function closeLightbox() {
  lightboxOpen.value = false
}

function prevImage() {
  if (!car.value) return
  selectedImageIndex.value =
    (selectedImageIndex.value - 1 + car.value.images.length) % car.value.images.length
}

function nextImage() {
  if (!car.value) return
  selectedImageIndex.value = (selectedImageIndex.value + 1) % car.value.images.length
}

function submitBooking() {
  if (!car.value || !bookingState.value.ready) return

  bookingResult.value =
    `Заявка-заглушка создана: ${car.value.name}, ${formatDate(bookingForm.startDate)} - ${formatDate(bookingForm.endDate)}, ` +
    `${bookingState.value.days} суток, сумма ${formatPrice(bookingState.value.totalPrice)}. Позже сюда можно подключить API бронирования.`
}

function formatPrice(price: number) {
  return `${new Intl.NumberFormat('ru-RU').format(price)} ₽`
}

function formatDate(dateString: string) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(new Date(`${dateString}T00:00:00`))
}

function formatRange(start: string, end: string) {
  return `${formatDate(start)} - ${formatDate(end)}`
}

function rangesOverlap(startA: string, endA: string, startB: string, endB: string) {
  return startA < endB && endA > startB
}

function getDaysBetween(start: string, end: string) {
  const startDate = new Date(`${start}T00:00:00`)
  const endDate = new Date(`${end}T00:00:00`)
  const millisecondsPerDay = 1000 * 60 * 60 * 24

  return Math.ceil((endDate.getTime() - startDate.getTime()) / millisecondsPerDay)
}

function getToday() {
  const now = new Date()
  const offset = now.getTimezoneOffset()
  const normalized = new Date(now.getTime() - offset * 60 * 1000)

  return normalized.toISOString().slice(0, 10)
}
</script>

<style scoped>
.car-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top right, rgba(176, 198, 255, 0.28), transparent 25%),
    radial-gradient(circle at bottom left, rgba(149, 158, 253, 0.12), transparent 24%),
    #f4faff;
}

.page-shell {
  width: min(1280px, calc(100% - 32px));
  margin: 0 auto;
  padding: 24px 0 48px;
}

.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 24px;
  color: #667788;
  font-size: 14px;
}

.breadcrumbs a {
  color: #163f77;
  text-decoration: none;
}

.gallery {
  display: grid;
  gap: 18px;
}

.gallery__hero {
  position: relative;
  height: 420px;
  min-height: 280px;
  overflow: hidden;
  border: 0;
  border-radius: 28px;
  padding: 0;
  background: #dfe8ef;
  cursor: pointer;
}

.gallery__hero img,
.gallery__thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gallery__thumbs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.gallery__thumb {
  height: 148px;
  overflow: hidden;
  padding: 0;
  border: 2px solid transparent;
  border-radius: 22px;
  background: #dfe8ef;
  cursor: pointer;
}

.gallery__thumb--active {
  border-color: #163f77;
  box-shadow: 0 0 0 6px rgba(22, 63, 119, 0.1);
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(320px, 0.9fr);
  gap: 28px;
  align-items: start;
}

.content-main {
  display: grid;
  gap: 24px;
}

.hero-copy,
.details-panel,
.booking-card,
.not-found {
  border: 1px solid rgba(17, 29, 35, 0.08);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 20px 50px rgba(14, 40, 64, 0.06);
}

.hero-copy {
  padding: 34px;
  border-radius: 30px;
}

.info-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 38px;
  padding: 0 16px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 800;
}

.hero-copy h1,
.details-panel h2,
.not-found h1,
.booking-result h3 {
  margin: 0;
}

.hero-copy h1 {
  color: #001944;
  font-size: clamp(2.6rem, 4vw, 4.4rem);
  line-height: 0.94;
  letter-spacing: -0.06em;
}

.hero-copy__text {
  margin: 18px 0 0;
  max-width: 760px;
  color: #465665;
  font-size: 1.05rem;
  line-height: 1.8;
}

.hero-copy__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 22px;
}

.info-chip {
  background: #edf4f9;
  color: #163f77;
}

.spec-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.spec-card {
  padding: 22px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(17, 29, 35, 0.08);
}

.spec-card__icon {
  color: #163f77;
}

.spec-card p,
.details-panel__header p,
.booking-card__top p,
.booking-note,
.booking-result p,
.not-found p,
.availability__title {
  margin: 0;
  color: #60707f;
}

.spec-card p {
  margin-top: 18px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.spec-card strong {
  display: block;
  margin-top: 10px;
  color: #001944;
  font-size: 1.15rem;
}

.details-panel {
  padding: 30px;
  border-radius: 28px;
}

.details-panel__header {
  margin-bottom: 20px;
}

.details-table {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 24px;
}

.details-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid rgba(17, 29, 35, 0.08);
}

.details-row span {
  color: #60707f;
}

.details-row strong {
  color: #001944;
  text-align: right;
}

.booking-card {
  position: sticky;
  top: 92px;
  display: grid;
  gap: 18px;
  padding: 28px;
  border-radius: 30px;
}

.booking-card__top strong {
  color: #001944;
  font-size: 2.4rem;
  line-height: 1;
  letter-spacing: -0.05em;
}

.booking-card__top span {
  color: #60707f;
  font-size: 1rem;
}

.booking-card__fields {
  display: grid;
  gap: 14px;
}

.booking-field {
  display: grid;
  gap: 10px;
}

.booking-field span,
.availability__title,
.booking-summary__title {
  color: #001944;
  font-size: 0.95rem;
  font-weight: 800;
}

.booking-field input {
  min-height: 52px;
  padding: 0 16px;
  border: 1px solid #d6dfe7;
  border-radius: 16px;
  background: #f8fbfd;
  color: #18222d;
  font: inherit;
}

.booking-field input:focus {
  outline: none;
  border-color: #3b63b7;
  box-shadow: 0 0 0 4px rgba(59, 99, 183, 0.12);
}

.availability {
  padding: 18px;
  border-radius: 20px;
  background: #f5f9fc;
}

.availability__list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
}

.availability__list span {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 12px;
  border-radius: 999px;
  background: rgba(0, 25, 68, 0.08);
  color: #36506b;
  font-size: 13px;
  font-weight: 600;
}

.booking-summary {
  padding: 18px;
  border-radius: 20px;
}

.booking-summary--muted {
  background: #f2f6fb;
}

.booking-summary--success {
  background: rgba(105, 255, 135, 0.16);
}

.booking-summary--error {
  background: rgba(255, 218, 214, 0.8);
}

.booking-summary__text {
  margin-top: 8px;
  color: #50606f;
  line-height: 1.6;
}

.booking-summary__price {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 16px;
  color: #001944;
  font-weight: 700;
}

.booking-button,
.not-found__link,
.lightbox__close,
.lightbox__nav {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  cursor: pointer;
  transition: 0.2s ease;
}

.booking-button {
  min-height: 56px;
  border-radius: 18px;
  background: #001944;
  color: #fff;
  font-size: 1rem;
  font-weight: 800;
}

.booking-button:disabled {
  background: #a4b1c0;
  cursor: not-allowed;
}

.booking-note {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  line-height: 1.6;
}

.booking-result {
  padding: 18px;
  border-radius: 20px;
  background: #edf4ff;
}

.booking-result p {
  margin-top: 10px;
  line-height: 1.7;
}

.not-found {
  max-width: 680px;
  margin: 48px auto 0;
  padding: 40px;
  border-radius: 30px;
  text-align: center;
}

.not-found__link {
  min-height: 48px;
  margin-top: 24px;
  padding: 0 20px;
  border-radius: 14px;
  background: #001944;
  color: #fff;
  text-decoration: none;
  font-weight: 700;
}

.lightbox {
  position: fixed;
  inset: 0;
  z-index: 140;
  display: grid;
  place-items: center;
  background: rgba(7, 14, 24, 0.82);
  backdrop-filter: blur(10px);
}

.lightbox__image {
  width: min(92vw, 1200px);
  max-height: 86vh;
  border-radius: 28px;
  object-fit: contain;
}

.lightbox__close,
.lightbox__nav {
  position: absolute;
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.16);
  color: #fff;
}

.lightbox__close {
  top: 32px;
  right: 32px;
}

.lightbox__nav--prev {
  left: 32px;
}

.lightbox__nav--next {
  right: 32px;
}

@media (max-width: 1100px) {
  .content-grid {
    grid-template-columns: 1fr;
  }

  .booking-card {
    position: static;
  }
}

@media (max-width: 840px) {
  .gallery__thumbs,
  .spec-grid,
  .details-table {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 680px) {
  .page-shell {
    width: min(100%, calc(100% - 24px));
    padding: 18px 0 32px;
  }

  .gallery__hero {
    min-height: 240px;
    border-radius: 24px;
  }

  .gallery__thumbs,
  .spec-grid,
  .details-table {
    grid-template-columns: 1fr;
  }

  .hero-copy,
  .details-panel,
  .booking-card,
  .not-found {
    padding: 22px;
    border-radius: 24px;
  }

  .details-panel__header,
  .details-row {
    display: grid;
    gap: 6px;
  }

  .details-row strong {
    text-align: left;
  }

  .lightbox__nav {
    bottom: 24px;
    top: auto;
  }

  .lightbox__nav--prev {
    left: calc(50% - 62px);
  }

  .lightbox__nav--next {
    right: calc(50% - 62px);
  }

  .lightbox__close {
    top: 18px;
    right: 18px;
  }
}
</style>
