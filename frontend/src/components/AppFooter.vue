<template>
  <footer class="app-footer">
    <div class="app-footer__container">
      <div class="app-footer__brand">
        <div class="app-footer__logo">Three Eights Rent</div>
        <p class="app-footer__text">
          © 2026 Three Eights Rent. Все права защищены.
        </p>
      </div>

      <div class="app-footer__column">
        <h4 class="app-footer__title">О сервисе</h4>
        <nav class="app-footer__nav">
          <a class="app-footer__link" href="#">О нас</a>
          <a class="app-footer__link" href="#">Условия аренды</a>
        </nav>
      </div>

      <div class="app-footer__column">
        <h4 class="app-footer__title">Помощь</h4>
        <nav class="app-footer__nav">
          <a class="app-footer__link" href="#">Контакты</a>
          <a class="app-footer__link" href="#">Политика конфиденциальности</a>
        </nav>
      </div>

      <div class="app-footer__column">
        <h4 class="app-footer__title">Адрес</h4>
        <p class="app-footer__text">
          Казань, ул. Кремлевская 8.
        </p>
    
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
</script>

<style scoped>
.app-footer {
  margin-top: 80px;
  border-top: 1px solid #e2e8f0;
  background-color: #f8fafc;
}

.app-footer__container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 48px 24px;
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr 1fr;
  gap: 32px;
}

.app-footer__logo {
  margin-bottom: 16px;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 20px;
  font-weight: 800;
  color: #001944;
}

.app-footer__text {
  margin: 0;
  color: #64748b;
  font-size: 14px;
  line-height: 1.6;
}

.app-footer__title {
  margin: 0 0 16px;
  color: #001944;
  font-size: 16px;
  font-weight: 700;
}

.app-footer__nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.app-footer__link {
  color: #64748b;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s ease;
}

.app-footer__link:hover {
  color: #001944;
}

.app-footer__socials {
  display: flex;
  gap: 12px;
}

.app-footer__social-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 10px;
  background-color: transparent;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease;
}

.app-footer__social-btn:hover {
  background-color: #e2e8f0;
  color: #001944;
}

@media (max-width: 900px) {
  .app-footer__container {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 640px) {
  .app-footer__container {
    padding: 40px 16px;
    grid-template-columns: 1fr;
  }
}
</style>
