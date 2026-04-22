<template>
  <section class="car-page">
    <AppHeader />

    <main class="page-shell">
      <section v-if="isLoading" class="state-card">
        <h1>Загружаем автомобиль</h1>
        <p>Получаем актуальную карточку из car-service.</p>
      </section>

      <section v-else-if="errorMessage" class="state-card">
        <h1>Не удалось загрузить карточку</h1>
        <p>{{ errorMessage }}</p>
        <button class="state-button" type="button" @click="loadCar">Повторить</button>
      </section>

      <section v-else-if="!car" class="state-card">
        <h1>Автомобиль не найден</h1>
        <p>Возможно, карточка была удалена или ссылка устарела.</p>
        <RouterLink class="not-found__link" :to="{ name: 'cars' }">Вернуться в каталог</RouterLink>
      </section>

      <section v-else class="car-details">
        <div class="breadcrumbs">
          <RouterLink :to="{ name: 'cars' }">Автомобили</RouterLink>
          <span>/</span>
          <span>{{ carTitle }}</span>
        </div>

        <div class="content-grid">
          <div class="content-main">
            <section class="gallery">
              <button
                v-if="selectedImageUrl"
                class="gallery__hero"
                type="button"
                @click="openLightbox(selectedImageIndex)"
              >
                <img :src="selectedImageUrl" :alt="carTitle" />
              </button>
              <div v-else class="gallery__placeholder">
                <span class="material-symbols-outlined">image</span>
                <p>Изображения пока отсутствуют</p>
              </div>

              <div v-if="visibleThumbnails.length" class="gallery__thumbs-wrap">
                <button
                  v-if="hasThumbPagination"
                  class="gallery__thumbs-nav gallery__thumbs-nav--prev"
                  type="button"
                  :disabled="!canScrollThumbsPrev"
                  aria-label="Показать предыдущие фото"
                  @click="scrollThumbsPrev"
                >
                  <span class="material-symbols-outlined">chevron_left</span>
                </button>
                <div class="gallery__thumbs">
                  <button
                    v-for="(image, index) in visibleThumbnails"
                    :key="image.id"
                    class="gallery__thumb"
                    :class="{ 'gallery__thumb--active': thumbnailWindowStart + index === selectedImageIndex }"
                    type="button"
                    @click="selectedImageIndex = thumbnailWindowStart + index"
                  >
                    <img :src="image.url" :alt="`${carTitle} фото ${thumbnailWindowStart + index + 1}`" />
                  </button>
                </div>
                <button
                  v-if="hasThumbPagination"
                  class="gallery__thumbs-nav gallery__thumbs-nav--next"
                  type="button"
                  :disabled="!canScrollThumbsNext"
                  aria-label="Показать следующие фото"
                  @click="scrollThumbsNext"
                >
                  <span class="material-symbols-outlined">chevron_right</span>
                </button>
              </div>
            </section>

            <header class="hero-copy">
              <div class="hero-copy__header">
                <div>
                  <h1>{{ carTitle }}</h1>
                  <p class="hero-copy__text">
                    {{ car.description || 'Описание для этого автомобиля пока не добавлено.' }}
                  </p>
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
                  <strong>{{ car.model }}</strong>
                </div>
                <div class="details-row">
                  <span>Кузов</span>
                  <strong>{{ car.body_type }}</strong>
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
                  <strong>{{ car.seats_count }}</strong>
                </div>
                <div class="details-row">
                  <span>Тип топлива</span>
                  <strong>{{ car.fuel_type }}</strong>
                </div>
              </div>
            </section>
          </div>

          <aside ref="bookingCardRef" class="booking-card">
            <div class="booking-card__top">
              <p>Стоимость аренды</p>
              <div>
                <strong>{{ formatPrice(car.price_per_day) }}</strong>
                <span>/ сутки</span>
              </div>
            </div>

            <div class="booking-summary" :class="bookingState.toneClass">
              <p class="booking-summary__title">{{ bookingState.title }}</p>
              <p class="booking-summary__text">{{ bookingState.description }}</p>
              <div v-if="bookingState.ready" class="booking-summary__price">
                <span>{{ bookingState.days }} суток</span>
                <strong>{{ formatPrice(bookingState.totalPrice) }}</strong>
              </div>
            </div>

            <template v-if="authStore.isAuthenticated">
              <div class="booking-card__fields">
                <div ref="startBookingFieldRef" class="booking-field booking-field--calendar">
                  <span>Дата начала</span>
                  <button
                    class="booking-field__trigger"
                    :class="{ 'booking-field__trigger--placeholder': !bookingForm.startDate }"
                    type="button"
                    :disabled="isAvailabilityLoading"
                    @click="toggleCalendar('start')"
                  >
                    {{ isAvailabilityLoading ? 'Загрузка...' : formatFieldValue(bookingForm.startDate, 'Выбрать') }}
                  </button>

                  <div v-if="openCalendar === 'start'" class="booking-calendar">
                    <div class="booking-calendar__header">
                      <button class="booking-calendar__nav" type="button" @click="shiftCalendarMonth(-1)">
                        <span class="material-symbols-outlined">chevron_left</span>
                      </button>
                      <strong>{{ calendarMonthLabel }}</strong>
                      <button class="booking-calendar__nav" type="button" @click="shiftCalendarMonth(1)">
                        <span class="material-symbols-outlined">chevron_right</span>
                      </button>
                    </div>

                    <div class="booking-calendar__weekdays">
                      <span v-for="day in weekDays" :key="day">{{ day }}</span>
                    </div>

                    <div class="booking-calendar__grid">
                      <button
                        v-for="cell in calendarCells"
                        :key="cell.key"
                        class="booking-calendar__day"
                        :class="{
                          'booking-calendar__day--muted': !cell.isCurrentMonth,
                          'booking-calendar__day--selected': cell.isSelected,
                        }"
                        type="button"
                        :disabled="cell.isDisabled"
                        @click="selectCalendarDate(cell.date)"
                      >
                        {{ cell.label }}
                      </button>
                    </div>
                  </div>
                </div>

                <div ref="endBookingFieldRef" class="booking-field booking-field--calendar">
                  <span>Дата завершения</span>
                  <button
                    class="booking-field__trigger"
                    :class="{ 'booking-field__trigger--placeholder': !bookingForm.endDate }"
                    type="button"
                    :disabled="isAvailabilityLoading || !bookingForm.startDate"
                    @click="toggleCalendar('end')"
                  >
                    {{ formatFieldValue(bookingForm.endDate, 'Выбрать') }}
                  </button>

                  <div v-if="openCalendar === 'end'" class="booking-calendar">
                    <div class="booking-calendar__header">
                      <button class="booking-calendar__nav" type="button" @click="shiftCalendarMonth(-1)">
                        <span class="material-symbols-outlined">chevron_left</span>
                      </button>
                      <strong>{{ calendarMonthLabel }}</strong>
                      <button class="booking-calendar__nav" type="button" @click="shiftCalendarMonth(1)">
                        <span class="material-symbols-outlined">chevron_right</span>
                      </button>
                    </div>

                    <div class="booking-calendar__weekdays">
                      <span v-for="day in weekDays" :key="day">{{ day }}</span>
                    </div>

                    <div class="booking-calendar__grid">
                      <button
                        v-for="cell in calendarCells"
                        :key="cell.key"
                        class="booking-calendar__day"
                        :class="{
                          'booking-calendar__day--muted': !cell.isCurrentMonth,
                          'booking-calendar__day--selected': cell.isSelected,
                        }"
                        type="button"
                        :disabled="cell.isDisabled"
                        @click="selectCalendarDate(cell.date)"
                      >
                        {{ cell.label }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <button
                class="booking-button"
                type="button"
                :disabled="!bookingState.ready || isBookingSubmitting"
                @click="submitBooking"
              >
                {{ isBookingSubmitting ? 'Оформляем...' : 'Забронировать' }}
              </button>
            </template>

            <template v-else>
              <div v-if="false" class="availability-panel">
                <div class="availability-panel__header">
                  <p class="availability-panel__title">Занятые даты</p>
                  <span v-if="busyIntervals.length" class="availability-panel__badge">
                    {{ busyIntervals.length }}
                  </span>
                </div>

                <p v-if="isAvailabilityLoading" class="availability-panel__note">
                  Обновляем занятые периоды...
                </p>

                <template v-else-if="availability">
                  <p class="availability-panel__note">
                    Показываем занятые интервалы до {{ availabilityWindowEndLabel }}. После входа можно
                    сразу выбрать свободный диапазон.
                  </p>

                  <ul v-if="visibleBusyIntervals.length" class="availability-list">
                    <li
                      v-for="interval in visibleBusyIntervals"
                      :key="`guest-${interval.start_date}-${interval.end_date}`"
                      class="availability-list__item"
                    >
                      <span class="material-symbols-outlined">event_busy</span>
                      <div>
                        <strong>{{ formatBusyInterval(interval) }}</strong>
                        <p>{{ getBusyIntervalDays(interval) }} суток</p>
                      </div>
                    </li>
                  </ul>

                  <p v-else class="availability-panel__empty">
                    В ближайшем окне занятых дат пока нет.
                  </p>
                </template>

                <p
                  v-else-if="availabilityErrorMessage"
                  class="availability-panel__note availability-panel__note--warning"
                >
                  {{ availabilityErrorMessage }}
                </p>
              </div>

              <div class="booking-guest">
                <p class="booking-guest__text">
                  Войдите, чтобы выбрать даты и оформить бронь.
                </p>
                <RouterLink class="booking-login-link" :to="loginLink">
                  Войти, чтобы забронировать
                </RouterLink>
              </div>
            </template>

            <div v-if="bookingErrorMessage" class="booking-result booking-result--error">
              <h3>Бронирование не создано</h3>
              <p>{{ bookingErrorMessage }}</p>
            </div>

            <div v-if="bookingResult" class="booking-result booking-result--success">
              <h3>Бронирование оформлено</h3>
              <p>{{ bookingResult }}</p>
              <RouterLink class="booking-result__link" :to="{ name: 'profile' }">
                Перейти в историю бронирований
              </RouterLink>
            </div>
          </aside>
        </div>
      </section>
    </main>

    <AppFooter />

    <div v-if="car && lightboxOpen && selectedImageUrl" class="lightbox" @click.self="closeLightbox">
      <button class="lightbox__close" type="button" @click="closeLightbox">
        <span class="material-symbols-outlined">close</span>
      </button>
      <button class="lightbox__nav lightbox__nav--prev" type="button" @click="prevImage">
        <span class="material-symbols-outlined">chevron_left</span>
      </button>
      <img class="lightbox__image" :src="selectedImageUrl" :alt="carTitle" />
      <button class="lightbox__nav lightbox__nav--next" type="button" @click="nextImage">
        <span class="material-symbols-outlined">chevron_right</span>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
import {
  createBooking,
  getBookingErrorMessage,
  getCarAvailability,
} from '@/services/bookings'
import type { BookingBusyInterval, CarAvailabilityResponse } from '@/services/bookings'
import type { CarDetailsResponse } from '@/services/cars'
import { getCarDetails } from '@/services/cars'
import { ApiError, useAuthStore } from '@/stores/auth'
import { useFavoritesStore } from '@/stores/favorites'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const favoritesStore = useFavoritesStore()

const car = ref<CarDetailsResponse | null>(null)
const availability = ref<CarAvailabilityResponse | null>(null)
const isLoading = ref(false)
const isAvailabilityLoading = ref(false)
const errorMessage = ref('')
const availabilityErrorMessage = ref('')
const selectedImageIndex = ref(0)
const thumbnailWindowStart = ref(0)
const lightboxOpen = ref(false)
const bookingResult = ref('')
const bookingErrorMessage = ref('')
const isBookingSubmitting = ref(false)
const preserveBookingFeedback = ref(false)
const bookingCardRef = ref<HTMLElement | null>(null)
const startBookingFieldRef = ref<HTMLElement | null>(null)
const endBookingFieldRef = ref<HTMLElement | null>(null)
const visibleThumbCount = 4
const visibleBusyIntervalsCount = 6
const today = getToday()
const availabilityRangeEnd = addDays(today, 365)
const openCalendar = ref<'start' | 'end' | null>(null)
const calendarMonth = ref(startOfMonth(today))
const weekDays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const bookingForm = reactive({
  startDate: '',
  endDate: '',
})

let currentController: AbortController | null = null

const carId = computed<number | null>(() => {
  const rawId = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  const parsedId = Number(rawId)

  return Number.isInteger(parsedId) && parsedId > 0 ? parsedId : null
})

const carTitle = computed(() => (car.value ? `${car.value.brand} ${car.value.model}` : 'Автомобиль'))
const selectedImageUrl = computed(() => car.value?.images[selectedImageIndex.value]?.url || '')
const hasThumbPagination = computed(() => (car.value?.images.length ?? 0) > visibleThumbCount)
const visibleThumbnails = computed(() =>
  car.value?.images.slice(thumbnailWindowStart.value, thumbnailWindowStart.value + visibleThumbCount) ?? [],
)
const canScrollThumbsPrev = computed(() => thumbnailWindowStart.value > 0)
const canScrollThumbsNext = computed(() => {
  if (!car.value) {
    return false
  }

  return thumbnailWindowStart.value + visibleThumbCount < car.value.images.length
})
const isFavorite = computed(() => (car.value ? favoritesStore.isFavorite(car.value.id) : false))
const favoriteButtonLabel = computed(() =>
  isFavorite.value ? 'Удалить из избранного' : 'Добавить в избранное',
)
const loginLink = computed(() => ({
  name: 'login' as const,
  query: {
    redirect: route.fullPath,
  },
}))
const busyIntervals = computed(() => availability.value?.busy_intervals ?? [])
const calendarMonthLabel = computed(() => {
  return new Intl.DateTimeFormat('ru-RU', {
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(parseDateOnly(calendarMonth.value))
})
const calendarCells = computed(() => {
  if (!openCalendar.value) {
    return []
  }

  const monthStartDate = parseDateOnly(calendarMonth.value)
  const firstDay = monthStartDate.getUTCDay()
  const offset = firstDay === 0 ? 6 : firstDay - 1
  const firstCellDate = addDays(calendarMonth.value, -offset)

  return Array.from({ length: 42 }, (_, index) => {
    const date = addDays(firstCellDate, index)

    return {
      key: `${openCalendar.value}-${date}`,
      date,
      label: parseDateOnly(date).getUTCDate(),
      isCurrentMonth: date.slice(0, 7) === calendarMonth.value.slice(0, 7),
      isDisabled: isCalendarDateDisabled(date),
      isSelected:
        (openCalendar.value === 'start' && bookingForm.startDate === date) ||
        (openCalendar.value === 'end' && bookingForm.endDate === date),
    }
  })
})
const visibleBusyIntervals = computed(() => busyIntervals.value.slice(0, visibleBusyIntervalsCount))
const availabilityWindowEndLabel = computed(() => {
  if (!availability.value) {
    return ''
  }

  return formatDate(addDays(availability.value.to, -1))
})
const startDateBusyInterval = computed(() =>
  findIntervalContainingDate(bookingForm.startDate, busyIntervals.value),
)
const selectedBusyInterval = computed(() =>
  findOverlappingInterval(bookingForm.startDate, bookingForm.endDate, busyIntervals.value),
)

const primaryFeatures = computed(() => {
  if (!car.value) {
    return []
  }

  return [
    { label: 'Год выпуска', value: String(car.value.year), icon: 'calendar_today' },
    { label: 'Топливо', value: car.value.fuel_type, icon: 'local_gas_station' },
    { label: 'Трансмиссия', value: car.value.transmission, icon: 'settings' },
    { label: 'Мест', value: String(car.value.seats_count), icon: 'group' },
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

  if (!authStore.isAuthenticated) {
    return {
      ready: false,
      title: 'Бронирование доступно после входа',
      description: 'Авторизуйтесь, чтобы выбрать даты, оформить бронь и видеть её историю в профиле.',
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--info',
    }
  }

  if (startDateBusyInterval.value) {
    return {
      ready: false,
      title: 'Дата начала недоступна',
      description: `Этот день входит в занятый интервал. Попробуйте начать аренду с ${formatDate(getBusyIntervalEndDate(startDateBusyInterval.value))}.`,
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--error',
    }
  }

  if (!bookingForm.startDate || !bookingForm.endDate) {
    return {
      ready: false,
      title: 'Выберите даты аренды',
      description: 'Укажите дату начала и дату завершения, чтобы увидеть предварительный расчет.',
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

  if (selectedBusyInterval.value) {
    return {
      ready: false,
      title: 'Диапазон занят',
      description: `Выбранный период пересекается с интервалом ${formatBusyInterval(selectedBusyInterval.value)}.`,
      days: 0,
      totalPrice: 0,
      toneClass: 'booking-summary--error',
    }
  }

  const days = getDaysBetween(bookingForm.startDate, bookingForm.endDate)

  return {
    ready: true,
    title: 'Предварительный расчет готов',
    description: `Стоимость аренды с ${formatDate(bookingForm.startDate)} по ${formatDate(bookingForm.endDate)}.`,
    days,
    totalPrice: days * car.value.price_per_day,
    toneClass: 'booking-summary--success',
  }
})

async function loadAvailability(carIdValue: number, signal?: AbortSignal) {
  isAvailabilityLoading.value = true
  availabilityErrorMessage.value = ''

  try {
    availability.value = await getCarAvailability(
      carIdValue,
      {
        from: today,
        to: availabilityRangeEnd,
      },
      signal,
    )
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      return
    }

    availability.value = null
    availabilityErrorMessage.value =
      error instanceof Error
        ? error.message
        : 'Не удалось загрузить информацию о занятых датах. Попробуйте снова, или войдите в аккаунт, чтобы видеть актуальную доступность при выборе дат.'
  } finally {
    isAvailabilityLoading.value = false
  }
}

async function loadCar() {
  currentController?.abort()
  currentController = new AbortController()
  isLoading.value = true
  errorMessage.value = ''
  car.value = null
  availability.value = null
  availabilityErrorMessage.value = ''

  if (carId.value === null) {
    isLoading.value = false
    return
  }

  try {
    car.value = await getCarDetails(carId.value, currentController.signal)
    selectedImageIndex.value = 0
    thumbnailWindowStart.value = 0
    await loadAvailability(carId.value, currentController.signal)

    if (authStore.isAuthenticated) {
      await favoritesStore.ensureLoaded()
    }
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      return
    }

    if (error instanceof ApiError && error.status === 404) {
      car.value = null
      return
    }

    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить карточку автомобиля.'
  } finally {
    isLoading.value = false
  }
}

async function handleFavoriteClick() {
  if (!car.value) {
    return
  }

  if (!authStore.isAuthenticated) {
    await router.push(loginLink.value)
    return
  }

  try {
    await favoritesStore.toggleFavorite(car.value.id)
  } catch {
    // The store keeps the latest backend error for views that need it.
  }
}

async function submitBooking() {
  if (!car.value || !bookingState.value.ready || isBookingSubmitting.value) {
    return
  }

  bookingErrorMessage.value = ''
  bookingResult.value = ''
  const startDate = bookingForm.startDate
  const endDate = bookingForm.endDate

  try {
    isBookingSubmitting.value = true

    await createBooking(authStore.authorizedRequest, {
      car_id: car.value.id,
      start_date: toApiDateTime(startDate),
      end_date: toApiDateTime(endDate),
    })

    bookingResult.value =
      `${carTitle.value} забронирован с ${formatDate(startDate)} по ${formatDate(endDate)}. ` +
      'Теперь бронь можно отслеживать и отменять в профиле.'

    preserveBookingFeedback.value = true  
    bookingForm.startDate = ''
    bookingForm.endDate = ''
    openCalendar.value = null
    await loadAvailability(car.value.id)
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      await router.push(loginLink.value)
      return
    }

    if (error instanceof ApiError && error.status === 409) {
      await loadAvailability(car.value.id)
    }

    bookingErrorMessage.value = getBookingErrorMessage(error, 'Не удалось создать бронирование.')
  } finally {
    isBookingSubmitting.value = false
  }
}

function toggleCalendar(kind: 'start' | 'end') {
  if (isAvailabilityLoading.value) {
    return
  }

  if (kind === 'end' && !bookingForm.startDate) {
    return
  }

  if (openCalendar.value === kind) {
    openCalendar.value = null
    return
  }

  openCalendar.value = kind
  calendarMonth.value = startOfMonth(
    kind === 'start' ? bookingForm.startDate || today : bookingForm.endDate || bookingForm.startDate || today,
  )
  void scrollBookingFieldIntoView(kind)
}

function shiftCalendarMonth(step: number) {
  const monthDate = parseDateOnly(calendarMonth.value)
  monthDate.setUTCMonth(monthDate.getUTCMonth() + step)
  calendarMonth.value = monthDate.toISOString().slice(0, 10)
}

function selectCalendarDate(date: string) {
  if (isCalendarDateDisabled(date)) {
    return
  }

  if (openCalendar.value === 'start') {
    bookingForm.startDate = date

    if (
      bookingForm.endDate &&
      (bookingForm.endDate <= date || findOverlappingInterval(date, bookingForm.endDate, busyIntervals.value))
    ) {
      bookingForm.endDate = ''
    }

    openCalendar.value = 'end'
    calendarMonth.value = startOfMonth(bookingForm.startDate || today)
    void scrollBookingFieldIntoView('end')
    return
  }

  bookingForm.endDate = date
  openCalendar.value = null
}

function isCalendarDateDisabled(date: string) {
  if (date < today) {
    return true
  }

  if (openCalendar.value === 'start') {
    return Boolean(findIntervalContainingDate(date, busyIntervals.value))
  }

  if (openCalendar.value === 'end') {
    if (!bookingForm.startDate) {
      return true
    }

    return date <= bookingForm.startDate || Boolean(findOverlappingInterval(bookingForm.startDate, date, busyIntervals.value))
  }

  return true
}

function formatFieldValue(dateString: string, fallback: string) {
  return dateString ? formatDate(dateString) : fallback
}

function openLightbox(index: number) {
  if (!selectedImageUrl.value) {
    return
  }

  selectedImageIndex.value = index
  lightboxOpen.value = true
}

function closeLightbox() {
  lightboxOpen.value = false
}

function prevImage() {
  if (!car.value?.images.length) {
    return
  }

  selectedImageIndex.value =
    (selectedImageIndex.value - 1 + car.value.images.length) % car.value.images.length
}

function nextImage() {
  if (!car.value?.images.length) {
    return
  }

  selectedImageIndex.value = (selectedImageIndex.value + 1) % car.value.images.length
}

function scrollThumbsPrev() {
  thumbnailWindowStart.value = Math.max(0, thumbnailWindowStart.value - 1)
}

function scrollThumbsNext() {
  if (!car.value) {
    return
  }

  thumbnailWindowStart.value = Math.min(
    car.value.images.length - visibleThumbCount,
    thumbnailWindowStart.value + 1,
  )
}

watch(
  () => route.params.id,
  () => {
    closeLightbox()
    openCalendar.value = null
    bookingResult.value = ''
    bookingErrorMessage.value = ''
    bookingForm.startDate = ''
    bookingForm.endDate = ''
    void loadCar()
  },
)

watch([() => bookingForm.startDate, () => bookingForm.endDate], () => {
    if (preserveBookingFeedback.value) {
    preserveBookingFeedback.value = false
    return
    }
  
  bookingResult.value = ''
  bookingErrorMessage.value = ''
})

watch(
  () => bookingForm.startDate,
  (startDate) => {
    if (
      bookingForm.endDate &&
      (!startDate || bookingForm.endDate <= startDate || findOverlappingInterval(startDate, bookingForm.endDate, busyIntervals.value))
    ) {
      bookingForm.endDate = ''
    }
  },
)

watch(selectedImageIndex, (index) => {
  if (index < thumbnailWindowStart.value) {
    thumbnailWindowStart.value = index
    return
  }

  const visibleWindowEnd = thumbnailWindowStart.value + visibleThumbCount - 1
  if (index > visibleWindowEnd) {
    thumbnailWindowStart.value = index - visibleThumbCount + 1
  }
})

onMounted(() => {
  calendarMonth.value = startOfMonth(today)
  void loadCar()
})

onBeforeUnmount(() => {
  currentController?.abort()
})

function formatPrice(price: number) {
  return `${new Intl.NumberFormat('ru-RU').format(price)} ₽`
}

function formatDate(dateString: string) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(parseDateOnly(dateString))
}

function formatBusyInterval(interval: BookingBusyInterval) {
  const busyUntil = addDays(interval.end_date, -1)

  if (interval.start_date === busyUntil) {
    return formatDate(interval.start_date)
  }

  return `${formatDate(interval.start_date)} - ${formatDate(busyUntil)}`
}

function getBusyIntervalDays(interval: BookingBusyInterval) {
  return getDaysBetween(interval.start_date, interval.end_date)
}

function getBusyIntervalEndDate(interval: BookingBusyInterval | null | undefined) {
  return interval?.end_date || today
}

function findIntervalContainingDate(dateString: string, intervals: BookingBusyInterval[]) {
  if (!dateString) {
    return null
  }

  return intervals.find((interval) => interval.start_date <= dateString && interval.end_date > dateString) ?? null
}

function findOverlappingInterval(
  startDate: string,
  endDate: string,
  intervals: BookingBusyInterval[],
) {
  if (!startDate || !endDate || endDate <= startDate) {
    return null
  }

  return intervals.find((interval) => interval.start_date < endDate && interval.end_date > startDate) ?? null
}

function getDaysBetween(start: string, end: string) {
  const startDate = parseDateOnly(start)
  const endDate = parseDateOnly(end)
  const millisecondsPerDay = 1000 * 60 * 60 * 24

  return Math.ceil((endDate.getTime() - startDate.getTime()) / millisecondsPerDay)
}

function parseDateOnly(dateString: string) {
  return new Date(`${dateString}T00:00:00.000Z`)
}

function addDays(dateString: string, days: number) {
  const date = parseDateOnly(dateString)
  date.setUTCDate(date.getUTCDate() + days)

  return date.toISOString().slice(0, 10)
}

function startOfMonth(dateString: string) {
  return `${dateString.slice(0, 7)}-01`
}

function getToday() {
  const now = new Date()
  const offset = now.getTimezoneOffset()
  const normalized = new Date(now.getTime() - offset * 60 * 1000)

  return normalized.toISOString().slice(0, 10)
}

function toApiDateTime(dateString: string) {
  return `${dateString}T00:00:00.000Z`
}

async function scrollBookingFieldIntoView(kind: 'start' | 'end') {
  await nextTick()

  const field = kind === 'start' ? startBookingFieldRef.value : endBookingFieldRef.value
  if (!field) {
    return
  }

  field.scrollIntoView({
    block: 'nearest',
    behavior: 'smooth',
  })
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

.state-card,
.hero-copy,
.details-panel,
.booking-card {
  border: 1px solid rgba(17, 29, 35, 0.08);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 20px 50px rgba(14, 40, 64, 0.06);
}

.state-card {
  max-width: 680px;
  margin: 48px auto 0;
  padding: 40px;
  border-radius: 30px;
  text-align: center;
}

.state-card h1 {
  margin: 0;
}

.state-card p {
  margin: 16px 0 0;
  color: #60707f;
}

.state-button,
.not-found__link,
.lightbox__close,
.lightbox__nav,
.favorite-button,
.booking-button,
.booking-login-link,
.booking-result__link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  cursor: pointer;
  transition: 0.2s ease;
}

.state-button,
.not-found__link,
.booking-login-link,
.booking-result__link {
  min-height: 48px;
  margin-top: 24px;
  padding: 0 20px;
  border-radius: 14px;
  background: #001944;
  color: #fff;
  text-decoration: none;
  font-weight: 700;
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

.gallery__hero,
.gallery__placeholder {
  position: relative;
  height: 420px;
  min-height: 280px;
  overflow: hidden;
  border-radius: 28px;
}

.gallery__hero {
  border: 0;
  padding: 0;
  background: #dfe8ef;
  cursor: pointer;
}

.gallery__placeholder {
  display: grid;
  place-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #eef4f8 0%, #d8e2eb 100%);
  color: #526170;
  text-align: center;
  padding: 24px;
}

.gallery__placeholder .material-symbols-outlined {
  font-size: 40px;
}

.gallery__hero img,
.gallery__thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gallery__thumbs-wrap {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 12px;
  align-items: center;
}

.gallery__thumbs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.gallery__thumbs-nav {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border: 1px solid rgba(17, 29, 35, 0.08);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.92);
  color: #163f77;
  box-shadow: 0 12px 30px rgba(14, 40, 64, 0.08);
}

.gallery__thumbs-nav:disabled {
  opacity: 0.42;
  cursor: not-allowed;
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

.hero-copy {
  padding: 34px;
  border-radius: 30px;
}

.hero-copy__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
}

.hero-copy h1,
.details-panel h2,
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

.favorite-button {
  width: 52px;
  height: 52px;
  flex-shrink: 0;
  border-radius: 18px;
  background: #edf4f9;
  color: #163f77;
}

.favorite-button:disabled {
  opacity: 0.72;
  cursor: not-allowed;
}

.favorite-button .material-symbols-outlined {
  font-size: 28px;
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

.hero-copy__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 22px;
}

.info-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 38px;
  padding: 0 16px;
  border-radius: 999px;
  background: #edf4f9;
  color: #163f77;
  font-size: 13px;
  font-weight: 800;
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
.booking-card__top p,
.booking-note,
.booking-result p,
.booking-guest__text {
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
  max-height: calc(100vh - 116px);
  padding: 28px;
  border-radius: 30px;
  overflow-y: auto;
  overscroll-behavior: contain;
  scrollbar-gutter: stable;
}

.booking-card::-webkit-scrollbar {
  width: 10px;
}

.booking-card::-webkit-scrollbar-track {
  background: transparent;
}

.booking-card::-webkit-scrollbar-thumb {
  border: 2px solid transparent;
  border-radius: 999px;
  background: rgba(22, 63, 119, 0.22);
  background-clip: padding-box;
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
.booking-summary__title {
  color: #001944;
  font-size: 0.95rem;
  font-weight: 800;
}

.booking-field__input {
  min-height: 52px;
  padding: 0 16px;
  border: 1px solid #d6dfe7;
  border-radius: 16px;
  background: #f8fbfd;
  color: #18222d;
  font: inherit;
}

.booking-field__input:focus {
  outline: none;
  border-color: #3b63b7;
  box-shadow: 0 0 0 4px rgba(59, 99, 183, 0.12);
}

.booking-field__input--invalid {
  border-color: rgba(184, 57, 38, 0.55);
  background: #fff6f4;
}

.booking-field--calendar {
  position: relative;
}

.booking-field__trigger {
  min-height: 52px;
  padding: 0 16px;
  border: 1px solid #d6dfe7;
  border-radius: 16px;
  background: #f8fbfd;
  color: #18222d;
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.booking-field__trigger--placeholder {
  color: #60707f;
}

.booking-field__trigger:disabled {
  cursor: not-allowed;
  opacity: 0.72;
}

.booking-calendar {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  z-index: 5;
  width: min(320px, 100%);
  padding: 16px;
  border: 1px solid rgba(17, 29, 35, 0.08);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 20px 50px rgba(14, 40, 64, 0.12);
}

.booking-calendar__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  color: #001944;
}

.booking-calendar__nav {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid rgba(17, 29, 35, 0.08);
  border-radius: 12px;
  background: #f7fbff;
  color: #163f77;
  cursor: pointer;
}

.booking-calendar__weekdays,
.booking-calendar__grid {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 6px;
}

.booking-calendar__weekdays {
  margin-bottom: 8px;
  color: #60707f;
  font-size: 12px;
  text-align: center;
}

.booking-calendar__day {
  min-height: 38px;
  border: 0;
  border-radius: 12px;
  background: #eef4f9;
  color: #18222d;
  cursor: pointer;
  font: inherit;
}

.booking-calendar__day--muted {
  opacity: 0.42;
}

.booking-calendar__day--selected {
  background: #001944;
  color: #fff;
}

.booking-calendar__day:disabled {
  background: #f3f5f8;
  color: #a4b1c0;
  cursor: not-allowed;
}

.booking-summary {
  padding: 18px;
  border-radius: 20px;
}

.booking-summary--muted {
  background: #f2f6fb;
}

.booking-summary--info {
  background: rgba(22, 63, 119, 0.08);
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

.availability-panel {
  display: grid;
  gap: 12px;
  padding: 16px;
  border: 1px solid rgba(17, 29, 35, 0.08);
  border-radius: 20px;
  background: #f7fbff;
}

.availability-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.availability-panel__title {
  margin: 0;
  color: #001944;
  font-size: 0.95rem;
  font-weight: 800;
}

.availability-panel__badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  min-height: 32px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(22, 63, 119, 0.12);
  color: #163f77;
  font-size: 0.85rem;
  font-weight: 800;
}

.availability-panel__note,
.availability-panel__empty,
.availability-list__item p,
.booking-inline-message {
  margin: 0;
  color: #60707f;
  line-height: 1.6;
}

.availability-panel__note--warning,
.booking-inline-message--error {
  color: #9d3928;
}

.availability-list {
  display: grid;
  gap: 10px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.availability-list__item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 12px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.82);
}

.availability-list__item .material-symbols-outlined {
  color: #163f77;
}

.availability-list__item strong {
  color: #001944;
}

.booking-note {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  line-height: 1.6;
}

.booking-guest {
  display: grid;
  gap: 14px;
  padding: 4px 0;
}

.booking-login-link {
  margin-top: 0;
}

.booking-result {
  display: grid;
  gap: 10px;
  padding: 18px;
  border-radius: 20px;
}

.booking-result--success {
  background: #edf4ff;
}

.booking-result--error {
  background: #fff0ee;
}

.booking-result p {
  line-height: 1.7;
}

.booking-result__link {
  width: 100%;
  margin-top: 4px;
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
    max-height: none;
    overflow: visible;
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

  .gallery__hero,
  .gallery__placeholder {
    min-height: 240px;
    border-radius: 24px;
  }

  .gallery__thumbs,
  .spec-grid,
  .details-table {
    grid-template-columns: 1fr;
  }

  .gallery__thumbs-wrap {
    grid-template-columns: 1fr;
  }

  .gallery__thumbs-nav {
    width: 100%;
  }

  .hero-copy,
  .details-panel,
  .booking-card,
  .state-card {
    padding: 22px;
    border-radius: 24px;
  }

  .hero-copy__header,
  .details-row {
    display: grid;
    gap: 12px;
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
