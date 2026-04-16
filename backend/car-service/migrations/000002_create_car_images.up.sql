CREATE TABLE IF NOT EXISTS car_images (
    id BIGSERIAL PRIMARY KEY,
    car_id BIGINT NOT NULL,
    bucket_name VARCHAR(100) NOT NULL,
    object_key VARCHAR(512) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    is_main BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_car_images_object_key UNIQUE (object_key),
    CONSTRAINT chk_car_images_file_size CHECK (file_size >= 0),
    CONSTRAINT chk_car_images_sort_order CHECK (sort_order >= 0),
    CONSTRAINT fk_car_images_car
        FOREIGN KEY (car_id)
        REFERENCES cars(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_car_images_car_id_sort_order
    ON car_images(car_id, sort_order, id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_car_images_one_main_per_car
    ON car_images(car_id)
    WHERE is_main = TRUE;
