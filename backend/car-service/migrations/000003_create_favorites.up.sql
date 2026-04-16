CREATE TABLE IF NOT EXISTS favorites (
    user_id BIGINT NOT NULL,
    car_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_favorites PRIMARY KEY (user_id, car_id),
    CONSTRAINT fk_favorites_car
        FOREIGN KEY (car_id)
        REFERENCES cars(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_favorites_car_id ON favorites(car_id);
