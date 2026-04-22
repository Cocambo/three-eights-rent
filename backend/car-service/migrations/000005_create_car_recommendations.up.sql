CREATE TABLE IF NOT EXISTS car_recommendations (
    source_car_id BIGINT NOT NULL,
    recommended_car_id BIGINT NOT NULL,
    score DOUBLE PRECISION NOT NULL,
    rank INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_car_recommendations PRIMARY KEY (source_car_id, recommended_car_id),
    CONSTRAINT fk_car_recommendations_source_car
        FOREIGN KEY (source_car_id)
        REFERENCES cars(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_car_recommendations_recommended_car
        FOREIGN KEY (recommended_car_id)
        REFERENCES cars(id)
        ON DELETE CASCADE,
    CONSTRAINT chk_car_recommendations_rank CHECK (rank > 0),
    CONSTRAINT chk_car_recommendations_distinct CHECK (source_car_id <> recommended_car_id)
);

CREATE INDEX IF NOT EXISTS idx_car_recommendations_source_car_id
    ON car_recommendations(source_car_id);
