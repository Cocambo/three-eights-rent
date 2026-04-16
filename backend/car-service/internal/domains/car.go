package domains

import "time"

type Car struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Brand        string     `gorm:"type:varchar(100);not null" json:"brand"`
	Model        string     `gorm:"type:varchar(100);not null" json:"model"`
	Year         int        `gorm:"type:smallint;not null" json:"year"`
	FuelType     string     `gorm:"column:fuel_type;type:varchar(50);not null" json:"fuel_type"`
	Transmission string     `gorm:"type:varchar(50);not null" json:"transmission"`
	BodyType     string     `gorm:"column:body_type;type:varchar(50);not null" json:"body_type"`
	Color        string     `gorm:"type:varchar(50);not null" json:"color"`
	SeatsCount   int        `gorm:"column:seats_count;not null" json:"seats_count"`
	PricePerDay  int64      `gorm:"column:price_per_day;type:bigint;not null" json:"price_per_day"`
	Purpose      string     `gorm:"type:varchar(100);not null" json:"purpose"`
	Description  string     `gorm:"type:text;not null" json:"description"`
	CreatedAt    time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null" json:"updated_at"`
	CarImages    []CarImage `gorm:"constraint:OnDelete:CASCADE;foreignKey:CarID;references:ID" json:"car_images,omitempty"`
	Favorites    []Favorite `gorm:"constraint:OnDelete:CASCADE;foreignKey:CarID;references:ID" json:"favorites,omitempty"`
}

func (Car) TableName() string {
	return "cars"
}

type CarImage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CarID       uint      `gorm:"column:car_id;not null;index" json:"car_id"`
	BucketName  string    `gorm:"column:bucket_name;type:varchar(100);not null" json:"bucket_name"`
	ObjectKey   string    `gorm:"column:object_key;type:varchar(512);not null;uniqueIndex" json:"object_key"`
	FileName    string    `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	ContentType string    `gorm:"column:content_type;type:varchar(255);not null" json:"content_type"`
	FileSize    int64     `gorm:"column:file_size;type:bigint;not null" json:"file_size"`
	IsMain      bool      `gorm:"column:is_main;not null;default:false" json:"is_main"`
	SortOrder   int       `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	Car         Car       `gorm:"constraint:OnDelete:CASCADE;foreignKey:CarID;references:ID" json:"-"`
}

func (CarImage) TableName() string {
	return "car_images"
}

type Favorite struct {
	UserID    uint      `gorm:"column:user_id;primaryKey;autoIncrement:false" json:"user_id"`
	CarID     uint      `gorm:"column:car_id;primaryKey;autoIncrement:false" json:"car_id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	Car       Car       `gorm:"constraint:OnDelete:CASCADE;foreignKey:CarID;references:ID" json:"-"`
}

func (Favorite) TableName() string {
	return "favorites"
}
