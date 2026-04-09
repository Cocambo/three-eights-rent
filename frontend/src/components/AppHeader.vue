<template>
  <header class="app-header">
    <div class="app-header__container">
      <RouterLink class="app-header__logo" :to="{ name: 'home' }">
        Three Eights Rent
      </RouterLink>

      <nav class="app-header__nav">
        <RouterLink
          to="/"
          class="app-header__link"
          active-class="app-header__link--active"
          exact-active-class="app-header__link--active"
        >
          Главная
        </RouterLink>

        <RouterLink
          v-if="authStore.isAuthenticated"
          to="/profile"
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Профиль
        </RouterLink>
      </nav>

      <div class="app-header__actions">
        <template v-if="authStore.isAuthenticated">
          <RouterLink
            class="app-header__profile-btn"
            :to="{ name: 'profile' }"
            :aria-label="`Перейти в профиль ${authStore.displayName}`"
            title="Профиль"
          >
            <span class="material-symbols-outlined">person</span>
          </RouterLink>

          <button
            class="app-header__login-btn app-header__logout-btn"
            type="button"
            :disabled="isLoggingOut"
            @click="handleLogout"
          >
            {{ isLoggingOut ? 'Выходим...' : 'Выйти' }}
          </button>
        </template>

        <RouterLink
          v-else
          class="app-header__login-btn"
          :to="{ name: 'login' }"
        >
          Войти
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { ApiError, useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const isLoggingOut = ref(false)

const handleLogout = async () => {
  try {
    isLoggingOut.value = true
    await authStore.logout()

    if (route.name !== 'home') {
      await router.push({ name: 'home' })
    }
  } catch (error) {
    if (!(error instanceof ApiError)) {
      console.error('Не удалось завершить сессию:', error)
    }
  } finally {
    isLoggingOut.value = false
  }
}
</script>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(20px);
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.04);
}

.app-header__container {
  max-width: 1200px;
  min-height: 72px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.app-header__logo {
  text-decoration: none;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 22px;
  font-weight: 800;
  color: #001944;
  letter-spacing: -0.03em;
}

.app-header__nav {
  display: flex;
  align-items: center;
  gap: 24px;
}

.app-header__link {
  color: #5f6b76;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s ease;
}

.app-header__link:hover,
.app-header__link--active {
  color: #001944;
}

.app-header__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-header__profile-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-bottom: 2px solid #001944;
  background: transparent;
  color: #001944;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  text-decoration: none;
}

.app-header__login-btn {
  min-height: 40px;
  padding: 0 18px;
  border: 1px solid rgba(0, 25, 68, 0.14);
  border-radius: 999px;
  background: #ffffff;
  color: #001944;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    opacity 0.2s ease;
}

.app-header__login-btn:hover {
  background: #eef5ff;
  border-color: rgba(0, 25, 68, 0.24);
}

.app-header__login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.app-header__logout-btn {
  border: none;
  background: #001944;
  color: #ffffff;
}

@media (max-width: 768px) {
  .app-header__container {
    padding: 0 16px;
  }

  .app-header__nav {
    display: none;
  }
}
</style>
