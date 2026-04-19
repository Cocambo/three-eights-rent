<template>
  <div>
    <div class="cars-grid">
      <CarCatalogCard v-for="car in cars" :key="car.id" :car="car" />
    </div>

    <div v-if="!cars.length" class="empty-state">
      <h3>По вашему запросу ничего не найдено</h3>
      <p>Попробуйте изменить фильтры или сбросить параметры поиска.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import CarCatalogCard from '@/components/CarCatalogCard.vue'
import type { Car } from '@/data/cars'

type CatalogCar = Car & {
  catalogFuelType: string
  catalogTransmission: string
}

defineProps<{
  cars: CatalogCar[]
}>()
</script>

<style scoped>
.cars-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 24px;
}

.empty-state {
  margin-top: 24px;
  padding: 32px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(17, 29, 35, 0.08);
  text-align: center;
}

.empty-state h3 {
  margin: 0;
}

.empty-state p {
  margin: 8px 0 0;
  color: #617080;
}

@media (max-width: 1100px) {
  .cars-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 768px) {
  .cars-grid {
    grid-template-columns: 1fr;
  }
}
</style>
