<template>
  <div class="cars-grid">
    <CarCatalogCard v-for="car in cars" :key="car.id" :car="car" />
  </div>
</template>

<script setup lang="ts">
import CarCatalogCard from '@/components/CarCatalogCard.vue'
import type { CarCatalogCardModel } from '@/services/cars'

defineProps<{
  cars: CarCatalogCardModel[]
}>()
</script>

<style scoped>
.cars-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 24px;
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
