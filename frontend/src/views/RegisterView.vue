<template>
  <main class="register-view">
    <div class="register-view__bg register-view__bg--top" aria-hidden="true"></div>
    <div class="register-view__bg register-view__bg--bottom" aria-hidden="true"></div>

    <section class="register-card">
      <aside class="register-card__promo">
        <div class="register-card__promo-image"></div>

        <div class="register-card__promo-content">
          <RouterLink class="register-card__brand-link" :to="{ name: 'home' }">
            Three Eights Rent
          </RouterLink>

          <p class="register-card__description">
            Присоединяйтесь к сообществу ценителей комфорта и исключительного сервиса.
          </p>
        </div>
      </aside>

      <section class="register-card__form-section">
        <header class="register-card__header">
          <h2 class="register-card__heading">Регистрация</h2>
        </header>

        <form class="register-form" @submit.prevent="handleSubmit">
          <div class="register-form__row">
            <div class="form-field">
              <label class="form-field__label" for="firstName">Имя</label>
              <input
                id="firstName"
                v-model="form.firstName"
                class="form-field__input"
                type="text"
                placeholder="Иван"
                autocomplete="given-name"
              />
            </div>

            <div class="form-field">
              <label class="form-field__label" for="lastName">Фамилия</label>
              <input
                id="lastName"
                v-model="form.lastName"
                class="form-field__input"
                type="text"
                placeholder="Иванов"
                autocomplete="family-name"
              />
            </div>
          </div>

          <div class="form-field">
            <label class="form-field__label" for="email">Email</label>
            <div class="form-field__control form-field__control--icon">
              <span class="material-symbols-outlined form-field__icon" aria-hidden="true">mail</span>
              <input
                id="email"
                v-model="form.email"
                class="form-field__input form-field__input--with-icon"
                type="email"
                placeholder="example@threeeights.ru"
                autocomplete="email"
              />
            </div>
          </div>

          <div class="form-field">
            <label class="form-field__label" for="password">Пароль</label>
            <div class="form-field__control form-field__control--icon">
              <span class="material-symbols-outlined form-field__icon" aria-hidden="true">lock</span>
              <input
                id="password"
                v-model="form.password"
                class="form-field__input form-field__input--with-icon"
                type="password"
                placeholder="••••••••"
                autocomplete="new-password"
              />
            </div>
          </div>

          <div class="form-field">
            <label class="form-field__label" for="birthDate">Дата рождения</label>
            <div class="form-field__control form-field__control--icon">
              <span class="material-symbols-outlined form-field__icon" aria-hidden="true">calendar_today</span>
              <input
                id="birthDate"
                v-model="form.birthDate"
                class="form-field__input form-field__input--with-icon"
                type="date"
              />
            </div>
          </div>

          <div class="register-form__actions">
            <button class="button button--primary" type="submit" :disabled="isLoading">
              {{ isLoading ? 'Регистрация...' : 'Зарегистрироваться' }}
            </button>

            <div class="register-form__divider" aria-hidden="true">
              <span></span>
              <span class="register-form__divider-text">или</span>
              <span></span>
            </div>

            <RouterLink class="button button--secondary" :to="{ name: 'login' }">
              Уже есть аккаунт? Войти
            </RouterLink>
          </div>
        </form>

        <p class="register-card__policy">
          Нажимая «Зарегистрироваться», вы принимаете наши
          <a href="#">Условия использования</a>
          и
          <a href="#">Политику конфиденциальности</a>.
        </p>
      </section>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'

const isLoading = ref(false)

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  birthDate: ''
})

const handleSubmit = async () => {
  try {
    isLoading.value = true

    // TODO: заменить на реальный API-запрос
    const payload = {
      firstName: form.firstName,
      lastName: form.lastName,
      email: form.email,
      password: form.password,
      birthDate: form.birthDate
    }

    console.log('Register payload:', payload)

    // Пример будущего запроса:
    // await api.post('/auth/register', payload)
  } catch (error) {
    console.error('Ошибка регистрации:', error)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Plus+Jakarta+Sans:wght@700;800&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap');

:global(*) {
  box-sizing: border-box;
}

:global(body) {
  margin: 0;
  font-family: 'Inter', sans-serif;
  background: #f4faff;
  color: #111d23;
}

:global(#app) {
  min-height: 100vh;
}

.register-view {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  overflow: hidden;
  background: #f4faff;
}

.register-view__bg {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  pointer-events: none;
}

.register-view__bg--top {
  top: -10%;
  right: -5%;
  width: 40rem;
  height: 40rem;
  opacity: 0.5;
  background: #ddeaf2;
}

.register-view__bg--bottom {
  bottom: -10%;
  left: -5%;
  width: 30rem;
  height: 30rem;
  opacity: 0.4;
  filter: blur(80px);
  background: #e9f6fd;
}

.register-card {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 5fr 7fr;
  width: 100%;
  max-width: 1200px;
  overflow: hidden;
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 40px 60px -15px rgba(0, 25, 68, 0.08);
}

.register-card__promo {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  min-height: 720px;
  padding: 40px;
  background: #001944;
}

.register-card__promo-image {
  position: absolute;
  inset: 0;
  opacity: 0.4;
  background-image: url('https://lh3.googleusercontent.com/aida-public/AB6AXuDkfYPZTJJm0qB6qMWHq-HdrhIpjIDY6jUhfURSvn70zN7OK6nVSAK2jFgpyR_8nuof0p4vrcBYQeGJS3gV6WhtZRMZvUduc9t4LClu7YyHUbIRDELR5PhR9A7KYl28VpvPv-mwSl4Iazul-DZ6wjiN6QoilZSJQc25MZ4gOb6VpZWpLTlO3fHcHBm1JSKJWJxejio3nh8Q2McHQSigy9Kq_8ODU4NQzpTnjWGePXzvmFk2cHSW0DeRR81ESmnFvzu4meplrJ_o-yuG');
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}

.register-card__promo-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 14px;
}

.register-card__brand-link {
  display: inline-flex;
  align-items: center;
  color: #ffffff;
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.03em;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.register-card__brand-link:hover {
  opacity: 0.82;
}

.register-card__brand-link:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.8);
  outline-offset: 6px;
  border-radius: 6px;
}

.register-card__description {
  max-width: 380px;
  margin: 0;
  color: rgba(217, 226, 255, 0.78);
  font-size: 18px;
  line-height: 1.6;
}

.register-card__form-section {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 28px 36px;
}

.register-card__header {
  margin-bottom: 22px;
}

.register-card__heading {
  margin: 0 0 8px;
  color: #001944;
  font-size: 32px;
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.03em;
}

.register-card__subheading {
  margin: 0;
  font-size: 0.84rem;
  line-height: 1.45;
  color: #454652;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.register-form__row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-field__label {
  padding-left: 4px;
  font-size: 0.68rem;
  font-weight: 600;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #001944;
}

.form-field__control {
  position: relative;
}

.form-field__icon,
.material-symbols-outlined {
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
}

.form-field__icon {
  position: absolute;
  top: 50%;
  left: 14px;
  transform: translateY(-50%);
  font-size: 18px;
  color: #767683;
}

.form-field__input {
  width: 100%;
  min-height: 44px;
  padding: 10px 12px;
  border: 2px solid transparent;
  border-radius: 10px;
  outline: none;
  font: inherit;
  font-size: 0.92rem;
  color: #111d23;
  background: #d7e4ec;
  transition: border-color 0.2s ease, background-color 0.2s ease;
}

.form-field__input::placeholder {
  color: #767683;
}

.form-field__input:focus {
  border-color: #001944;
  background: #e9f6fd;
}

.form-field__input--with-icon {
  padding-left: 40px;
}

.register-form__actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-top: 2px;
}

.register-form__divider {
  display: flex;
  align-items: center;
  gap: 12px;
}

.register-form__divider span:first-child,
.register-form__divider span:last-child {
  flex: 1;
  height: 1px;
  background: #d7e4ec;
}

.register-form__divider-text {
  font-size: 0.68rem;
  font-weight: 500;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #767683;
}

.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 44px;
  padding: 10px 14px;
  border-radius: 10px;
  border: none;
  font: inherit;
  font-size: 0.92rem;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.25s ease;
}

.button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.button--primary {
  color: #ffffff;
  background: #001944;
  box-shadow: 0 10px 24px rgba(0, 25, 68, 0.15);
}

.button--primary:hover:not(:disabled) {
  background: #002c6d;
}

.button--secondary {
  color: #001944;
  background: transparent;
  border: 2px solid #c6c5d4;
}

.button--secondary:hover {
  background: #e9f6fd;
}

.register-card__policy {
  margin: 18px 0 0;
  text-align: center;
  font-size: 0.68rem;
  line-height: 1.5;
  color: #767683;
}

.register-card__policy a {
  color: inherit;
}

.register-card__policy a:hover {
  color: #001944;
}

@media (max-width: 960px) {
  .register-card {
    grid-template-columns: 1fr;
    max-width: 720px;
  }

  .register-card__promo {
    display: none;
  }

  .register-card__form-section {
    padding: 24px 24px;
  }
}

@media (max-width: 640px) {
  .register-view {
    padding: 12px;
  }

  .register-card__form-section {
    padding: 24px 18px;
  }

  .register-form__row {
    grid-template-columns: 1fr;
  }

  .register-card__heading {
    font-size: 1.75rem;
  }
}
</style>
