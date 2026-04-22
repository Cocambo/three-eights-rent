<template>
  <div class="home-page">
    <AppHeader />

    <main>
      <HomeRecommendationsSection v-if="authStore.isAuthenticated" />

      <section class="hero-section">
        <div class="container hero-section__layout">
          <div class="hero-copy">
            <span class="hero-copy__eyebrow">Premium car rental</span>
            <h1 class="hero-copy__title">
              Аренда авто
              <br />
              вашей мечты
            </h1>
            <p class="hero-copy__text">
              Эксклюзивный парк премиальных автомобилей для тех, кто ценит безупречный
              стиль, мощь и абсолютный комфорт в каждой поездке.
            </p>

            <div class="hero-copy__actions">
              <RouterLink class="button button--primary" :to="{ name: 'cars' }">
                Забронировать сейчас
                <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span>
              </RouterLink>

              <RouterLink class="button button--secondary" :to="{ name: 'cars' }">
                Посмотреть парк
              </RouterLink>
            </div>

            <div class="hero-stats">
              <article class="hero-stats__card">
                <strong>50+</strong>
                <span>автомобилей премиум-класса</span>
              </article>
              <article class="hero-stats__card">
                <strong>24/7</strong>
                <span>поддержка и сопровождение</span>
              </article>
              <article class="hero-stats__card">
                <strong>2026</strong>
                <span>модельный год автопарка</span>
              </article>
            </div>
          </div>

          <div class="hero-visual" aria-hidden="true">
            <div class="hero-visual__glow hero-visual__glow--primary"></div>
            <div class="hero-visual__glow hero-visual__glow--accent"></div>

            <div class="hero-visual__frame">
              <img
                class="hero-visual__image"
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuDY16pgpKLJaz3b3f1MNB6JWyU8A-hUz7GK0oVO_HrZ97ShAhEAQaXgH9WSaPh-v2nrK9hY-S6-YkHpXNomSDcOFiTbZBaT6mhrXEw305pddWQNlzVlAQE1CMum-s_ioGj_DjFhiTuiewiY_n2N75fWJvq3wZVwsc1_9O6qoFO3x9lV_RU1BFWMiKoCodQnr4J1xivwn6UoABJJ1gvsxuCZxs3XMnPA4QBuHleRxR02g1LA-oMvNmaFwDmE3d9AaB5is6MgG4wSs9uC"
                alt="Премиальный автомобиль"
              />

              <div class="hero-visual__badge hero-visual__badge--top">
                <span class="material-symbols-outlined" aria-hidden="true">workspace_premium</span>
                <div>
                  <strong>Concierge service</strong>
                  <span>Персональный менеджер на каждом этапе</span>
                </div>
              </div>

              <div class="hero-visual__badge hero-visual__badge--bottom">
                <strong>Полная страховка включена</strong>
                <span>Подача автомобиля в удобную точку города</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="advantages-section">
        <div class="container">
          <div class="section-heading">
            <span class="section-heading__eyebrow">Наши преимущества</span>
            <h2 class="section-heading__title">Почему выбирают нас</h2>
            <p class="section-heading__text">
              Продумали сервис так, чтобы каждая аренда ощущалась как персональный premium
              experience, а не просто выдача ключей.
            </p>
          </div>

          <div class="advantages-grid">
            <article class="advantage-card">
              <div class="advantage-card__icon advantage-card__icon--blue">
                <span class="material-symbols-outlined" aria-hidden="true">support_agent</span>
              </div>
              <h3>Круглосуточная поддержка</h3>
              <p>
                Ваш персональный консьерж на связи 24/7, чтобы решить любой вопрос в любой
                точке мира.
              </p>
            </article>

            <article class="advantage-card">
              <div class="advantage-card__icon advantage-card__icon--green">
                <span class="material-symbols-outlined" aria-hidden="true">verified_user</span>
              </div>
              <h3>Страховка включена</h3>
              <p>
                Полное КАСКО премиум-класса уже в стоимости. Ваше спокойствие остаётся нашим
                приоритетом.
              </p>
            </article>

            <article class="advantage-card">
              <div class="advantage-card__icon advantage-card__icon--violet">
                <span class="material-symbols-outlined" aria-hidden="true">new_releases</span>
              </div>
              <h3>Новые авто</h3>
              <p>
                Наш парк состоит из автомобилей актуального модельного года с минимальным
                пробегом и без компромиссов по оснащению.
              </p>
            </article>
          </div>
        </div>
      </section>
    </main>

    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
import HomeRecommendationsSection from '@/components/HomeRecommendationsSection.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top right, rgba(176, 198, 255, 0.34), transparent 24%),
    radial-gradient(circle at 15% 40%, rgba(105, 255, 135, 0.12), transparent 18%),
    linear-gradient(180deg, #f7fbff 0%, #eef4f8 100%);
}

.container {
  width: min(1200px, calc(100% - 32px));
  margin: 0 auto;
}

.hero-section {
  position: relative;
  overflow: hidden;
  padding: 40px 0 28px;
}

.hero-section__layout {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(320px, 0.95fr);
  align-items: center;
  gap: 40px;
}

.hero-copy {
  position: relative;
  z-index: 1;
  padding: 28px 0;
}

.hero-copy__eyebrow,
.section-heading__eyebrow {
  display: inline-flex;
  align-items: center;
  margin-bottom: 18px;
  color: #0058ca;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.hero-copy__title,
.section-heading__title,
.advantage-card h3 {
  font-family: 'Plus Jakarta Sans', sans-serif;
}

.hero-copy__title {
  margin: 0;
  color: #001944;
  font-size: clamp(3.2rem, 6vw, 5.4rem);
  font-weight: 800;
  line-height: 0.94;
  letter-spacing: -0.05em;
}

.hero-copy__text,
.section-heading__text,
.advantage-card p,
.hero-stats__card span,
.hero-visual__badge span {
  color: #556372;
  line-height: 1.7;
}

.hero-copy__text {
  max-width: 560px;
  margin: 24px 0 0;
  font-size: 18px;
}

.hero-copy__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 34px;
}

.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 56px;
  padding: 0 24px;
  border-radius: 16px;
  font-size: 15px;
  font-weight: 700;
  text-decoration: none;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease,
    background-color 0.2s ease,
    border-color 0.2s ease;
}

.button:hover {
  transform: translateY(-1px);
}

.button--primary {
  gap: 10px;
  background: #001944;
  color: #ffffff;
  box-shadow: 0 18px 30px rgba(0, 25, 68, 0.14);
}

.button--primary:hover {
  background: #002c6d;
}

.button--secondary {
  border: 1px solid rgba(0, 25, 68, 0.12);
  background: rgba(255, 255, 255, 0.72);
  color: #001944;
}

.button--secondary:hover {
  border-color: rgba(0, 25, 68, 0.22);
  background: #ffffff;
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  margin-top: 38px;
}

.hero-stats__card,
.hero-visual__badge,
.advantage-card {
  border: 1px solid rgba(0, 25, 68, 0.08);
  backdrop-filter: blur(14px);
  box-shadow: 0 20px 40px rgba(14, 40, 64, 0.06);
}

.hero-stats__card {
  padding: 18px 18px 20px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.72);
}

.hero-stats__card strong {
  display: block;
  margin-bottom: 8px;
  color: #001944;
  font-size: 28px;
  font-weight: 800;
  line-height: 1;
}

.hero-stats__card span {
  display: block;
  font-size: 14px;
}

.hero-visual {
  position: relative;
  min-height: 640px;
}

.hero-visual__glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
}

.hero-visual__glow--primary {
  top: 10%;
  right: 10%;
  width: 260px;
  height: 260px;
  background: rgba(176, 198, 255, 0.46);
}

.hero-visual__glow--accent {
  left: 10%;
  bottom: 10%;
  width: 180px;
  height: 180px;
  background: rgba(105, 255, 135, 0.18);
}

.hero-visual__frame {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 640px;
  padding: 44px 26px;
  border-radius: 36px;
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.82), rgba(227, 240, 248, 0.72)),
    rgba(255, 255, 255, 0.46);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 32px 60px rgba(0, 25, 68, 0.08);
  overflow: hidden;
}

.hero-visual__frame::before {
  content: '';
  position: absolute;
  inset: 18px;
  border-radius: 28px;
  border: 1px solid rgba(0, 25, 68, 0.06);
}

.hero-visual__image {
  position: relative;
  z-index: 1;
  width: min(100%, 720px);
  object-fit: contain;
  filter: drop-shadow(0 28px 26px rgba(0, 25, 68, 0.16));
}

.hero-visual__badge {
  position: absolute;
  z-index: 2;
  max-width: 250px;
  padding: 16px 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.82);
}

.hero-visual__badge--top {
  top: 32px;
  left: 20px;
  display: flex;
  gap: 12px;
}

.hero-visual__badge--top .material-symbols-outlined {
  color: #001944;
}

.hero-visual__badge--bottom {
  right: 20px;
  bottom: 30px;
}

.hero-visual__badge strong {
  display: block;
  color: #001944;
  font-size: 15px;
}

.hero-visual__badge span {
  display: block;
  margin-top: 4px;
  font-size: 13px;
}

.advantages-section {
  padding: 52px 0 0;
}

.section-heading {
  max-width: 620px;
  margin-bottom: 30px;
}

.section-heading__title {
  margin: 0;
  color: #001944;
  font-size: clamp(2.2rem, 4vw, 3.2rem);
  font-weight: 800;
  line-height: 1.02;
  letter-spacing: -0.04em;
}

.section-heading__text {
  margin: 16px 0 0;
  font-size: 17px;
}

.advantages-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 22px;
}

.advantage-card {
  padding: 30px;
  border-radius: 26px;
  background: rgba(255, 255, 255, 0.78);
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.advantage-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 26px 44px rgba(14, 40, 64, 0.1);
}

.advantage-card__icon {
  width: 58px;
  height: 58px;
  border-radius: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 22px;
}

.advantage-card__icon .material-symbols-outlined {
  font-variation-settings: 'FILL' 1, 'wght' 500, 'GRAD' 0, 'opsz' 24;
}

.advantage-card__icon--blue {
  background: #d9e2ff;
  color: #001944;
}

.advantage-card__icon--green {
  background: rgba(105, 255, 135, 0.32);
  color: #003912;
}

.advantage-card__icon--violet {
  background: #e0e0ff;
  color: #27308a;
}

.advantage-card h3 {
  margin: 0;
  color: #001944;
  font-size: 24px;
  font-weight: 800;
  line-height: 1.15;
  letter-spacing: -0.03em;
}

.advantage-card p {
  margin: 14px 0 0;
  font-size: 15px;
}

@media (max-width: 1100px) {
  .hero-section__layout,
  .advantages-grid {
    grid-template-columns: 1fr;
  }

  .hero-visual {
    min-height: auto;
  }

  .hero-visual__frame {
    min-height: 520px;
  }
}

@media (max-width: 768px) {
  .container {
    width: min(100% - 24px, 1200px);
  }

  .hero-section {
    padding-top: 24px;
  }

  .hero-copy__actions,
  .hero-stats {
    grid-template-columns: 1fr;
  }

  .hero-copy__actions {
    flex-direction: column;
  }

  .button {
    width: 100%;
  }

  .hero-visual__frame {
    min-height: 420px;
    padding: 28px 16px;
  }

  .hero-visual__badge {
    position: static;
    max-width: none;
    margin-top: 14px;
  }

  .hero-visual__frame {
    display: grid;
    align-content: center;
  }

  .advantage-card {
    padding: 24px;
  }
}
</style>
