package domains

import "time"

const (
	BookingStatusActive    = "active"
	BookingStatusCancelled = "cancelled"
	BookingStatusCompleted = "completed"
)

type Booking struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"column:user_id;not null;index" json:"user_id"`
	CarID       uint       `gorm:"column:car_id;not null;index" json:"car_id"`
	StartDate   time.Time  `gorm:"column:start_date;not null" json:"start_date"`
	EndDate     time.Time  `gorm:"column:end_date;not null" json:"end_date"`
	Status      string     `gorm:"column:status;type:varchar(20);not null;default:active" json:"status"`
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null" json:"updated_at"`
	CancelledAt *time.Time `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`
	Car         Car        `gorm:"constraint:OnDelete:RESTRICT;foreignKey:CarID;references:ID" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}
