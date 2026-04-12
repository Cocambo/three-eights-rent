package dto

import "time"

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=255"`
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	BirthDate string `json:"birth_date" binding:"required,datetime=2006-01-02"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	BirthDate string `json:"birth_date" binding:"required,datetime=2006-01-02"`
}

type CreateDriverLicenseCategoryRequest struct {
	CategoryCode string `json:"category_code" binding:"required,max=10"`
	IssuedAt     string `json:"issued_at" binding:"omitempty,datetime=2006-01-02"`
	ExpiresAt    string `json:"expires_at" binding:"omitempty,datetime=2006-01-02"`
}

type CreateDriverLicenseRequest struct {
	LicenseNumber string                               `json:"license_number" binding:"required,max=100"`
	IssuedAt      string                               `json:"issued_at" binding:"required,datetime=2006-01-02"`
	ExpiresAt     string                               `json:"expires_at" binding:"required,datetime=2006-01-02"`
	Categories    []CreateDriverLicenseCategoryRequest `json:"categories" binding:"required,min=1,dive"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type ProfileResponse struct {
	UserID    uint      `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time `json:"birth_date"`
}

type TokenPairResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type AuthResponse struct {
	User   UserResponse      `json:"user"`
	Tokens TokenPairResponse `json:"tokens"`
}

type DriverLicenseCategoryResponse struct {
	CategoryCode string     `json:"category_code"`
	IssuedAt     *time.Time `json:"issued_at,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

type DriverLicenseResponse struct {
	UserID        uint                            `json:"user_id"`
	LicenseNumber string                          `json:"license_number"`
	IssuedAt      time.Time                       `json:"issued_at"`
	ExpiresAt     time.Time                       `json:"expires_at"`
	Categories    []DriverLicenseCategoryResponse `json:"categories"`
}
