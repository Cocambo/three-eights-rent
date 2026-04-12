<template>
  <div class="login-page">
    <div class="login-card">
      <section class="login-hero">
        <img
          class="login-hero__image"
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuBTb-Iv6kVlGdyNh5Q6jImMVD-65uca7c_KvHiM3x7G0WV7Q_klOOgqKvrBHmtYxrhNWl_dRRdLQ61f45NSqJoYcndWmO268HmYX-VsRBrVkK6SP3bw6ayIrhvw31Mh-6FpqBCFbZ9Dx9dswWBn7G0ceq2K-BXaZ7h3Whyk1E9MCzd_vtKpJ5NDUhkPqRMbgTE9-codsWl2vIyw9VVRnCibC3AUb8gSaQiHEn5Ww7Pc5p43CcwZ34zqcH8od7sVp22GUjbYOC_TiavG"
          alt="Luxury dark car profile"
        />

        <div class="login-hero__overlay"></div>

        <div class="login-hero__content">
          <div class="login-hero__top">
            <RouterLink class="brand-title" :to="{ name: 'home' }">
              Three Eights Rent
            </RouterLink>
          </div>
          <div class="login-hero__middle">
            <h1 class="hero-title">
              Искусство <br />
              движения.
            </h1>
            <p class="hero-text">
              Ваш доступ к эксклюзивному парку автомобилей. Начните свое
              путешествие прямо сейчас.
            </p>
          </div>
          <div class="login-hero__bottom">
            <span class="clients-text">Более 500+ довольных клиентов</span>
          </div>
        </div>

        <div class="login-hero__blur"></div>
      </section>
      <section class="login-form-section">
        <div class="login-form-wrapper">
          <div class="login-header">
            <h2 class="login-title">Вход в личный кабинет</h2>
          </div>

          <form class="login-form" @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="email" class="form-label">Email</label>

              <div class="input-wrapper">
                <span class="input-icon">@</span>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  class="form-input"
                  placeholder="example@threeights.ru"
                  autocomplete="email"
                />
              </div>
            </div>

            <div class="form-group">
              <div class="form-group__row">
                <label for="password" class="form-label">Пароль</label>
                <button
                  type="button"
                  class="link-button"
                  @click="handleForgotPassword"
                >
                  Забыли пароль?
                </button>
              </div>

              <div class="input-wrapper">
                <span class="input-icon">#</span>
                <input
                  id="password"
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  class="form-input form-input--password"
                  placeholder="••••••••"
                  autocomplete="current-password"
                />

                <button
                  type="button"
                  class="password-toggle"
                  @click="showPassword = !showPassword"
                >
                  {{ showPassword ? 'Скрыть' : 'Показать' }}
                </button>
              </div>
            </div>

            <button
              type="submit"
              class="submit-button"
              :disabled="isLoading"
            >
              <span>{{ isLoading ? 'Входим...' : 'Войти' }}</span>
              <span class="submit-button__arrow">→</span>
            </button>
          </form>

          <div class="register-block">
            <p>
              Нет аккаунта?
              <RouterLink
                class="register-link"
                :to="{ name: 'register' }"
              >
                Зарегистрироваться
              </RouterLink>
            </p>
          </div>

          <p v-if="errorMessage" class="error-message">
            {{ errorMessage }}
          </p>
        </div>
      </section>
    </div>

    <div class="page-decor page-decor--top"></div>
    <div class="page-decor page-decor--bottom"></div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { ApiError, useAuthStore } from '@/stores/auth'

const form = reactive({
  email: '',
  password: '',
})

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const showPassword = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

const resolveRedirect = () => {
  const redirect = route.query.redirect
  return typeof redirect === 'string' && redirect.startsWith('/') ? redirect : '/'
}

const handleSubmit = async () => {
  errorMessage.value = ''

  if (!form.email || !form.password) {
    errorMessage.value = 'Заполните email и пароль.'
    return
  }

  try {
    isLoading.value = true

    await authStore.login({
      email: form.email,
      password: form.password,
    })

    await router.push(resolveRedirect())
  } catch (error) {
    if (error instanceof ApiError) {
      errorMessage.value =
        error.status === 401 ? 'Логин или пароль неверный.' : error.message
      return
    }

    errorMessage.value = 'Не удалось выполнить вход.'
  } finally {
    isLoading.value = false
  }
}

const handleForgotPassword = () => {
  errorMessage.value = 'Восстановление пароля пока не подключено.'
}
</script>

<style scoped>
:global(*) {
  box-sizing: border-box;
}

:global(body) {
  margin: 0;
  font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: #f4faff;
  color: #111d23;
}

:global(#app) {
  min-height: 100vh;
}

.login-page {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  overflow: hidden;
  background: #f4faff;
}

.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 1200px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: #ffffff;
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 40px 60px -15px rgba(0, 25, 68, 0.08);
}

.login-hero {
  position: relative;
  min-height: 720px;
  background: #001944;
  overflow: hidden;
}

.login-hero__image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0.5;
  mix-blend-mode: overlay;
}

.login-hero__overlay {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(
      180deg,
      rgba(0, 25, 68, 0.2) 0%,
      rgba(0, 25, 68, 0.5) 100%
    );
}

.login-hero__content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 48px;
}

.brand-title {
  display: inline-flex;
  align-items: center;
  margin: 0;
  color: #ffffff;
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.03em;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.brand-title:hover {
  opacity: 0.82;
}

.brand-title:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.8);
  outline-offset: 6px;
  border-radius: 6px;
}

.login-hero__middle {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.hero-title {
  margin: 0;
  color: #ffffff;
  font-size: 56px;
  font-weight: 800;
  line-height: 0.95;
  letter-spacing: -0.04em;
}

.hero-text {
  max-width: 380px;
  margin: 0;
  color: rgba(217, 226, 255, 0.78);
  font-size: 18px;
  line-height: 1.6;
}

.login-hero__bottom {
  display: flex;
  align-items: center;
  gap: 16px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
}

.login-hero__blur {
  position: absolute;
  right: -80px;
  bottom: 25%;
  width: 260px;
  height: 260px;
  border-radius: 50%;
  background: rgba(0, 44, 109, 0.35);
  filter: blur(60px);
}

.login-form-section {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 56px;
  background: #ffffff;
}

.login-form-wrapper {
  width: 100%;
  max-width: 380px;
}

.login-header {
  margin-bottom: 40px;
}

.login-title {
  margin: 0 0 8px;
  color: #001944;
  font-size: 32px;
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.03em;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.form-label {
  color: #767683;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.12em;
}

.link-button {
  padding: 0;
  border: none;
  background: transparent;
  color: #0058ca;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.link-button:hover {
  color: #00429b;
}

.input-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  top: 50%;
  left: 16px;
  transform: translateY(-50%);
  color: #767683;
  font-size: 18px;
  pointer-events: none;
}

.form-input {
  width: 100%;
  min-height: 56px;
  padding: 16px 16px 16px 46px;
  border: 1px solid transparent;
  border-radius: 14px;
  background: #d7e4ec;
  color: #111d23;
  font-size: 15px;
  font-weight: 500;
  outline: none;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.form-input::placeholder {
  color: rgba(118, 118, 131, 0.7);
}

.form-input:focus {
  background: #e3f0f8;
  border-color: #0058ca;
  box-shadow: 0 0 0 4px rgba(0, 88, 202, 0.08);
}

.form-input--password {
  padding-right: 92px;
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 14px;
  transform: translateY(-50%);
  border: none;
  background: transparent;
  color: #767683;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.password-toggle:hover {
  color: #0058ca;
}

.submit-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  min-height: 56px;
  padding: 14px 20px;
  border: none;
  border-radius: 14px;
  background: #001944;
  color: #ffffff;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 8px 20px -6px rgba(0, 25, 68, 0.3);
  transition:
    transform 0.2s ease,
    background-color 0.2s ease,
    opacity 0.2s ease;
}

.submit-button:hover {
  transform: translateY(-1px);
  background: #002c6d;
}

.submit-button:active {
  transform: translateY(0);
}

.submit-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.submit-button__arrow {
  font-size: 18px;
}

.register-block {
  text-align: center;
  margin: 10px;
}

.register-block p {
  margin: 0;
  color: #454652;
  font-size: 14px;
  font-weight: 500;
}

.register-link {
  margin-left: 4px;
  padding: 0;
  border: none;
  background: transparent;
  color: #0058ca;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
}

.register-link:hover {
  text-decoration: underline;
  text-underline-offset: 4px;
}

.error-message {
  margin: 20px 0 0;
  color: #ba1a1a;
  font-size: 14px;
  font-weight: 600;
  text-align: center;
}

.page-decor {
  position: absolute;
  z-index: 0;
  border-radius: 50%;
  pointer-events: none;
}

.page-decor--top {
  top: -10%;
  right: -10%;
  width: 40vw;
  height: 40vw;
  background: rgba(221, 234, 242, 0.75);
  filter: blur(120px);
}

.page-decor--bottom {
  bottom: -10%;
  left: -10%;
  width: 30vw;
  height: 30vw;
  background: rgba(149, 158, 253, 0.1);
  filter: blur(100px);
}

@media (max-width: 960px) {
  .login-card {
    grid-template-columns: 1fr;
  }

  .login-hero {
    display: none;
  }

  .login-form-section {
    padding: 40px 24px;
  }
}

@media (max-width: 640px) {
  .login-title {
    font-size: 30px;
  }

  .login-form-wrapper {
    max-width: 100%;
  }
}
</style>
