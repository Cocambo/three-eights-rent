CREATE TABLE IF NOT EXISTS cars (
    id BIGSERIAL PRIMARY KEY,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year SMALLINT NOT NULL,
    fuel_type VARCHAR(50) NOT NULL,
    transmission VARCHAR(50) NOT NULL,
    body_type VARCHAR(50) NOT NULL,
    color VARCHAR(50) NOT NULL,
    seats_count INT NOT NULL,
    price_per_day BIGINT NOT NULL,
    purpose VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_cars_year CHECK (year >= 1900),
    CONSTRAINT chk_cars_seats_count CHECK (seats_count > 0),
    CONSTRAINT chk_cars_price_per_day CHECK (price_per_day >= 0)
);

CREATE INDEX IF NOT EXISTS idx_cars_catalog_brand_model ON cars(brand, model);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_year ON cars(year);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_fuel_type ON cars(fuel_type);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_transmission ON cars(transmission);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_body_type ON cars(body_type);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_color ON cars(color);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_seats_count ON cars(seats_count);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_price_per_day ON cars(price_per_day);
CREATE INDEX IF NOT EXISTS idx_cars_catalog_purpose ON cars(purpose);
