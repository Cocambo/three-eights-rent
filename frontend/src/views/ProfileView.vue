<template>
  <div class="profile-page">
    <AppHeader />

    <main class="profile-page__main">
      <section class="profile-card">
        <div class="profile-card__header">
          <div class="avatar-block">
            <div class="avatar-block__image-wrap">
              <img
                class="avatar-block__image"
                :src="avatar"
                :alt="authStore.displayName"
              />
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
            <input
              id="birthDate"
              v-model="form.birthDate"
              class="form-field__input"
              type="date"
            />
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
            <input
              id="issueDate"
              v-model="form.issueDate"
              class="form-field__input"
              type="date"
            />
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
            <label
              v-for="category in categories"
              :key="category"
              class="category-chip"
            >
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
          :disabled="isSaving || authStore.isProfileLoading"
          @click="saveProfile"
        >
          {{ isSaving ? 'Сохраняем...' : 'Сохранить изменения' }}
        </button>
      </div>
    </main>

    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import avatar from '@/assets/profile_avatar.png'
import { ApiError, formatDateForInput, useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const isSaving = ref(false)
const isLoggingOut = ref(false)
const feedbackMessage = ref('')
const feedbackType = ref<'success' | 'error'>('success')

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

onMounted(() => {
  fillFormFromStore()
  void fetchProfile()
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
  background-color: #f4faff;
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
  background-color: #f4faff;
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
.license-card {
  background-color: #ffffff;
  border-radius: 16px;
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
  border-radius: 16px;
  overflow: hidden;
  border: 4px solid #e3f0f8;
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
  border-radius: 12px;
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

@media (max-width: 768px) {
  .profile-page__main {
    padding: 32px 16px 56px;
    gap: 24px;
  }

  .profile-card,
  .license-card {
    padding: 24px;
  }

  .profile-card__header {
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
