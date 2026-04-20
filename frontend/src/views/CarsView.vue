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
            <select id="brand" v-model="draftFilters.brand" class="field">
              <option value="">Все марки</option>
              <option v-for="brand in brandOptions" :key="brand.value" :value="brand.value">
                {{ brand.label }}
              </option>
            </select>
          </div>

          <div class="filter-group">
            <p class="filter-label">Назначение</p>
            <label v-for="purpose in purposeOptions" :key="purpose.value" class="checkbox-row">
              <input
                :checked="draftFilters.purpose === purpose.value"
                type="checkbox"
                :value="purpose.value"
                @change="toggleSingleValueFilter('purpose', purpose.value)"
              />
              <span>{{ purpose.label }}</span>
            </label>
          </div>

          <div class="filter-group">
            <p class="filter-label">Тип топлива</p>
            <label v-for="fuel in fuelTypeOptions" :key="fuel.value" class="checkbox-row">
              <input
                :checked="draftFilters.fuelType === fuel.value"
                type="checkbox"
                :value="fuel.value"
                @change="toggleSingleValueFilter('fuelType', fuel.value)"
              />
              <span>{{ fuel.label }}</span>
            </label>
          </div>

          <div class="filter-group">
            <p class="filter-label">Тип кузова</p>
            <div class="chips">
              <button
                type="button"
                class="chip"
                :class="{ 'chip--active': draftFilters.bodyType === '' }"
                @click="draftFilters.bodyType = ''"
              >
                Все
              </button>
              <button
                v-for="bodyType in bodyTypeOptions"
                :key="bodyType.value"
                type="button"
                class="chip"
                :class="{ 'chip--active': draftFilters.bodyType === bodyType.value }"
                @click="draftFilters.bodyType = bodyType.value"
              >
                {{ bodyType.label }}
              </button>
            </div>
          </div>

          <div class="filter-group">
            <p class="filter-label">Количество мест</p>
            <div class="chips">
              <button
                type="button"
                class="chip"
                :class="{ 'chip--active': draftFilters.seatsCount === null }"
                @click="draftFilters.seatsCount = null"
              >
                Все
              </button>
              <button
                v-for="seat in seatsOptions"
                :key="seat.value"
                type="button"
                class="chip"
                :class="{ 'chip--active': draftFilters.seatsCount === seat.value }"
                @click="draftFilters.seatsCount = seat.value"
              >
                {{ seat.label }}
              </button>
            </div>
          </div>

          <div class="filter-group">
            <p class="filter-label">Цена в сутки (₽)</p>
            <div class="price-grid">
              <input
                v-model.number="draftFilters.priceMin"
                type="number"
                class="field"
                placeholder="От"
                min="0"
              />
              <input
                v-model.number="draftFilters.priceMax"
                type="number"
                class="field"
                placeholder="До"
                min="0"
              />
            </div>
          </div>

          <div class="filter-group">
            <p class="filter-label">КПП</p>
            <label
              v-for="transmission in transmissionOptions"
              :key="transmission.value"
              class="checkbox-row"
            >
              <input
                :checked="draftFilters.transmission === transmission.value"
                type="checkbox"
                :value="transmission.value"
                @change="toggleSingleValueFilter('transmission', transmission.value)"
              />
              <span>{{ transmission.label }}</span>
            </label>
          </div>
        </div>

        <button class="primary-button" type="button" @click="applyFilters">Применить</button>
      </aside>

      <section class="catalog-content">
        <div class="catalog-content__header">
          <div>
            <h1>Доступные автомобили</h1>
          </div>
          <span class="catalog-count">Найдено {{ pagination.total }} моделей</span>
        </div>

        <div v-if="isLoading" class="catalog-state">
          <h3>Загружаем каталог</h3>
          <p>Получаем актуальные автомобили из car-service.</p>
        </div>

        <div v-else-if="errorMessage" class="catalog-state catalog-state--error">
          <h3>Не удалось загрузить каталог</h3>
          <p>{{ errorMessage }}</p>
          <button class="state-button" type="button" @click="reloadCatalog">Повторить</button>
        </div>

        <div v-else-if="!cars.length" class="catalog-state">
          <h3>По вашему запросу ничего не найдено</h3>
          <p>Попробуйте изменить фильтры или сбросить параметры поиска.</p>
        </div>

        <CarCatalogList v-else :cars="carCards" />
      </section>
    </main>

    <AppFooter />
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'

import AppFooter from '@/components/AppFooter.vue'
import AppHeader from '@/components/AppHeader.vue'
import CarCatalogList from '@/components/CarCatalogList.vue'
import {
  getCarsCatalog,
  mapCarCatalogItemToCardModel,
  type CarCatalogItem,
  type CarsCatalogQuery,
  type PaginationMeta,
} from '@/services/cars'
import { useAuthStore } from '@/stores/auth'
import { useFavoritesStore } from '@/stores/favorites'

interface FilterDraftState {
  brand: string
  purpose: string
  fuelType: string
  bodyType: string
  seatsCount: number | null
  priceMin: number | null
  priceMax: number | null
  transmission: string
}

interface FilterOption {
  value: string
  label: string
}

const defaultPagination: PaginationMeta = {
  total: 0,
  limit: 100,
  offset: 0,
}

const seatsOptions = [
  { label: '2', value: 2 },
  { label: '4', value: 4 },
  { label: '5', value: 5 },
  { label: '7', value: 7 },
]

function createDefaultDraftFilters(): FilterDraftState {
  return {
    brand: '',
    purpose: '',
    fuelType: '',
    bodyType: '',
    seatsCount: null,
    priceMin: null,
    priceMax: null,
    transmission: '',
  }
}

const search = ref('')
const cars = ref<CarCatalogItem[]>([])
const pagination = ref<PaginationMeta>(defaultPagination)
const isLoading = ref(false)
const errorMessage = ref('')
const brandOptions = ref<FilterOption[]>([])
const purposeOptions = ref<FilterOption[]>([])
const fuelTypeOptions = ref<FilterOption[]>([])
const bodyTypeOptions = ref<FilterOption[]>([])
const transmissionOptions = ref<FilterOption[]>([])
const draftFilters = reactive<FilterDraftState>(createDefaultDraftFilters())
const appliedFilters = ref<FilterDraftState>(createDefaultDraftFilters())
const activeQuery = ref<CarsCatalogQuery>({
  limit: 100,
  offset: 0,
})
const authStore = useAuthStore()
const favoritesStore = useFavoritesStore()

let searchTimer: ReturnType<typeof setTimeout> | null = null
let currentController: AbortController | null = null

const carCards = computed(() => cars.value.map((car) => mapCarCatalogItemToCardModel(car)))

function snapshotDraftFilters(): FilterDraftState {
  return {
    brand: draftFilters.brand,
    purpose: draftFilters.purpose,
    fuelType: draftFilters.fuelType,
    bodyType: draftFilters.bodyType,
    seatsCount: draftFilters.seatsCount,
    priceMin: draftFilters.priceMin,
    priceMax: draftFilters.priceMax,
    transmission: draftFilters.transmission,
  }
}

function assignDraftFilters(nextFilters: FilterDraftState) {
  draftFilters.brand = nextFilters.brand
  draftFilters.purpose = nextFilters.purpose
  draftFilters.fuelType = nextFilters.fuelType
  draftFilters.bodyType = nextFilters.bodyType
  draftFilters.seatsCount = nextFilters.seatsCount
  draftFilters.priceMin = nextFilters.priceMin
  draftFilters.priceMax = nextFilters.priceMax
  draftFilters.transmission = nextFilters.transmission
}

function toSentenceCase(value: string) {
  if (!value) {
    return ''
  }

  return value.charAt(0).toUpperCase() + value.slice(1)
}

function toFilterOption(value: string): FilterOption {
  return {
    value,
    label: toSentenceCase(value),
  }
}

function buildUniqueFilterOptions(values: string[]) {
  return [...new Set(values.map((value) => value.trim()).filter(Boolean))]
    .sort((left, right) => left.localeCompare(right, 'ru'))
    .map((value) => toFilterOption(value))
}

function buildCatalogQuery(filters: FilterDraftState, searchQuery: string): CarsCatalogQuery {
  return {
    q: searchQuery.trim() || undefined,
    brand: filters.brand || undefined,
    fuel_type: filters.fuelType || undefined,
    transmission: filters.transmission || undefined,
    body_type: filters.bodyType || undefined,
    seats_count: filters.seatsCount ?? undefined,
    price_min: filters.priceMin ?? undefined,
    price_max: filters.priceMax ?? undefined,
    purpose: filters.purpose || undefined,
    limit: 100,
    offset: 0,
  }
}

async function loadFilterOptions() {
  try {
    const response = await getCarsCatalog({
      limit: 100,
      offset: 0,
    })

    const items = response.items
    brandOptions.value = buildUniqueFilterOptions(items.map((item) => item.brand))
    purposeOptions.value = buildUniqueFilterOptions(items.map((item) => item.purpose))
    fuelTypeOptions.value = buildUniqueFilterOptions(items.map((item) => item.fuel_type))
    bodyTypeOptions.value = buildUniqueFilterOptions(items.map((item) => item.body_type))
    transmissionOptions.value = buildUniqueFilterOptions(items.map((item) => item.transmission))
  } catch (error) {
    if (!(error instanceof DOMException && error.name === 'AbortError')) {
      console.error('Failed to load filter options', error)
    }
  }
}

async function loadCatalog(query: CarsCatalogQuery) {
  currentController?.abort()
  currentController = new AbortController()
  isLoading.value = true
  errorMessage.value = ''

  try {
    const response = await getCarsCatalog(query, currentController.signal)
    cars.value = response.items
    pagination.value = response.pagination
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      return
    }

    cars.value = []
    pagination.value = defaultPagination
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить каталог автомобилей.'
  } finally {
    isLoading.value = false
  }
}

function applyFilters() {
  appliedFilters.value = snapshotDraftFilters()
  const query = buildCatalogQuery(appliedFilters.value, search.value)
  activeQuery.value = query
  window.scrollTo({ top: 0, behavior: 'smooth' })
  void loadCatalog(query)
}

function reloadCatalog() {
  void loadCatalog(activeQuery.value)
}

function resetFilters() {
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }

  search.value = ''
  const emptyFilters = createDefaultDraftFilters()
  assignDraftFilters(emptyFilters)
  appliedFilters.value = emptyFilters
  const query = buildCatalogQuery(emptyFilters, '')
  activeQuery.value = query
  void loadCatalog(query)
}

function toggleSingleValueFilter(
  key: 'purpose' | 'fuelType' | 'transmission',
  value: string,
) {
  draftFilters[key] = draftFilters[key] === value ? '' : value
}

watch(search, () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }

  searchTimer = setTimeout(() => {
    const query = buildCatalogQuery(appliedFilters.value, search.value)
    activeQuery.value = query
    void loadCatalog(query)
  }, 300)
})

onMounted(() => {
  if (authStore.isAuthenticated) {
    void favoritesStore.ensureLoaded()
  }

  void loadFilterOptions()
  void loadCatalog(activeQuery.value)
})

onBeforeUnmount(() => {
  currentController?.abort()

  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})
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

.primary-button,
.state-button {
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
.state-button:hover {
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

.catalog-state {
  padding: 32px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(17, 29, 35, 0.08);
  box-shadow: 0 16px 40px rgba(14, 40, 64, 0.06);
  text-align: center;
}

.catalog-state--error {
  border-color: rgba(186, 26, 26, 0.16);
}

.catalog-state h3 {
  margin: 0;
}

.catalog-state p {
  margin: 10px 0 0;
  color: #617080;
}

.state-button {
  min-height: 46px;
  margin-top: 18px;
  padding: 0 18px;
  border-radius: 14px;
  background: #163f77;
  color: #fff;
  font-weight: 700;
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
