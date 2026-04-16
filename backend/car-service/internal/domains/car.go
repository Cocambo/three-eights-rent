package domains

import "time"

type Car struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Brand          string    `gorm:"size:100;not null" json:"brand"`
	Model          string    `gorm:"size:100;not null" json:"model"`
	Year           int       `gorm:"not null" json:"year"`
	RegistrationNo string    `gorm:"size:32;uniqueIndex;not null" json:"registration_no"`
	ImageURL       string    `gorm:"size:255" json:"image_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Car) TableName() string {
	return "cars"
}
