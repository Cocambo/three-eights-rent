<template>
  <header class="app-header">
    <div class="app-header__container">
      <div class="app-header__logo">Three Eights Rent</div>

      <!-- <nav class="app-header__nav">
        <a class="app-header__link" href="#">Главная</a>
        <a class="app-header__link" href="#">Автопарк</a>
        <a class="app-header__link app-header__link--active" href="#">Профиль</a>
      </nav> -->

      <nav class="app-header__nav">
        <router-link
          to="/"
          class="app-header__link"
          active-class="app-header__link--active"
          exact-active-class="app-header__link--active"
        >
          Главная
        </router-link>

        <router-link
          to="/catalog"
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Автопарк
        </router-link>

        <router-link
          to="/profile"
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Профиль
        </router-link>
      </nav>

      <button class="app-header__profile-btn" type="button" aria-label="Профиль">
        <span class="material-symbols-outlined">person</span>
      </button>
    </div>
  </header>
</template>

<script setup>
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