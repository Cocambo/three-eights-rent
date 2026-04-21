CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    car_id BIGINT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_bookings_car
        FOREIGN KEY (car_id)
        REFERENCES cars(id),
    CONSTRAINT chk_bookings_dates CHECK (start_date < end_date),
    CONSTRAINT chk_bookings_status CHECK (status IN ('active', 'cancelled', 'completed'))
);

CREATE INDEX IF NOT EXISTS idx_bookings_car_dates
    ON bookings(car_id, start_date, end_date);

CREATE INDEX IF NOT EXISTS idx_bookings_user_created_at
    ON bookings(user_id, created_at DESC);
