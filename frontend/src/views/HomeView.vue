<template>
<AppHeader />
  <main class="home-view">
    <section class="home-card">
      <span class="home-card__eyebrow">Three Eights Rent</span>
      <h1 class="home-card__title">Главная страница пока в разработке</h1>
      <p class="home-card__text">
        Это временная заглушка. Отсюда позже можно будет перейти к каталогу,
        бронированию и личному кабинету.
      </p>

      <div class="home-card__actions">
        <RouterLink class="home-card__button home-card__button--primary" :to="{ name: 'login' }">
          Войти
        </RouterLink>
        <RouterLink class="home-card__button home-card__button--secondary" :to="{ name: 'register' }">
          Регистрация
        </RouterLink>
      </div>
    </section>
  </main>
  <AppFooter />
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<style scoped>
:global(body) {
  margin: 0;
  font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: #f4faff;
  color: #111d23;
}

:global(#app) {
  min-height: 100vh;
}

.home-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background:
    radial-gradient(circle at top right, rgba(221, 234, 242, 0.9), transparent 35%),
    radial-gradient(circle at bottom left, rgba(149, 158, 253, 0.12), transparent 28%),
    #f4faff;
}

.home-card {
  width: 100%;
  max-width: 760px;
  padding: 48px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 40px 60px -15px rgba(0, 25, 68, 0.08);
  backdrop-filter: blur(12px);
}

.home-card__eyebrow {
  display: inline-block;
  margin-bottom: 16px;
  color: #0058ca;
  font-size: 0.85rem;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.home-card__title {
  margin: 0 0 16px;
  color: #001944;
  font-size: 3rem;
  line-height: 1;
  letter-spacing: -0.05em;
}

.home-card__text {
  max-width: 520px;
  margin: 0;
  color: #454652;
  font-size: 1rem;
  line-height: 1.7;
}

.home-card__actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
}

.home-card__button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
  padding: 0 20px;
  border-radius: 12px;
  font-size: 0.95rem;
  font-weight: 700;
  text-decoration: none;
}

.home-card__button--primary {
  color: #ffffff;
  background: #001944;
}

.home-card__button--secondary {
  color: #001944;
  border: 1px solid rgba(0, 25, 68, 0.15);
  background: #ffffff;
}

@media (max-width: 640px) {
  .home-card {
    padding: 32px 24px;
  }

  .home-card__title {
    font-size: 2.25rem;
  }

  .home-card__actions {
    flex-direction: column;
  }
}
</style>
