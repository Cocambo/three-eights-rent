<template>
  <section class="catalog-page">
    <AppHeader v-model:search-query="search" :show-search="true" />

    <main class="container catalog-layout">
      <aside class="filters-card">
        <div class="filters-card__header">
          <h2>Фильтры</h2>
          <button class="filters-reset" type="button" @click="resetFilters">Сбросить</button>
        </div>

        <div class="filters-list">
          <div class="filter-group">
            <label class="filter-label" for="brand">Марка автомобиля</label>
            <select id="brand" v-model="filters.brand" class="field">
              <option value="">Все марки</option>
              <option v-for="brand in brands" :key="brand" :value="brand">
                {{ brand }}
              </option>
            </select>
          </div>

          <div class="filter-group">
            <p class="filter-label">Назначение</p>
            <label v-for="purpose in purposes" :key="purpose" class="checkbox-row">
              <input v-model="filters.purposes" type="checkbox" :value="purpose" />
              <span>{{ purpose }}</span>
            </label>
          </div>

          <div class="filter-group">
            <p class="filter-label">Тип топлива</p>
            <label v-for="fuel in fuelTypes" :key="fuel" class="checkbox-row">
              <input v-model="filters.fuelTypes" type="checkbox" :value="fuel" />
              <span>{{ fuel }}</span>
            </label>
          </div>

          <div class="filter-group">
            <p class="filter-label">Тип кузова</p>
            <div class="chips">
              <button
                type="button"
                class="chip"
                :class="{ 'chip--active': filters.bodyType === '' }"
                @click="filters.bodyType = ''"
              >
                Все
              </button>
              <button
                v-for="bodyType in bodyTypes"
                :key="bodyType"
                type="button"
                class="chip"
                :class="{ 'chip--active': filters.bodyType === bodyType }"
                @click="filters.bodyType = bodyType"
              >
                {{ bodyType }}
              </button>
            </div>
          </div>

          <div class="filter-group">
            <p class="filter-label">Количество мест</p>
            <div class="chips">
              <button
                type="button"
                class="chip"
                :class="{ 'chip--active': filters.seats === null }"
                @click="filters.seats = null"
              >
                Все
              </button>
              <button
                v-for="seat in seatsOptions"
                :key="seat.value"
                type="button"
                class="chip"
                :class="{ 'chip--active': filters.seats === seat.value }"
                @click="filters.seats = seat.value"
              >
                {{ seat.label }}
              </button>
            </div>
          </div>

          <div class="filter-group">
            <p class="filter-label">Цена в сутки (₽)</p>
            <div class="price-grid">
              <input
                v-model.number="filters.priceFrom"
                type="number"
                class="field"
                placeholder="От"
                min="0"
              />
              <input
                v-model.number="filters.priceTo"
                type="number"
                class="field"
                placeholder="До"
                min="0"
              />
            </div>
          </div>

          <div class="filter-group">
            <p class="filter-label">КПП</p>
            <label v-for="transmission in transmissions" :key="transmission" class="checkbox-row">
              <input v-model="filters.transmissions" type="checkbox" :value="transmission" />
              <span>{{ transmission }}</span>
            </label>
          </div>
        </div>

        <button class="primary-button" type="button">Применить</button>
      </aside>

      <section class="catalog-content">
        <div class="catalog-content__header">
          <div>
            <h1>Доступные автомобили</h1>
          </div>
          <span class="catalog-count">Найдено {{ filteredCars.length }} моделей</span>
        </div>

        <CarCatalogList :cars="filteredCars" />
      </section>
    </main>

    <AppFooter />
  </section>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
import CarCatalogList from '@/components/CarCatalogList.vue'
import { cars, type Car } from '@/data/cars'

type CatalogCar = Car & {
  catalogFuelType: string
  catalogTransmission: string
}

const search = ref('')

const brands = ['BMW', 'Mercedes-Benz', 'Audi', 'Porsche']
const purposes = ['Для повседневной езды', 'Торжество', 'Деловые', 'Путешествия', 'Эксклюзив']
const fuelTypes = ['Бензин', 'Дизель', 'Электро', 'Гибрид']
const bodyTypes = ['Седан', 'SUV', 'Купе']
const transmissions = ['АКПП', 'МКПП']
const seatsOptions = [
  { label: '2', value: 2 },
  { label: '4', value: 4 },
  { label: '5', value: 5 },
  { label: '7+', value: 7 },
]

const filters = reactive({
  brand: '',
  purposes: [] as string[],
  fuelTypes: [] as string[],
  bodyType: '',
  seats: null as number | null,
  priceFrom: null as number | null,
  priceTo: null as number | null,
  transmissions: [] as string[],
})

const catalogCars = computed<CatalogCar[]>(() =>
  cars.map((car) => ({
    ...car,
    catalogFuelType:
      car.fuelType === 'Бензин АИ-100'
        ? 'Бензин'
        : car.fuelType === 'Гибрид'
          ? 'Гибрид'
          : car.fuelType === 'Дизель'
            ? 'Дизель'
            : car.fuelType,
    catalogTransmission: car.transmission === 'PDK 8-ступ.' ? 'АКПП' : car.transmission,
  })),
)

const filteredCars = computed(() => {
  return catalogCars.value.filter((car) => {
    const matchesSearch =
      !search.value ||
      `${car.name} ${car.category} ${car.brand}`.toLowerCase().includes(search.value.toLowerCase())

    const matchesBrand = !filters.brand || car.brand === filters.brand

    const matchesPurposes =
      !filters.purposes.length || filters.purposes.some((purpose) => car.purposes.includes(purpose))

    const matchesFuel = !filters.fuelTypes.length || filters.fuelTypes.includes(car.catalogFuelType)

    const matchesBodyType = !filters.bodyType || car.bodyType === filters.bodyType

    const matchesSeats =
      filters.seats === null || (filters.seats === 7 ? car.seats >= 7 : car.seats === filters.seats)

    const matchesPriceFrom = filters.priceFrom === null || car.pricePerDay >= filters.priceFrom

    const matchesPriceTo = filters.priceTo === null || car.pricePerDay <= filters.priceTo

    const matchesTransmission =
      !filters.transmissions.length || filters.transmissions.includes(car.catalogTransmission)

    return (
      matchesSearch &&
      matchesBrand &&
      matchesPurposes &&
      matchesFuel &&
      matchesBodyType &&
      matchesSeats &&
      matchesPriceFrom &&
      matchesPriceTo &&
      matchesTransmission
    )
  })
})

function resetFilters() {
  search.value = ''
  filters.brand = ''
  filters.purposes = []
  filters.fuelTypes = []
  filters.bodyType = ''
  filters.seats = null
  filters.priceFrom = null
  filters.priceTo = null
  filters.transmissions = []
}
</script>

<style scoped>
.catalog-page {
  min-height: 100vh;
  background: linear-gradient(180deg, #f7fbff 0%, #eef4f8 100%);
}

.container {
  width: min(1280px, calc(100% - 32px));
  margin: 0 auto;
}

.catalog-layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 32px;
  padding: 32px 0 48px;
}

.filters-card {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(17, 29, 35, 0.08);
  box-shadow: 0 16px 40px rgba(14, 40, 64, 0.06);
}

.filters-card {
  position: sticky;
  top: 100px;
  align-self: start;
  padding: 24px;
  border-radius: 24px;
}

.filters-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 24px;
}

.filters-card__header h2,
.catalog-content__header h1 {
  margin: 0;
}

.filters-reset {
  padding: 0;
  border: 0;
  background: none;
  color: #3b63b7;
  font-weight: 600;
  cursor: pointer;
}

.filters-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.filter-label {
  margin: 0;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6f7c88;
}

.field {
  width: 100%;
  height: 46px;
  padding: 0 14px;
  border: 1px solid #d6dfe7;
  border-radius: 14px;
  outline: none;
  background: #f8fbfd;
  color: #18222d;
}

.field:focus {
  border-color: #3b63b7;
  box-shadow: 0 0 0 4px rgba(59, 99, 183, 0.12);
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #405160;
  cursor: pointer;
}

.checkbox-row input {
  width: 18px;
  height: 18px;
  accent-color: #163f77;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.chip {
  padding: 10px 14px;
  border: 1px solid #d8e1ea;
  border-radius: 999px;
  background: #f7fafc;
  color: #4e5e6d;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s ease;
}

.chip--active {
  background: #163f77;
  border-color: #163f77;
  color: #fff;
}

.price-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.primary-button {
  border: 0;
  cursor: pointer;
  transition: 0.2s ease;
}

.primary-button {
  width: 100%;
  margin-top: 24px;
  height: 52px;
  border-radius: 16px;
  background: #163f77;
  color: #fff;
  font-weight: 700;
}

.primary-button:hover,
.primary-button:hover {
  transform: translateY(-1px);
}

.catalog-content {
  min-width: 0;
}

.catalog-content__header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
}

.catalog-count {
  color: #617080;
  font-size: 14px;
}

@media (max-width: 1100px) {
  .catalog-layout {
    grid-template-columns: 1fr;
  }

  .filters-card {
    position: static;
  }
}

@media (max-width: 768px) {
  .catalog-content__header {
    flex-direction: column;
    align-items: stretch;
  }

  .price-grid {
    grid-template-columns: 1fr;
  }
}
</style>
