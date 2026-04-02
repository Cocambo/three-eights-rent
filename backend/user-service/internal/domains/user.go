package domains

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Profile       UserProfile   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	DriverLicense DriverLicense `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}

type UserProfile struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_profiles_user_id"`
	FirstName string    `gorm:"type:varchar(100);not null"`
	LastName  string    `gorm:"type:varchar(100);not null"`
	BirthDate time.Time `gorm:"type:date"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type DriverLicense struct {
	ID            uint                    `gorm:"primaryKey"`
	UserID        uint                    `gorm:"not null;uniqueIndex:idx_driver_licenses_user_id"`
	LicenseNumber string                  `gorm:"type:varchar(100);not null;uniqueIndex:idx_driver_licenses_license_number"`
	IssuedAt      time.Time               `gorm:"type:date;not null"`
	ExpiresAt     time.Time               `gorm:"type:date;not null"`
	CreatedAt     time.Time               `gorm:"not null"`
	UpdatedAt     time.Time               `gorm:"not null"`
	Categories    []DriverLicenseCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:DriverLicenseID;references:ID"`
}

type DriverLicenseCategory struct {
	ID              uint      `gorm:"primaryKey"`
	DriverLicenseID uint      `gorm:"not null;index:idx_driver_license_categories_driver_license_id"`
	CategoryCode    string    `gorm:"type:varchar(10);not null"`
	IssuedAt        time.Time `gorm:"type:date"`
	ExpiresAt       time.Time `gorm:"type:date"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

func (DriverLicense) TableName() string {
	return "driver_licenses"
}

func (DriverLicenseCategory) TableName() string {
	return "driver_license_categories"
}
