package domains

import "time"

type CarRecommendation struct {
	SourceCarID      uint      `gorm:"column:source_car_id;primaryKey;autoIncrement:false" json:"source_car_id"`
	RecommendedCarID uint      `gorm:"column:recommended_car_id;primaryKey;autoIncrement:false" json:"recommended_car_id"`
	Score            float64   `gorm:"column:score;type:double precision;not null" json:"score"`
	Rank             int       `gorm:"column:rank;not null" json:"rank"`
	CreatedAt        time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"not null" json:"updated_at"`
}

func (CarRecommendation) TableName() string {
	return "car_recommendations"
}
