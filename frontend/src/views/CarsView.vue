<template>
    <section class="catalog-page">
        <AppHeader v-model:search-query="search" :show-search="true" />

        <main class="container catalog-layout">
            <aside class="filters-card">
                <div class="filters-card__header">
                    <h2>Фильтры</h2>
                    <button class="filters-reset" type="button" @click="resetFilters">
                        Сбросить
                    </button>
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
                            <button type="button" class="chip" :class="{ 'chip--active': filters.bodyType === '' }"
                                @click="filters.bodyType = ''">
                                Все
                            </button>
                            <button v-for="bodyType in bodyTypes" :key="bodyType" type="button" class="chip"
                                :class="{ 'chip--active': filters.bodyType === bodyType }"
                                @click="filters.bodyType = bodyType">
                                {{ bodyType }}
                            </button>
                        </div>
                    </div>

                    <div class="filter-group">
                        <p class="filter-label">Количество мест</p>
                        <div class="chips">
                            <button type="button" class="chip" :class="{ 'chip--active': filters.seats === null }"
                                @click="filters.seats = null">
                                Все
                            </button>
                            <button v-for="seat in seatsOptions" :key="seat.value" type="button" class="chip"
                                :class="{ 'chip--active': filters.seats === seat.value }"
                                @click="filters.seats = seat.value">
                                {{ seat.label }}
                            </button>
                        </div>
                    </div>

                    <div class="filter-group">
                        <p class="filter-label">Цена в сутки (₽)</p>
                        <div class="price-grid">
                            <input v-model.number="filters.priceFrom" type="number" class="field" placeholder="От"
                                min="0" />
                            <input v-model.number="filters.priceTo" type="number" class="field" placeholder="До"
                                min="0" />
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

                <button class="primary-button" type="button">
                    Применить
                </button>
            </aside>

            <section class="catalog-content">
                <div class="catalog-content__header">
                    <div>
                        <h1>Доступные автомобили</h1>
                    </div>
                    <span class="catalog-count">Найдено {{ filteredCars.length }} моделей</span>
                </div>

                <div class="cars-grid">
                    <article v-for="car in filteredCars" :key="car.id" class="car-card">
                        <div class="car-card__image-wrap">
                            <img :src="car.image" :alt="car.name" class="car-card__image" />

                            <button class="favorite-button" type="button" aria-label="Добавить в избранное">
                                <span class="material-symbols-outlined">
                                    {{ car.favorite ? 'favorite' : 'favorite' }}
                                </span>
                            </button>

                            <div class="car-card__badges">
                                <span v-for="badge in car.tags" :key="badge" class="badge">
                                    {{ badge }}
                                </span>
                            </div>
                        </div>

                        <div class="car-card__body">
                            <div class="car-card__top">
                                <div>
                                    <h3>{{ car.name }}</h3>
                                    <p>{{ car.category }}</p>
                                </div>
                                <div class="car-card__price">
                                    <strong>{{ formatPrice(car.pricePerDay) }}</strong>
                                    <span>в сутки</span>
                                </div>
                            </div>

                            <div class="car-card__specs">
                                <div class="spec-item">
                                    <span class="material-symbols-outlined">group</span>
                                    <span>{{ car.seats }}</span>
                                </div>
                                <div class="spec-item">
                                    <span class="material-symbols-outlined">luggage</span>
                                    <span>{{ car.luggage }}</span>
                                </div>
                                <div class="spec-item">
                                    <span class="material-symbols-outlined">settings_suggest</span>
                                    <span>{{ car.transmission }}</span>
                                </div>
                            </div>

                            <button class="secondary-button" type="button">
                                Подробнее
                            </button>
                        </div>
                    </article>
                </div>

                <div v-if="!filteredCars.length" class="empty-state">
                    <h3>По вашему запросу ничего не найдено</h3>
                    <p>Попробуйте изменить фильтры или сбросить параметры поиска.</p>
                </div>
            </section>
        </main>

        <AppFooter />
    </section>
</template>

<script setup>
    import { computed, reactive, ref } from 'vue'
    import AppFooter from '@/components/AppFooter.vue'
    import AppHeader from '@/components/AppHeader.vue'

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
        { label: '7+', value: 7 }
    ]

    const filters = reactive({
        brand: '',
        purposes: [],
        fuelTypes: [],
        bodyType: '',
        seats: null,
        priceFrom: null,
        priceTo: null,
        transmissions: []
    })

    const cars = ref([
        {
            id: 1,
            name: 'BMW 5 Series',
            brand: 'BMW',
            category: 'Business Class',
            pricePerDay: 12500,
            seats: 5,
            luggage: 3,
            transmission: 'АКПП',
            fuelType: 'Бензин',
            bodyType: 'Седан',
            tags: ['Premium', 'Деловые'],
            purposes: ['Деловые'],
            favorite: false,
            image:
                'https://lh3.googleusercontent.com/aida-public/AB6AXuATbMiiFzDVbyOH4dxedupRIw8ps3h8HWNzkD-O-ajZWiK3V5HWKpK0qkRMtsdyW0mjG2dtE-zThPXm6tmLj0akHeGHQ7a3qkbElDjMSoEQaByr2Y_6_JQKfXWf7ZYGDLgI4R17mWOnIuOvm2h9GCsAXsrxySGe8tQ5NZCROtnzzRbE7OWWqyojKhEotsNIwoBvB0_zlp44KjhliOf6e79gbFoqrOr547PggsbCYK9urjHtF5RncJCs-ONW56wr-uf0Af74IKvXiAkx'
        },
        {
            id: 2,
            name: 'Audi Q7 Quattro',
            brand: 'Audi',
            category: 'Luxury SUV',
            pricePerDay: 15000,
            seats: 7,
            luggage: 5,
            transmission: 'АКПП',
            fuelType: 'Дизель',
            bodyType: 'SUV',
            tags: ['Путешествия'],
            purposes: ['Путешествия'],
            favorite: true,
            image:
                'https://lh3.googleusercontent.com/aida-public/AB6AXuALlODOR1o8SfbAHgwWxeVVp_cjVkJQYbLk5mwXOjwGf2UvAK-Gv_KkbkLjHKtR-sDVEQvGvN885JNsMzn7tGKB_7A0ElcCdyYElPnnGvztbFEVtu0q55InTyYeb8PKssRH2KEfgEHb_B8TIgiC15jWdCmMl3w5BcUx9r-styygOa7sEzz2OuJQmCkHJ7q5xE7OYREeE8wjOZ208vsjjSU7BMyNgPPuxyobfCp4fFC9szgRkYYDWvJiwObWBxQq6_UNh-8SAFYsOj49'
        },
        {
            id: 3,
            name: 'Porsche 911',
            brand: 'Porsche',
            category: 'Sport',
            pricePerDay: 45000,
            seats: 2,
            luggage: 1,
            transmission: 'АКПП',
            fuelType: 'Бензин',
            bodyType: 'Купе',
            tags: ['Эксклюзив'],
            purposes: ['Эксклюзив'],
            favorite: false,
            image:
                'https://lh3.googleusercontent.com/aida-public/AB6AXuASwknNcUmfA4vFQUcaiqzFUiEvW_pDx0HuugvPerFSaJPTbAByW913uQiTx9rLm0MeO-2ovSTunj9cFvnjO9ufcKqkskq36HUjU9yCbdzUkB1NZ13ilYHZSZOSxhAVhDnZjWEYq5rT_cKgfVLr2YFFid79GQ_pEQrz8cfrgcWaFirrnjabolwS6Ot2Ge8OnBC6x3eqW3CG4iTKro8krws0wmZTUz89doxW8U3uGgb6XEtBHraw2PaSicN_NwJiNBoKqM1m4-t7D5Wt'
        },
        {
            id: 4,
            name: 'Mercedes E-Class',
            brand: 'Mercedes-Benz',
            category: 'Business Class',
            pricePerDay: 11000,
            seats: 5,
            luggage: 3,
            transmission: 'АКПП',
            fuelType: 'Гибрид',
            bodyType: 'Седан',
            tags: ['Для повседневной езды', 'Торжество'],
            purposes: ['Для повседневной езды', 'Торжество'],
            favorite: false,
            image:
                'https://lh3.googleusercontent.com/aida-public/AB6AXuDurqKFTSW23S7kbn-rtS3eYI_bHB_cPzFQDyjqCN_r2feZejL1LPGE9vPFDC92PHfeKWJ9_idAQPwFMTwpgi-FHbWRrruMj8ua2QAqz8Nkq-W3y0EemBYiQqnGdO-5gl4pkDqD_YLgo0XkfYIOWaQZjYvrXJigZQksmJ6G5V6g1h8S73ZSY3f2HbdvT0TLMZ8RXUgbu_ErA7Qko8pb3VfGKwICZcfA1glKzJUM0VFFhS3NWNpYsRUllWe9xe9zIA6wJKxImnvW8cnX'
        }
    ])

    const filteredCars = computed(() => {
        return cars.value.filter((car) => {
            const matchesSearch = !search.value || `${car.name} ${car.category} ${car.brand}`
                .toLowerCase()
                .includes(search.value.toLowerCase())

            const matchesBrand = !filters.brand || car.brand === filters.brand

            const matchesPurposes = !filters.purposes.length ||
                filters.purposes.some((purpose) => car.purposes.includes(purpose))

            const matchesFuel = !filters.fuelTypes.length ||
                filters.fuelTypes.includes(car.fuelType)

            const matchesBodyType = !filters.bodyType || car.bodyType === filters.bodyType

            const matchesSeats = filters.seats === null ||
                (filters.seats === 7 ? car.seats >= 7 : car.seats === filters.seats)

            const matchesPriceFrom = filters.priceFrom === null || filters.priceFrom === '' ||
                car.pricePerDay >= filters.priceFrom

            const matchesPriceTo = filters.priceTo === null || filters.priceTo === '' ||
                car.pricePerDay <= filters.priceTo

            const matchesTransmission = !filters.transmissions.length ||
                filters.transmissions.includes(car.transmission)

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

    function formatPrice(price) {
        return new Intl.NumberFormat('ru-RU').format(price) + ' ₽'
    }

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

    // Здесь можно будет подключить API:
    // onMounted(async () => {
    //   const response = await fetch('/api/cars')
    //   cars.value = await response.json()
    // })
</script>

<style scoped>
    @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');
    @import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght@400');

    :global(*) {
        box-sizing: border-box;
    }

    :global(body) {
        margin: 0;
        font-family: 'Inter', sans-serif;
        background: #f4f8fb;
        color: #18222d;
    }

    :global(a) {
        color: inherit;
        text-decoration: none;
    }

    :global(button),
    :global(input),
    :global(select) {
        font: inherit;
    }

    .catalog-page {
        min-height: 100vh;
        background: linear-gradient(180deg, #f7fbff 0%, #eef4f8 100%);
    }

    .container {
        width: min(1280px, calc(100% - 32px));
        margin: 0 auto;
    }

    .topbar {
        position: sticky;
        top: 0;
        z-index: 30;
        backdrop-filter: blur(18px);
        background: rgba(255, 255, 255, 0.9);
        border-bottom: 1px solid rgba(17, 29, 35, 0.08);
    }

    .topbar__inner {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 24px;
        min-height: 76px;
    }

    .brand {
        font-size: 24px;
        font-weight: 800;
        letter-spacing: -0.04em;
        color: #11284b;
    }

    .topbar__nav {
        display: flex;
        align-items: center;
        gap: 28px;
    }

    .topbar__link {
        position: relative;
        color: #617080;
        font-weight: 500;
        transition: color 0.2s ease;
    }

    .topbar__link:hover,
    .topbar__link--active {
        color: #163f77;
    }

    .topbar__link--active::after {
        content: '';
        position: absolute;
        left: 0;
        right: 0;
        bottom: -10px;
        height: 2px;
        border-radius: 999px;
        background: #163f77;
    }

    .topbar__actions {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .search-box {
        display: flex;
        align-items: center;
        gap: 10px;
        min-width: 280px;
        padding: 0 14px;
        height: 46px;
        background: #eef3f7;
        border: 1px solid transparent;
        border-radius: 14px;
    }

    .search-box input {
        width: 100%;
        border: 0;
        outline: none;
        background: transparent;
    }

    .catalog-layout {
        display: grid;
        grid-template-columns: 320px minmax(0, 1fr);
        gap: 32px;
        padding: 32px 0 48px;
    }

    .filters-card,
    .car-card {
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
    .catalog-content__header h1,
    .car-card__top h3,
    .empty-state h3,
    .site-footer h4 {
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
    .secondary-button,
    .icon-button {
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
    .secondary-button:hover,
    .icon-button:hover,
    .favorite-button:hover {
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

    .catalog-content__header p,
    .car-card__top p,
    .footer-text,
    .empty-state p {
        margin: 8px 0 0;
        color: #617080;
    }

    .catalog-count {
        color: #617080;
        font-size: 14px;
    }

    .cars-grid {
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: 24px;
    }

    .car-card {
        overflow: hidden;
        border-radius: 24px;
    }

    .car-card__image-wrap {
        position: relative;
        height: 250px;
        overflow: hidden;
    }

    .car-card__image {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.35s ease;
    }

    .car-card:hover .car-card__image {
        transform: scale(1.04);
    }

    .favorite-button,
    .icon-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 44px;
        height: 44px;
        border-radius: 14px;
        background: rgba(255, 255, 255, 0.9);
        color: #163f77;
    }

    .favorite-button {
        position: absolute;
        top: 16px;
        right: 16px;
    }

    .icon-button--primary {
        background: #163f77;
        color: #fff;
    }

    .car-card__badges {
        position: absolute;
        left: 16px;
        right: 16px;
        bottom: 16px;
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
    }

    .badge {
        display: inline-flex;
        align-items: center;
        min-height: 28px;
        padding: 0 10px;
        border-radius: 999px;
        background: rgba(255, 255, 255, 0.92);
        color: #163f77;
        font-size: 12px;
        font-weight: 700;
    }

    .car-card__body {
        padding: 20px;
    }

    .car-card__top {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 16px;
    }

    .car-card__price {
        text-align: right;
    }

    .car-card__price strong {
        display: block;
        font-size: 22px;
        color: #11284b;
    }

    .car-card__price span {
        font-size: 12px;
        color: #617080;
    }

    .car-card__specs {
        display: flex;
        align-items: center;
        gap: 16px;
        margin: 18px 0;
        padding: 14px 0;
        border-top: 1px solid #ebf0f4;
        border-bottom: 1px solid #ebf0f4;
    }

    .spec-item {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #526170;
        font-size: 14px;
    }

    .secondary-button {
        width: 100%;
        height: 48px;
        border-radius: 14px;
        background: #edf3f7;
        color: #163f77;
        font-weight: 700;
    }

    .empty-state {
        margin-top: 24px;
        padding: 32px;
        border-radius: 24px;
        background: rgba(255, 255, 255, 0.9);
        border: 1px solid rgba(17, 29, 35, 0.08);
        text-align: center;
    }

    .site-footer {
        margin-top: 56px;
        background: #fff;
        border-top: 1px solid rgba(17, 29, 35, 0.08);
    }

    .footer-grid {
        display: grid;
        grid-template-columns: 1.4fr 1fr 1fr 1.2fr;
        gap: 24px;
        padding: 40px 0;
    }

    .brand--footer {
        margin-bottom: 12px;
    }

    .site-footer ul {
        margin: 0;
        padding: 0;
        list-style: none;
    }

    .site-footer li+li {
        margin-top: 10px;
    }

    .site-footer a {
        color: #617080;
    }

    .subscribe-box {
        display: flex;
        gap: 10px;
        margin-top: 14px;
    }

    .footer-bottom {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 20px;
        padding: 20px 0 28px;
        border-top: 1px solid #edf2f6;
        color: #7d8894;
        font-size: 13px;
    }

    .footer-icons {
        display: flex;
        align-items: center;
        gap: 16px;
    }

    @media (max-width: 1100px) {
        .catalog-layout {
            grid-template-columns: 1fr;
        }

        .filters-card {
            position: static;
        }

        .cars-grid,
        .footer-grid {
            grid-template-columns: 1fr 1fr;
        }
    }

    @media (max-width: 768px) {

        .topbar__inner,
        .catalog-content__header,
        .footer-bottom {
            flex-direction: column;
            align-items: stretch;
        }

        .topbar__nav {
            order: 3;
            justify-content: space-between;
        }

        .topbar__actions,
        .search-box {
            width: 100%;
        }

        .search-box {
            min-width: 0;
        }

        .cars-grid,
        .footer-grid,
        .price-grid {
            grid-template-columns: 1fr;
        }

        .car-card__top,
        .car-card__specs {
            flex-direction: column;
            align-items: flex-start;
        }
    }
</style>
