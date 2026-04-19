export interface BookingRange {
    start: string
    end: string
}

export interface CarFeature {
    label: string
    value: string
    icon?: string
}

export interface Car {
    id: string
    name: string
    brand: string
    category: string
    classLabel: string
    pricePerDay: number
    rating: number
    reviewsCount: number
    seats: number
    luggage: number
    transmission: string
    fuelType: string
    bodyType: string
    year: number
    color: string
    drive: string
    engine: string
    power: string
    acceleration: string
    maxSpeed: string
    deposit: string
    location: string
    purpose: string
    description: string
    shortDescription: string
    tags: string[]
    purposes: string[]
    favorite: boolean
    image: string
    images: string[]
    features: CarFeature[]
    highlights: string[]
    included: string[]
    bookingNotice: string
    unavailableRanges: BookingRange[]
}

export const cars: Car[] = [
    {
        id: 'bmw-5-series',
        name: 'BMW 5 Series',
        brand: 'BMW',
        category: 'Бизнес-класс',
        classLabel: 'Business Class',
        pricePerDay: 12500,
        rating: 4.8,
        reviewsCount: 28,
        seats: 5,
        luggage: 3,
        transmission: 'АКПП',
        fuelType: 'Бензин',
        bodyType: 'Седан',
        year: 2024,
        color: 'Mineral White',
        drive: 'Задний привод',
        engine: '2.0 Turbo',
        power: '258 л.с.',
        acceleration: '6.2 сек до 100',
        maxSpeed: '250 км/ч',
        deposit: '50 000 ₽',
        location: 'Казань, подача по городу и в аэропорт',
        purpose: 'Деловые поездки / ежедневная аренда',
        description:
            'Сдержанный представительский седан с премиальной отделкой салона, мягкой подвеской и отличной шумоизоляцией. Подходит для встреч, трансферов и комфортных поездок по городу.',
        shortDescription: 'Комфортный седан для города, встреч и деловых поездок.',
        tags: ['Деловые'],
        purposes: ['Для повседневной езды', 'Деловые'],
        favorite: false,
        image:
            'https://lh3.googleusercontent.com/aida-public/AB6AXuATbMiiFzDVbyOH4dxedupRIw8ps3h8HWNzkD-O-ajZWiK3V5HWKpK0qkRMtsdyW0mjG2dtE-zThPXm6tmLj0akHeGHQ7a3qkbElDjMSoEQaByr2Y_6_JQKfXWf7ZYGDLgI4R17mWOnIuOvm2h9GCsAXsrxySGe8tQ5NZCROtnzzRbE7OWWqyojKhEotsNIwoBvB0_zlp44KjhliOf6e79gbFoqrOr547PggsbCYK9urjHtF5RncJCs-ONW56wr-uf0Af74IKvXiAkx',
        images: [
            'https://lh3.googleusercontent.com/aida-public/AB6AXuATbMiiFzDVbyOH4dxedupRIw8ps3h8HWNzkD-O-ajZWiK3V5HWKpK0qkRMtsdyW0mjG2dtE-zThPXm6tmLj0akHeGHQ7a3qkbElDjMSoEQaByr2Y_6_JQKfXWf7ZYGDLgI4R17mWOnIuOvm2h9GCsAXsrxySGe8tQ5NZCROtnzzRbE7OWWqyojKhEotsNIwoBvB0_zlp44KjhliOf6e79gbFoqrOr547PggsbCYK9urjHtF5RncJCs-ONW56wr-uf0Af74IKvXiAkx',
            'https://images.unsplash.com/photo-1555215695-3004980ad54e?auto=format&fit=crop&w=1200&q=80',
            'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?auto=format&fit=crop&w=1200&q=80',
            'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&w=1200&q=80',
        ],
        features: [
            { label: 'Год выпуска', value: '2024', icon: 'calendar_today' },
            { label: 'Топливо', value: 'Бензин', icon: 'local_gas_station' },
            { label: 'Трансмиссия', value: 'АКПП', icon: 'settings' },
            { label: 'Мест', value: '5', icon: 'group' },
            { label: 'Привод', value: 'Задний привод', icon: 'trip_origin' },
            { label: 'Багаж', value: '3 чемодана', icon: 'luggage' },
        ],
        highlights: [
            'Вентилируемые сиденья и премиальная акустика',
            'Идеален для трансферов и представительских задач',
            'Комфортная посадка для дальних поездок',
        ],
        included: ['Полный бак при выдаче', 'Страховка ОСАГО', 'Поддержка 24/7'],
        bookingNotice: 'Онлайн-бронирование работает в тестовом режиме и пока не отправляет заявку менеджеру.',
        unavailableRanges: [
            { start: '2026-04-22', end: '2026-04-25' },
            { start: '2026-05-03', end: '2026-05-06' },
        ],
    },
    {
        id: 'audi-q7-quattro',
        name: 'Audi Q7 Quattro',
        brand: 'Audi',
        category: 'Luxury SUV',
        classLabel: 'SUV',
        pricePerDay: 15000,
        rating: 4.9,
        reviewsCount: 34,
        seats: 7,
        luggage: 5,
        transmission: 'АКПП',
        fuelType: 'Дизель',
        bodyType: 'SUV',
        year: 2023,
        color: 'Daytona Gray',
        drive: 'Полный привод quattro',
        engine: '3.0 TDI',
        power: '286 л.с.',
        acceleration: '6.5 сек до 100',
        maxSpeed: '241 км/ч',
        deposit: '70 000 ₽',
        location: 'Казань, подача по Татарстану',
        purpose: 'Путешествия / семья / трансферы',
        description:
            'Просторный премиальный SUV для большой семьи и дальних поездок. В салоне много места, удобная посадка и продуманная конфигурация для чемоданов и ручной клади.',
        shortDescription: 'Большой премиальный SUV для семьи и дальних маршрутов.',
        tags: ['Путешествия'],
        purposes: ['Путешествия', 'Для повседневной езды'],
        favorite: true,
        image:
            'https://lh3.googleusercontent.com/aida-public/AB6AXuALlODOR1o8SfbAHgwWxeVVp_cjVkJQYbLk5mwXOjwGf2UvAK-Gv_KkbkLjHKtR-sDVEQvGvN885JNsMzn7tGKB_7A0ElcCdyYElPnnGvztbFEVtu0q55InTyYeb8PKssRH2KEfgEHb_B8TIgiC15jWdCmMl3w5BcUx9r-styygOa7sEzz2OuJQmCkHJ7q5xE7OYREeE8wjOZ208vsjjSU7BMyNgPPuxyobfCp4fFC9szgRkYYDWvJiwObWBxQq6_UNh-8SAFYsOj49',
        images: [
            'https://lh3.googleusercontent.com/aida-public/AB6AXuALlODOR1o8SfbAHgwWxeVVp_cjVkJQYbLk5mwXOjwGf2UvAK-Gv_KkbkLjHKtR-sDVEQvGvN885JNsMzn7tGKB_7A0ElcCdyYElPnnGvztbFEVtu0q55InTyYeb8PKssRH2KEfgEHb_B8TIgiC15jWdCmMl3w5BcUx9r-styygOa7sEzz2OuJQmCkHJ7q5xE7OYREeE8wjOZ208vsjjSU7BMyNgPPuxyobfCp4fFC9szgRkYYDWvJiwObWBxQq6_UNh-8SAFYsOj49',
            'https://images.unsplash.com/photo-1511919884226-fd3cad34687c?auto=format&fit=crop&w=1200&q=80',
            'https://images.unsplash.com/photo-1544636331-e26879cd4d9b?auto=format&fit=crop&w=1200&q=80',
            'https://images.unsplash.com/photo-1494976388531-d1058494cdd8?auto=format&fit=crop&w=1200&q=80',
        ],
        features: [
            { label: 'Год выпуска', value: '2023', icon: 'calendar_today' },
            { label: 'Топливо', value: 'Дизель', icon: 'local_gas_station' },
            { label: 'Трансмиссия', value: 'АКПП', icon: 'settings' },
            { label: 'Мест', value: '7', icon: 'group' },
            { label: 'Привод', value: 'quattro AWD', icon: 'trip_origin' },
            { label: 'Багаж', value: '5 чемоданов', icon: 'luggage' },
        ],
        highlights: [
            'Полноразмерный третий ряд сидений',
            'Большой багажник для поездок и семейных путешествий',
            'Уверенная управляемость на трассе и в городе',
        ],
        included: ['Расширенная страховка', 'Детское кресло по запросу', 'Поддержка 24/7'],
        bookingNotice: 'Выбор дат уже работает как интерфейс, интеграция с реальным календарем доступности будет добавлена позже.',
        unavailableRanges: [
            { start: '2026-04-26', end: '2026-04-29' },
            { start: '2026-05-08', end: '2026-05-11' },
        ],
    },
    {
        id: 'porsche-911-turbo-s',
        name: 'Porsche 911 Turbo S',
        brand: 'Porsche',
        category: 'Спорт',
        classLabel: 'Премиум класс',
        pricePerDay: 45000,
        rating: 4.9,
        reviewsCount: 42,
        seats: 2,
        luggage: 1,
        transmission: 'PDK 8-ступ.',
        fuelType: 'Бензин АИ-100',
        bodyType: 'Купе',
        year: 2024,
        color: 'Jet Black Metallic',
        drive: 'Полный AWD',
        engine: '3.8 Twin-Turbo',
        power: '650 л.с.',
        acceleration: '2.7 сек до 100',
        maxSpeed: '330 км/ч',
        deposit: '300 000 ₽',
        location: 'Казань, возможна подача в аэропорт',
        purpose: 'Бизнес / путешествия / особые события',
        description:
            'Вершина инженерной мысли и символ динамического превосходства. Этот 911 Turbo S сочетает повседневный комфорт, трековую точность и эмоции суперкара в одном автомобиле.',
        shortDescription: 'Икона спортивной аренды с мощностью суперкара и комфортом гран-турера.',
        tags: ['Premium', 'Эксклюзив'],
        purposes: ['Эксклюзив', 'Путешествия', 'Деловые'],
        favorite: false,
        image:
            'https://lh3.googleusercontent.com/aida-public/AB6AXuASwknNcUmfA4vFQUcaiqzFUiEvW_pDx0HuugvPerFSaJPTbAByW913uQiTx9rLm0MeO-2ovSTunj9cFvnjO9ufcKqkskq36HUjU9yCbdzUkB1NZ13ilYHZSZOSxhAVhDnZjWEYq5rT_cKgfVLr2YFFid79GQ_pEQrz8cfrgcWaFirrnjabolwS6Ot2Ge8OnBC6x3eqW3CG4iTKro8krws0wmZTUz89doxW8U3uGgb6XEtBHraw2PaSicN_NwJiNBoKqM1m4-t7D5Wt',
        images: [
            'https://lh3.googleusercontent.com/aida-public/AB6AXuASwknNcUmfA4vFQUcaiqzFUiEvW_pDx0HuugvPerFSaJPTbAByW913uQiTx9rLm0MeO-2ovSTunj9cFvnjO9ufcKqkskq36HUjU9yCbdzUkB1NZ13ilYHZSZOSxhAVhDnZjWEYq5rT_cKgfVLr2YFFid79GQ_pEQrz8cfrgcWaFirrnjabolwS6Ot2Ge8OnBC6x3eqW3CG4iTKro8krws0wmZTUz89doxW8U3uGgb6XEtBHraw2PaSicN_NwJiNBoKqM1m4-t7D5Wt',
            'https://lh3.googleusercontent.com/aida-public/AB6AXuCx248Q-Eo9fcOPUT-MUdNjb7f-WieRVnIIxPOyAHcPB2eB_jPE1DnZ-5LgtT9utN-oGtW-s9Mr9RYp4c28iaC8s36-lwXoKLl8NTFbCHl5oTLNyezHbBMF5SS6ocdjYRQJLKiQvbN5Usu-l0eHfG2NSonVRzXB_kCaeDyZBDhwyaiOqmXCFocyXAQvNdaXWRvLYm21lYCgHms4p3P1_ki7FW6OQSr6KbUuYB7rkk46RJUq2lpOdIltXCY1wxtQ4LhWsN6pPB4XqeRT',
            'https://lh3.googleusercontent.com/aida-public/AB6AXuACgs3etkX9QXUuIDxp5ymTHFEBBUfg3ARZzyQpo-ztIuvCL-hvyFOoqlZfZFzumdVu_6kjDhHcfcF4CBG10sGVdT8WRz2gT8DUWOC5hrGaIFR3TFHdiylJdoTGk7zoIhNMzNHdyvOemXFxwEmy6xPIP3-HSvPd8gPB67IPtZ7pWVbeqE4019EyS_c3GJX6dPK26TdiwLd06wABIeK1hUf6TWhakH76Z9v-YbzPWgZIaUKZWmCNTfzJf8ekRShKxDHuK1-v5bDsrC6K',
            'https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&w=1200&q=80',
        ],
        features: [
            { label: 'Год выпуска', value: '2024', icon: 'calendar_today' },
            { label: 'Топливо', value: 'Бензин АИ-100', icon: 'local_gas_station' },
            { label: 'Трансмиссия', value: 'PDK 8-ступ.', icon: 'settings' },
            { label: 'Мест', value: '2', icon: 'group' },
            { label: 'Мощность', value: '650 л.с.', icon: 'bolt' },
            { label: 'Разгон', value: '2.7 сек', icon: 'speed' },
        ],
        highlights: [
            'Полный привод и идеальный баланс для динамичной езды',
            'Премиальная отделка салона и спортивные кресла',
            'Подходит для вау-эффекта, съемок и особых поездок',
        ],
        included: ['Безлимит по эмоциям', 'Страховка', 'Консьерж 24/7'],
        bookingNotice: 'Форма бронирования пока работает как прототип: диапазон дат можно выбрать, а дальше будет подключен ваш сервис.',
        unavailableRanges: [
            { start: '2026-04-21', end: '2026-04-23' },
            { start: '2026-04-28', end: '2026-05-02' },
            { start: '2026-05-09', end: '2026-05-12' },
        ],
    },
    {
        id: 'mercedes-e-class',
        name: 'Mercedes E-Class',
        brand: 'Mercedes-Benz',
        category: 'Бизнес-класс',
        classLabel: 'Executive',
        pricePerDay: 11000,
        rating: 4.7,
        reviewsCount: 22,
        seats: 5,
        luggage: 3,
        transmission: 'АКПП',
        fuelType: 'Гибрид',
        bodyType: 'Седан',
        year: 2024,
        color: 'Obsidian Black',
        drive: 'Задний привод',
        engine: '2.0 Hybrid',
        power: '313 л.с.',
        acceleration: '6.1 сек до 100',
        maxSpeed: '245 км/ч',
        deposit: '60 000 ₽',
        location: 'Казань, подача по городу',
        purpose: 'Торжества / деловые встречи',
        description:
            'Современный бизнес-седан с мягким ходом, просторным вторым рядом и статусным внешним видом. Отличный выбор для свадеб, мероприятий и важных поездок.',
        shortDescription: 'Статусный седан для церемоний, встреч и делового ритма.',
        tags: ['Торжество'],
        purposes: ['Торжество', 'Деловые'],
        favorite: false,
        image:
            'https://lh3.googleusercontent.com/aida-public/AB6AXuDurqKFTSW23S7kbn-rtS3eYI_bHB_cPzFQDyjqCN_r2feZejL1LPGE9vPFDC92PHfeKWJ9_idAQPwFMTwpgi-FHbWRrruMj8ua2QAqz8Nkq-W3y0EemBYiQqnGdO-5gl4pkDqD_YLgo0XkfYIOWaQZjYvrXJigZQksmJ6G5V6g1h8S73ZSY3f2HbdvT0TLMZ8RXUgbu_ErA7Qko8pb3VfGKwICZcfA1glKzJUM0VFFhS3NWNpYsRUllWe9xe9zIA6wJKxImnvW8cnX',
        images: [
            'https://lh3.googleusercontent.com/aida-public/AB6AXuDurqKFTSW23S7kbn-rtS3eYI_bHB_cPzFQDyjqCN_r2feZejL1LPGE9vPFDC92PHfeKWJ9_idAQPwFMTwpgi-FHbWRrruMj8ua2QAqz8Nkq-W3y0EemBYiQqnGdO-5gl4pkDqD_YLgo0XkfYIOWaQZjYvrXJigZQksmJ6G5V6g1h8S73ZSY3f2HbdvT0TLMZ8RXUgbu_ErA7Qko8pb3VfGKwICZcfA1glKzJUM0VFFhS3NWNpYsRUllWe9xe9zIA6wJKxImnvW8cnX',
            'https://images.unsplash.com/photo-1617814076668-8dfc67a0e099?auto=format&fit=crop&w=1200&q=80',
            'https://images.unsplash.com/photo-1525609004556-c46c7d6cf023?auto=format&fit=crop&w=1200&q=80',
            'https://images.unsplash.com/photo-1553440569-bcc63803a83d?auto=format&fit=crop&w=1200&q=80',
        ],
        features: [
            { label: 'Год выпуска', value: '2024', icon: 'calendar_today' },
            { label: 'Топливо', value: 'Гибрид', icon: 'local_gas_station' },
            { label: 'Трансмиссия', value: 'АКПП', icon: 'settings' },
            { label: 'Мест', value: '5', icon: 'group' },
            { label: 'Привод', value: 'Задний привод', icon: 'trip_origin' },
            { label: 'Багаж', value: '3 чемодана', icon: 'luggage' },
        ],
        highlights: [
            'Мягкая подвеска и тихий салон',
            'Статусный внешний вид для мероприятий',
            'Гибридная установка для комфортной городской езды',
        ],
        included: ['Подача по городу', 'Страховка ОСАГО', 'Поддержка 24/7'],
        bookingNotice: 'Это тестовый модуль бронирования: даты проверяются только на клиенте и никуда не отправляются.',
        unavailableRanges: [
            { start: '2026-04-24', end: '2026-04-27' },
            { start: '2026-05-05', end: '2026-05-07' },
        ],
    },
]

export const findCarById = (id: string) => cars.find((car) => car.id === id)
