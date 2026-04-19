<template>
  <div class="favorites-page">
    <AppHeader />

    <main class="favorites-page__main">
      <section class="favorites-section">
    
        <div v-if="favoritesStore.isLoading" class="favorites-state">
          <h2>Загружаем избранное</h2>
          <p>Получаем сохраненные автомобили из car-service.</p>
        </div>

        <div v-else-if="favoritesStore.errorMessage" class="favorites-state favorites-state--error">
          <h2>Не удалось загрузить избранное</h2>
          <p>{{ favoritesStore.errorMessage }}</p>
          <button class="favorites-card__button" type="button" @click="loadFavorites">Повторить</button>
        </div>

        <div v-else-if="!favoritesStore.favoriteCards.length" class="favorites-state">
          <h2>Список избранного пока пуст</h2>
          <p>Перейдите в каталог и добавьте автомобили, которые хотите сохранить.</p>

          <RouterLink class="favorites-card__button" :to="{ name: 'cars' }">
            Перейти в каталог
          </RouterLink>
        </div>

        <CarCatalogList v-else :cars="favoritesStore.favoriteCards" />
      </section>
    </main>

    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
import CarCatalogList from '@/components/CarCatalogList.vue'
import { ApiError } from '@/stores/auth'
import { useFavoritesStore } from '@/stores/favorites'

const favoritesStore = useFavoritesStore()
const router = useRouter()
const route = useRoute()

async function loadFavorites() {
  try {
    await favoritesStore.loadFavorites({ force: true })
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      await router.push({
        name: 'login',
        query: {
          redirect: route.fullPath,
        },
      })
    }
  }
}

onMounted(() => {
  void loadFavorites()
})
</script>

<style scoped>
.favorites-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top right, rgba(216, 232, 246, 0.8), transparent 32%),
    #f4faff;
}

.favorites-page__main {
  padding: 32px 16px 56px;
}

.favorites-section {
  width: min(1200px, 100%);
  margin: 0 auto;
  display: grid;
  gap: 24px;
}

.favorites-section__header,
.favorites-state {
  width: 100%;
  padding: 32px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 48px rgba(0, 25, 68, 0.08);
}

.favorites-section__eyebrow {
  display: inline-block;
  margin-bottom: 16px;
  color: #0058ca;
  font-size: 0.85rem;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.favorites-section__title,
.favorites-state h2 {
  margin: 0;
  color: #001944;
  font-family: 'Plus Jakarta Sans', sans-serif;
}

.favorites-section__title {
  font-size: clamp(2rem, 3vw, 3rem);
  line-height: 1.05;
}

.favorites-section__text,
.favorites-state p {
  margin: 14px 0 0;
  max-width: 620px;
  color: #55606f;
  font-size: 1rem;
  line-height: 1.7;
}

.favorites-state {
  text-align: center;
}

.favorites-state--error {
  border: 1px solid rgba(186, 26, 26, 0.12);
}

.favorites-card__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
  margin-top: 28px;
  padding: 0 20px;
  border: none;
  border-radius: 12px;
  background: #001944;
  color: #ffffff;
  font-size: 0.95rem;
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
}

@media (max-width: 640px) {
  .favorites-page__main {
    padding: 24px 12px 40px;
  }

  .favorites-section__header,
  .favorites-state {
    padding: 24px;
  }
}
</style>
