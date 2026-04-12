<template>
  <div class="favorites-page">
    <AppHeader />

    <main class="favorites-page__main">
      <section class="favorites-card">
        <span class="favorites-card__eyebrow">Избранное</span>
        <h1 class="favorites-card__title">Список избранных автомобилей появится здесь</h1>

        <RouterLink class="favorites-card__button" :to="{ name: 'cars' }">
          Перейти к автомобилям
        </RouterLink>
      </section>
    </main>

    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
</script>

<style scoped>
.favorites-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top right, rgba(216, 232, 246, 0.8), transparent 32%),
    #f4faff;
}

.favorites-page__main {
  min-height: calc(100vh - 72px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}

.favorites-card {
  width: 100%;
  max-width: 720px;
  padding: 40px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 48px rgba(0, 25, 68, 0.08);
  text-align: center;
}

.favorites-card__eyebrow {
  display: inline-block;
  margin-bottom: 16px;
  color: #0058ca;
  font-size: 0.85rem;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.favorites-card__title {
  margin: 0 0 16px;
  color: #001944;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 2.5rem;
  line-height: 1.05;
}

.favorites-card__text {
  margin: 0 auto;
  max-width: 520px;
  color: #55606f;
  font-size: 1rem;
  line-height: 1.7;
}

.favorites-card__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
  margin-top: 28px;
  padding: 0 20px;
  border-radius: 12px;
  background: #001944;
  color: #ffffff;
  font-size: 0.95rem;
  font-weight: 700;
  text-decoration: none;
}

@media (max-width: 640px) {
  .favorites-card {
    padding: 32px 24px;
  }

  .favorites-card__title {
    font-size: 2rem;
  }
}
</style>
