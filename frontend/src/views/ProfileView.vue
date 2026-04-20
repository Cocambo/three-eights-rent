<template>
  <div class="profile-page">
    <AppHeader />

    <main class="profile-page__main">
      <section class="profile-card">
        <div class="profile-card__header">
          <div class="avatar-block">
            <div class="avatar-block__image-wrap">
              <img class="avatar-block__image" :src="avatar" :alt="authStore.displayName" />
            </div>
          </div>

          <div class="profile-card__title-wrap">
            <h1 class="profile-card__title">{{ authStore.displayName }}</h1>
            <p class="profile-card__subtitle">{{ authStore.user?.email }}</p>
          </div>
        </div>

        <div class="profile-form">
          <div class="form-field">
            <label class="form-field__label" for="firstName">Имя</label>
            <input
              id="firstName"
              v-model="form.firstName"
              class="form-field__input"
              type="text"
              placeholder="Введите имя"
            />
          </div>

          <div class="form-field">
            <label class="form-field__label" for="lastName">Фамилия</label>
            <input
              id="lastName"
              v-model="form.lastName"
              class="form-field__input"
              type="text"
              placeholder="Введите фамилию"
            />
          </div>

          <div class="form-field">
            <label class="form-field__label" for="birthDate">Дата рождения</label>
            <input id="birthDate" v-model="form.birthDate" class="form-field__input" type="date" />
          </div>

          <div class="form-field">
            <label class="form-field__label" for="email">Электронная почта</label>
            <input
              id="email"
              :value="authStore.user?.email || ''"
              class="form-field__input"
              type="email"
              readonly
            />
          </div>
        </div>
      </section>

      <section class="license-card">
        <div class="section-title">
          <span class="material-symbols-outlined section-title__icon">badge</span>
          <h2 class="section-title__text">Водительское удостоверение</h2>
        </div>

        <div class="license-form">
          <div class="form-field">
            <label class="form-field__label" for="licenseNumber">Серия и номер</label>
            <input
              id="licenseNumber"
              v-model="form.licenseNumber"
              class="form-field__input"
              type="text"
              placeholder="77 12 345678"
            />
          </div>

          <div class="form-field">
            <label class="form-field__label" for="issueDate">Дата выдачи</label>
            <input id="issueDate" v-model="form.issueDate" class="form-field__input" type="date" />
          </div>

          <div class="form-field">
            <label class="form-field__label" for="expirationDate">Дата окончания действия</label>
            <input
              id="expirationDate"
              v-model="form.expirationDate"
              class="form-field__input"
              type="date"
            />
          </div>
        </div>

        <div class="categories">
          <span class="categories__label">Открытые категории</span>

          <div class="categories__list">
            <label v-for="category in categories" :key="category" class="category-chip">
              <input
                v-model="form.categories"
                class="category-chip__input"
                type="checkbox"
                :value="category"
              />
              <span class="category-chip__text">{{ category }}</span>
            </label>
          </div>
        </div>
      </section>

      <p v-if="feedbackMessage" :class="feedbackClass">
        {{ feedbackMessage }}
      </p>

      <div class="profile-actions">
        <button
          class="profile-actions__logout"
          type="button"
          :disabled="isLoggingOut"
          @click="handleLogout"
        >
          {{ isLoggingOut ? 'Выходим...' : 'Выйти' }}
        </button>

        <button
          class="profile-actions__save"
          type="button"
          :disabled="!hasUnsavedChanges || isSaving || authStore.isProfileLoading"
          @click="saveProfile"
        >
          {{ isSaving ? 'Сохраняем...' : 'Сохранить изменения' }}
        </button>
      </div>

      <section class="bookings-card">
        <div class="section-title section-title--space">
          <div class="section-title__copy">
            <span class="material-symbols-outlined section-title__icon">event_note</span>
            <div>
              <h2 class="section-title__text">История бронирований</h2>
              <p class="section-title__description">
                Здесь можно посмотреть активные и завершённые брони, а также отменить текущую.
              </p>
            </div>
          </div>

          <button
            class="bookings-refresh"
            type="button"
            :disabled="isBookingsLoading"
            @click="loadBookings"
          >
            {{ isBookingsLoading ? 'Обновляем...' : 'Обновить' }}
          </button>
        </div>

        <p v-if="bookingActionMessage" :class="bookingActionClass">
          {{ bookingActionMessage }}
        </p>

        <div v-if="isBookingsLoading" class="bookings-state">
          <h3>Загружаем историю</h3>
          <p>Получаем ваши бронирования из car-service.</p>
        </div>

        <div v-else-if="bookingsErrorMessage" class="bookings-state bookings-state--error">
          <h3>Не удалось получить бронирования</h3>
          <p>{{ bookingsErrorMessage }}</p>
        </div>

        <div v-else-if="!bookings.length" class="bookings-state">
          <h3>Пока нет бронирований</h3>
          <p>Когда вы забронируете автомобиль, история появится здесь.</p>
          <RouterLink class="bookings-state__link" :to="{ name: 'cars' }">Перейти в каталог</RouterLink>
        </div>

        <div v-else class="bookings-list">
          <article v-for="booking in bookings" :key="booking.id" class="booking-item">
            <RouterLink
              class="booking-item__media"
              :to="{ name: 'car-details', params: { id: booking.car_id } }"
            >
              <img
                v-if="booking.car.main_image_url"
                :src="booking.car.main_image_url"
                :alt="`${booking.car.brand} ${booking.car.model}`"
              />
              <div v-else class="booking-item__placeholder">
                <span class="material-symbols-outlined">directions_car</span>
              </div>
            </RouterLink>

            <div class="booking-item__body">
              <div class="booking-item__header">
                <div>
                  <RouterLink
                    class="booking-item__title"
                    :to="{ name: 'car-details', params: { id: booking.car_id } }"
                  >
                    {{ booking.car.brand }} {{ booking.car.model }}
                  </RouterLink>
                  <p class="booking-item__subtitle">
                    {{ booking.car.year }} • {{ booking.car.body_type }} • {{ booking.car.transmission }}
                  </p>
                </div>
                <span class="booking-status" :class="statusClass(booking.status)">
                  {{ statusLabel(booking.status) }}
                </span>
              </div>

              <div class="booking-meta">
                <div class="booking-meta__item">
                  <span>Период</span>
                  <strong>{{ formatDateRange(booking.start_date, booking.end_date) }}</strong>
                </div>
                <div class="booking-meta__item">
                  <span>Длительность</span>
                  <strong>{{ getDurationDays(booking.start_date, booking.end_date) }} суток</strong>
                </div>
                <div class="booking-meta__item">
                  <span>Стоимость в сутки</span>
                  <strong>{{ formatPrice(booking.car.price_per_day) }}</strong>
                </div>
                <div class="booking-meta__item">
                  <span>Создано</span>
                  <strong>{{ formatDateTime(booking.created_at) }}</strong>
                </div>
              </div>

              <div class="booking-item__footer">
                <p v-if="booking.cancelled_at" class="booking-item__hint">
                  Отменено {{ formatDateTime(booking.cancelled_at) }}
                </p>
                <p v-else-if="booking.status === 'completed'" class="booking-item__hint">
                  Поездка завершена
                </p>
                <p v-else class="booking-item__hint">
                  Активное бронирование можно отменить в один клик.
                </p>

                <button
                  v-if="booking.status === 'active'"
                  class="booking-item__cancel"
                  type="button"
                  :disabled="isBookingPending(booking.id)"
                  @click="handleCancelBooking(booking.id)"
                >
                  {{ isBookingPending(booking.id) ? 'Отменяем...' : 'Отменить бронь' }}
                </button>
              </div>
            </div>
          </article>
        </div>
      </section>
    </main>

    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'

import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import avatar from '@/assets/profile_avatar.png'
import {
  cancelBooking,
  getBookingErrorMessage,
  getBookings,
  type BookingHistoryItem,
} from '@/services/bookings'
import { ApiError, formatDateForInput, useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const isSaving = ref(false)
const isLoggingOut = ref(false)
const isBookingsLoading = ref(false)
const feedbackMessage = ref('')
const feedbackType = ref<'success' | 'error'>('success')
const bookingActionMessage = ref('')
const bookingActionType = ref<'success' | 'error'>('success')
const bookingsErrorMessage = ref('')
const bookings = ref<BookingHistoryItem[]>([])
const pendingBookingIds = ref<number[]>([])

const categories = ['A', 'B', 'C', 'D', 'E']

const form = reactive({
  firstName: '',
  lastName: '',
  birthDate: '',
  licenseNumber: '',
  issueDate: '',
  expirationDate: '',
  categories: [] as string[],
})

const feedbackClass = computed(() =>
  feedbackType.value === 'success'
    ? 'profile-feedback profile-feedback--success'
    : 'profile-feedback profile-feedback--error',
)

const bookingActionClass = computed(() =>
  bookingActionType.value === 'success'
    ? 'profile-feedback profile-feedback--success'
    : 'profile-feedback profile-feedback--error',
)

const hasUnsavedChanges = computed(() => {
  return JSON.stringify(getCurrentFormSnapshot()) !== JSON.stringify(getStoredFormSnapshot())
})

const fillFormFromStore = () => {
  form.firstName = authStore.profile?.first_name || ''
  form.lastName = authStore.profile?.last_name || ''
  form.birthDate = formatDateForInput(authStore.profile?.birth_date)
  form.licenseNumber = authStore.driverLicense?.license_number || ''
  form.issueDate = formatDateForInput(authStore.driverLicense?.issued_at)
  form.expirationDate = formatDateForInput(authStore.driverLicense?.expires_at)
  form.categories = authStore.driverLicense?.categories.map((category) => category.category_code) || []
}

const fetchProfile = async () => {
  try {
    feedbackMessage.value = ''
    await authStore.fetchProfile()
    fillFormFromStore()
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      await router.push({ name: 'login', query: { redirect: '/profile' } })
      return
    }

    feedbackType.value = 'error'
    feedbackMessage.value = 'Не удалось загрузить профиль.'
  }
}

const loadBookings = async () => {
  try {
    isBookingsLoading.value = true
    bookingsErrorMessage.value = ''
    const response = await getBookings(authStore.authorizedRequest)
    bookings.value = response.items
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      await router.push({ name: 'login', query: { redirect: '/profile' } })
      return
    }

    bookingsErrorMessage.value = getBookingErrorMessage(
      error,
      'Не удалось загрузить историю бронирований.',
    )
  } finally {
    isBookingsLoading.value = false
  }
}

const saveProfile = async () => {
  feedbackMessage.value = ''

  if (!form.firstName || !form.lastName || !form.birthDate) {
    feedbackType.value = 'error'
    feedbackMessage.value = 'Заполните имя, фамилию и дату рождения.'
    return
  }

  const hasAnyLicenseField = Boolean(
    form.licenseNumber || form.issueDate || form.expirationDate || form.categories.length,
  )

  if (
    hasAnyLicenseField
    && (!form.licenseNumber || !form.issueDate || !form.expirationDate || form.categories.length === 0)
  ) {
    feedbackType.value = 'error'
    feedbackMessage.value =
      'Чтобы сохранить водительское удостоверение, заполните номер, даты и выберите хотя бы одну категорию.'
    return
  }

  try {
    isSaving.value = true

    await authStore.updateProfile({
      firstName: form.firstName,
      lastName: form.lastName,
      birthDate: form.birthDate,
    })

    if (hasAnyLicenseField) {
      await authStore.saveDriverLicense({
        licenseNumber: form.licenseNumber,
        issuedAt: form.issueDate,
        expiresAt: form.expirationDate,
        categories: form.categories,
      })
    }

    fillFormFromStore()
    feedbackType.value = 'success'
    feedbackMessage.value = 'Профиль успешно сохранён.'
  } catch (error) {
    if (error instanceof ApiError) {
      feedbackType.value = 'error'
      feedbackMessage.value = error.message
      return
    }

    feedbackType.value = 'error'
    feedbackMessage.value = 'Не удалось сохранить профиль.'
  } finally {
    isSaving.value = false
  }
}

const handleCancelBooking = async (bookingId: number) => {
  if (pendingBookingIds.value.includes(bookingId)) {
    return
  }

  bookingActionMessage.value = ''
  pendingBookingIds.value = [...pendingBookingIds.value, bookingId]

  try {
    const response = await cancelBooking(authStore.authorizedRequest, bookingId)

    bookings.value = bookings.value.map((booking) =>
      booking.id === bookingId
        ? {
            ...booking,
            status: response.status,
            cancelled_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          }
        : booking,
    )

    bookingActionType.value = 'success'
    bookingActionMessage.value = 'Бронирование отменено.'
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      await router.push({ name: 'login', query: { redirect: '/profile' } })
      return
    }

    bookingActionType.value = 'error'
    bookingActionMessage.value = getBookingErrorMessage(error, 'Не удалось отменить бронирование.')
  } finally {
    pendingBookingIds.value = pendingBookingIds.value.filter((id) => id !== bookingId)
  }
}

const handleLogout = async () => {
  try {
    isLoggingOut.value = true
    await authStore.logout()
    await router.push({ name: 'login' })
  } catch {
    feedbackType.value = 'error'
    feedbackMessage.value = 'Не удалось завершить сессию.'
  } finally {
    isLoggingOut.value = false
  }
}

function isBookingPending(bookingId: number) {
  return pendingBookingIds.value.includes(bookingId)
}

function statusLabel(status: string) {
  switch (status) {
    case 'active':
      return 'Активно'
    case 'cancelled':
      return 'Отменено'
    case 'completed':
      return 'Завершено'
    default:
      return status
  }
}

function statusClass(status: string) {
  switch (status) {
    case 'active':
      return 'booking-status--active'
    case 'cancelled':
      return 'booking-status--cancelled'
    case 'completed':
      return 'booking-status--completed'
    default:
      return ''
  }
}

function formatBookingDate(value: string) {
  const dateKey = toDateKey(value)

  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(`${dateKey}T00:00:00Z`))
}

function formatDateTime(value: string) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(value))
}

function formatDateRange(startDate: string, endDate: string) {
  return `${formatBookingDate(startDate)} - ${formatBookingDate(endDate)}`
}

function getDurationDays(startDate: string, endDate: string) {
  const start = new Date(`${toDateKey(startDate)}T00:00:00Z`)
  const end = new Date(`${toDateKey(endDate)}T00:00:00Z`)
  const millisecondsPerDay = 1000 * 60 * 60 * 24

  return Math.max(1, Math.ceil((end.getTime() - start.getTime()) / millisecondsPerDay))
}

function formatPrice(price: number) {
  return `${new Intl.NumberFormat('ru-RU').format(price)} ₽`
}

function toDateKey(value: string) {
  return value.slice(0, 10)
}

function getStoredFormSnapshot() {
  return {
    firstName: authStore.profile?.first_name || '',
    lastName: authStore.profile?.last_name || '',
    birthDate: formatDateForInput(authStore.profile?.birth_date),
    licenseNumber: authStore.driverLicense?.license_number || '',
    issueDate: formatDateForInput(authStore.driverLicense?.issued_at),
    expirationDate: formatDateForInput(authStore.driverLicense?.expires_at),
    categories: [...(authStore.driverLicense?.categories.map((category) => category.category_code) || [])].sort(),
  }
}

function getCurrentFormSnapshot() {
  return {
    firstName: form.firstName.trim(),
    lastName: form.lastName.trim(),
    birthDate: form.birthDate,
    licenseNumber: form.licenseNumber.trim(),
    issueDate: form.issueDate,
    expirationDate: form.expirationDate,
    categories: [...form.categories].sort(),
  }
}

onMounted(() => {
  fillFormFromStore()
  void fetchProfile()
  void loadBookings()
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@700;800&family=Inter:wght@400;500;600&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap');

:global(*) {
  box-sizing: border-box;
}

:global(body) {
  margin: 0;
  font-family: 'Inter', sans-serif;
  background:
    radial-gradient(circle at top left, rgba(208, 227, 255, 0.4), transparent 28%),
    #f4faff;
  color: #111d23;
}

:global(.material-symbols-outlined) {
  font-variation-settings:
    'FILL' 0,
    'wght' 400,
    'GRAD' 0,
    'opsz' 24;
}

.profile-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top right, rgba(208, 227, 255, 0.32), transparent 24%),
    #f4faff;
}

.profile-page__main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 48px 24px 80px;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.profile-card,
.license-card,
.bookings-card {
  background-color: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(0, 25, 68, 0.08);
  border-radius: 24px;
  padding: 32px;
  box-shadow: 0 20px 40px rgba(0, 25, 68, 0.06);
}

.license-card {
  border-left: 4px solid #002c6d;
}

.profile-card__header {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 32px;
}

.avatar-block {
  position: relative;
  width: fit-content;
}

.avatar-block__image-wrap {
  width: 96px;
  height: 96px;
  border-radius: 22px;
  overflow: hidden;
  border: 4px solid #e3f0f8;
  box-shadow: 0 16px 32px rgba(0, 25, 68, 0.08);
}

.avatar-block__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-card__title-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.profile-card__title {
  margin: 0;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 32px;
  font-weight: 800;
  color: #001944;
}

.profile-card__subtitle {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #767683;
}

.profile-form,
.license-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 24px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-field__label {
  padding-left: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #767683;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.form-field__input {
  width: 100%;
  min-height: 52px;
  padding: 14px 16px;
  border: 1px solid transparent;
  border-radius: 14px;
  background-color: #d7e4ec;
  color: #111d23;
  font-size: 15px;
  outline: none;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.form-field__input:focus {
  background-color: #f4faff;
  border-color: #001944;
  box-shadow: 0 0 0 3px rgba(0, 88, 202, 0.1);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.section-title--space {
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.section-title__copy {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.section-title__icon {
  color: #001944;
}

.section-title__text {
  margin: 0;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 24px;
  font-weight: 800;
  color: #001944;
}

.section-title__description {
  margin: 8px 0 0;
  color: #667788;
  line-height: 1.6;
}

.categories {
  margin-top: 28px;
}

.categories__label {
  display: inline-block;
  margin-bottom: 16px;
  padding-left: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #767683;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.categories__list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.category-chip {
  cursor: pointer;
}

.category-chip__input {
  display: none;
}

.category-chip__text {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 56px;
  padding: 12px 20px;
  border-radius: 12px;
  background-color: #e3f0f8;
  color: #454652;
  font-weight: 700;
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    transform 0.2s ease;
}

.category-chip__input:checked + .category-chip__text {
  background-color: #69ff87;
  color: #002108;
}

.category-chip:hover .category-chip__text {
  transform: translateY(-1px);
}

.bookings-refresh,
.booking-item__cancel,
.bookings-state__link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 46px;
  border-radius: 14px;
  font-weight: 700;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease,
    opacity 0.2s ease,
    background-color 0.2s ease;
}

.bookings-refresh {
  padding: 0 18px;
  border: 1px solid rgba(0, 25, 68, 0.12);
  background: #ffffff;
  color: #001944;
  cursor: pointer;
}

.bookings-refresh:hover,
.booking-item__cancel:hover,
.bookings-state__link:hover {
  transform: translateY(-1px);
}

.bookings-refresh:disabled,
.booking-item__cancel:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.bookings-state {
  padding: 28px;
  border-radius: 20px;
  background: #f6f9fc;
  text-align: center;
}

.bookings-state--error {
  background: #fff3f1;
}

.bookings-state h3 {
  margin: 0;
  color: #001944;
}

.bookings-state p {
  margin: 10px 0 0;
  color: #667788;
}

.bookings-state__link {
  margin-top: 18px;
  padding: 0 18px;
  background: #001944;
  color: #fff;
  text-decoration: none;
}

.bookings-list {
  display: grid;
  gap: 18px;
}

.booking-item {
  display: grid;
  grid-template-columns: 180px minmax(0, 1fr);
  gap: 20px;
  padding: 18px;
  border-radius: 22px;
  background: linear-gradient(180deg, rgba(247, 251, 255, 0.92), rgba(236, 244, 251, 0.92));
  border: 1px solid rgba(0, 25, 68, 0.08);
}

.booking-item__media {
  display: block;
  height: 150px;
  overflow: hidden;
  border-radius: 18px;
  background: #dbe7ef;
}

.booking-item__media img,
.booking-item__placeholder {
  width: 100%;
  height: 100%;
}

.booking-item__media img {
  object-fit: cover;
}

.booking-item__placeholder {
  display: grid;
  place-items: center;
  color: #60707f;
}

.booking-item__placeholder .material-symbols-outlined {
  font-size: 44px;
}

.booking-item__body {
  display: grid;
  gap: 18px;
}

.booking-item__header,
.booking-item__footer {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.booking-item__title {
  color: #001944;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 24px;
  font-weight: 800;
  text-decoration: none;
}

.booking-item__subtitle,
.booking-item__hint {
  margin: 8px 0 0;
  color: #667788;
  line-height: 1.6;
}

.booking-status {
  display: inline-flex;
  align-items: center;
  min-height: 34px;
  padding: 0 14px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 800;
}

.booking-status--active {
  background: rgba(105, 255, 135, 0.22);
  color: #045b1f;
}

.booking-status--cancelled {
  background: rgba(255, 195, 186, 0.5);
  color: #9b2413;
}

.booking-status--completed {
  background: rgba(216, 226, 235, 0.9);
  color: #455463;
}

.booking-meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.booking-meta__item {
  padding: 14px 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(0, 25, 68, 0.06);
}

.booking-meta__item span {
  display: block;
  color: #72808d;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.booking-meta__item strong {
  display: block;
  margin-top: 8px;
  color: #001944;
}

.booking-item__cancel {
  padding: 0 18px;
  border: none;
  background: #001944;
  color: #ffffff;
  cursor: pointer;
  box-shadow: 0 12px 24px rgba(0, 25, 68, 0.14);
}

.profile-actions {
  margin-top: 8px;
  display: flex;
  justify-content: center;
  gap: 12px;
}

.profile-actions__save {
  min-width: 220px;
  min-height: 52px;
  padding: 0 24px;
  border: none;
  border-radius: 12px;
  background-color: #001944;
  color: #ffffff;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease,
    opacity 0.2s ease;
}

.profile-actions__save:hover {
  background-color: #002c6d;
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(0, 25, 68, 0.16);
}

.profile-actions__save:disabled,
.profile-actions__logout:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.profile-actions__logout {
  min-height: 52px;
  padding: 0 24px;
  border: 1px solid rgba(0, 25, 68, 0.14);
  border-radius: 12px;
  background: #ffffff;
  color: #001944;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
}

.profile-feedback {
  margin: 0;
  text-align: center;
  font-size: 14px;
  font-weight: 600;
}

.profile-feedback--success {
  color: #005d2d;
}

.profile-feedback--error {
  color: #ba1a1a;
}

@media (max-width: 900px) {
  .booking-item {
    grid-template-columns: 1fr;
  }

  .booking-item__media {
    height: 220px;
  }

  .booking-meta {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .profile-page__main {
    padding: 32px 16px 56px;
    gap: 24px;
  }

  .profile-card,
  .license-card,
  .bookings-card {
    padding: 24px;
  }

  .profile-card__header,
  .section-title--space,
  .booking-item__header,
  .booking-item__footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .profile-card__title {
    font-size: 26px;
  }

  .profile-form,
  .license-form {
    grid-template-columns: 1fr;
  }

  .profile-actions {
    flex-direction: column;
  }
}
</style>
