<template>
  <header class="app-header">
    <div class="app-header__container">
      <RouterLink class="app-header__logo" :to="{ name: 'home' }">
        Three Eights Rent
      </RouterLink>

      <nav class="app-header__nav">
        <RouterLink
          :to="{ name: 'home' }"
          class="app-header__link"
          active-class="app-header__link--active"
          exact-active-class="app-header__link--active"
        >
          Главная
        </RouterLink>

        <RouterLink
          :to="{ name: 'cars' }"
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Автомобили
        </RouterLink>

        <RouterLink
          :to="{ name: 'favorites' }"
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Избранное
        </RouterLink>
      </nav>

      <div class="app-header__actions">
        <label v-if="showSearch" class="app-header__search">
          <span class="material-symbols-outlined">search</span>
          <input
            :value="searchQuery"
            type="search"
            placeholder="Поиск автомобилей"
            @input="updateSearchQuery"
          />
        </label>

        <RouterLink
          v-if="authStore.isAuthenticated"
          class="app-header__profile-btn"
          :to="{ name: 'profile' }"
          :aria-label="`Перейти в профиль ${authStore.displayName}`"
          title="Профиль"
        >
          <span class="material-symbols-outlined">person</span>
        </RouterLink>

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
import { RouterLink } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

interface Props {
  showSearch?: boolean
  searchQuery?: string
}

withDefaults(defineProps<Props>(), {
  showSearch: false,
  searchQuery: '',
})

const emit = defineEmits<{
  (event: 'update:searchQuery', value: string): void
}>()

const authStore = useAuthStore()

const updateSearchQuery = (event: Event) => {
  emit('update:searchQuery', (event.target as HTMLInputElement).value)
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
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
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
  justify-content: center;
  gap: 24px;
  flex-wrap: wrap;
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
  justify-content: flex-end;
  min-width: 0;
  flex-wrap: nowrap;
  gap: 12px;
}

.app-header__search {
  display: inline-flex;
  align-items: center;
  flex: 1 1 220px;
  gap: 10px;
  width: 100%;
  max-width: 280px;
  min-height: 42px;
  padding: 0 14px;
  border: 1px solid rgba(0, 25, 68, 0.1);
  border-radius: 999px;
  background: #ffffff;
  color: #5f6b76;
}

.app-header__search input {
  width: 100%;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  color: #001944;
}

.app-header__search input::placeholder {
  color: #7b8794;
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

@media (max-width: 768px) {
  .app-header__container {
    display: flex;
    flex-wrap: wrap;
    padding: 12px 16px;
    gap: 16px;
  }

  .app-header__nav {
    order: 3;
    width: 100%;
    gap: 16px;
  }

  .app-header__actions {
    flex-basis: 100%;
    width: 100%;
    flex-wrap: wrap;
  }

  .app-header__search {
    width: 100%;
    max-width: none;
  }
}
</style>
