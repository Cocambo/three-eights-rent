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
                :src="profile.avatar"
                :alt="profile.fullName"
              />
            </div>
          </div>

          <div class="profile-card__title-wrap">
            <h1 class="profile-card__title">{{ profile.fullName }}</h1>
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
              v-model="form.email"
              class="form-field__input"
              type="email"
              placeholder="Введите email"
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

      <div class="profile-actions">
        <button class="profile-actions__save" type="button" @click="saveProfile">
          Сохранить изменения
        </button>
      </div>
    </main>

    <AppFooter />
  </div>
</template>

<script setup>
import { reactive, computed, onMounted } from 'vue'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import avatar from '@/assets/profile_avatar.png'

const profile = reactive({
  firstName: 'Иван',
  lastName: 'Иванов',
  birthDate: '1992-07-18',
  email: 'example@threeeights.ru',
  avatar,
  licenseNumber: '',
  issueDate: '2018-05-12',
  expirationDate: '2028-05-12',
  categories: ['B', 'C', 'D'],
})

const categories = ['A', 'B', 'C', 'D', 'E']

const form = reactive({
  firstName: '',
  lastName: '',
  birthDate: '',
  email: '',
  licenseNumber: '',
  issueDate: '',
  expirationDate: '',
  categories: [],
})

const fullName = computed(() => `${profile.firstName} ${profile.lastName}`)

profile.fullName = fullName.value

const fillFormFromProfile = () => {
  form.firstName = profile.firstName
  form.lastName = profile.lastName
  form.birthDate = profile.birthDate
  form.email = profile.email
  form.licenseNumber = profile.licenseNumber
  form.issueDate = profile.issueDate
  form.expirationDate = profile.expirationDate
  form.categories = [...profile.categories]
}

const fetchProfile = async () => {
  try {
    // TODO: заменить на реальный API-запрос
    // const response = await api.get('/profile')
    // Object.assign(profile, response.data)

    fillFormFromProfile()
  } catch (error) {
    console.error('Ошибка при загрузке профиля:', error)
  }
}

const saveProfile = async () => {
  try {
    const payload = {
      firstName: form.firstName,
      lastName: form.lastName,
      birthDate: form.birthDate,
      email: form.email,
      licenseNumber: form.licenseNumber,
      issueDate: form.issueDate,
      expirationDate: form.expirationDate,
      categories: form.categories,
    }

    // TODO: заменить на реальный API-запрос
    // await api.put('/profile', payload)

    Object.assign(profile, payload)
    profile.fullName = `${payload.firstName} ${payload.lastName}`

    console.log('Профиль сохранен:', payload)
  } catch (error) {
    console.error('Ошибка при сохранении профиля:', error)
  }
}

onMounted(() => {
  fetchProfile()
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

.avatar-block__edit-btn {
  position: absolute;
  right: -8px;
  bottom: -8px;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 12px;
  background-color: #001944;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 10px 20px rgba(0, 25, 68, 0.18);
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
}

.profile-actions__save {
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
    box-shadow 0.2s ease;
}

.profile-actions__save:hover {
  background-color: #002c6d;
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(0, 25, 68, 0.16);
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

}
</style>
