package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"user-service/internal/domains"
	"user-service/internal/dto"
	apperrors "user-service/internal/errors"
	"user-service/internal/repository"
)

const (
	minPasswordLength = 8
	defaultBcryptCost = bcrypt.DefaultCost
	birthDateLayout   = "2006-01-02"
)

type UserServiceConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	BcryptCost    int
}

type UserService struct {
	userRepository repository.UserRepository
	accessSecret   []byte
	refreshSecret  []byte
	accessTTL      time.Duration
	refreshTTL     time.Duration
	bcryptCost     int
	now            func() time.Time
}

type tokenClaims struct {
	TokenType string `json:"token_type"`
	UserID    uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewUserService(userRepository repository.UserRepository, cfg UserServiceConfig) *UserService {
	bcryptCost := cfg.BcryptCost
	if bcryptCost == 0 {
		bcryptCost = defaultBcryptCost
	}
	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = 20 * time.Minute
	}
	if cfg.RefreshTTL <= 0 {
		cfg.RefreshTTL = 30 * 24 * time.Hour
	}

	return &UserService{
		userRepository: userRepository,
		accessSecret:   []byte(cfg.AccessSecret),
		refreshSecret:  []byte(cfg.RefreshSecret),
		accessTTL:      cfg.AccessTTL,
		refreshTTL:     cfg.RefreshTTL,
		bcryptCost:     bcryptCost,
		now:            time.Now,
	}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email, err := normalizeAndValidateEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	birthDate, err := parseBirthDate(req.BirthDate)
	if err != nil {
		return nil, err
	}
	firstName, err := validateName(req.FirstName, "first_name")
	if err != nil {
		return nil, err
	}
	lastName, err := validateName(req.LastName, "last_name")
	if err != nil {
		return nil, err
	}

	if _, err := s.userRepository.GetUserByEmail(ctx, email); err == nil {
		return nil, apperrors.ErrDuplicateEmail
	} else if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domains.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Profile: domains.UserProfile{
			FirstName: firstName,
			LastName:  lastName,
			BirthDate: birthDate,
		},
	}

	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	tokenPair, err := s.issueAndPersistTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:   toUserResponse(user),
		Tokens: *tokenPair,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email, err := normalizeAndValidateEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, validationError("password is required")
	}

	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	tokenPair, err := s.issueAndPersistTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:   toUserResponse(user),
		Tokens: *tokenPair,
	}, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (*dto.TokenPairResponse, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, validationError("refresh token is required")
	}

	claims, err := s.parseToken(refreshToken, s.refreshSecret, "refresh")
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidToken
		}
		return nil, err
	}

	if user.RefreshTokenExpiresAt == nil || user.RefreshTokenExpiresAt.Before(s.now()) {
		return nil, apperrors.ErrUnauthorized
	}
	if claims.Subject != fmt.Sprintf("%d", user.ID) {
		return nil, apperrors.ErrInvalidToken
	}

	tokenPair, err := s.issueAndPersistTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.userRepository.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return validationError("refresh token is required")
	}

	if _, err := s.parseToken(refreshToken, s.refreshSecret, "refresh"); err != nil {
		return err
	}

	return s.userRepository.DeleteRefreshToken(ctx, refreshToken)
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*dto.ProfileResponse, error) {
	if userID == 0 {
		return nil, validationError("user_id must be greater than zero")
	}

	profile, err := s.userRepository.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toProfileResponse(profile), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	if userID == 0 {
		return nil, validationError("user_id must be greater than zero")
	}

	birthDate, err := parseBirthDate(req.BirthDate)
	if err != nil {
		return nil, err
	}
	firstName, err := validateName(req.FirstName, "first_name")
	if err != nil {
		return nil, err
	}
	lastName, err := validateName(req.LastName, "last_name")
	if err != nil {
		return nil, err
	}

	profile, err := s.userRepository.GetProfile(ctx, userID)
	if err != nil {
		if !errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}
		profile = &domains.UserProfile{UserID: userID}
	}

	profile.FirstName = firstName
	profile.LastName = lastName
	profile.BirthDate = birthDate

	if err := s.userRepository.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return toProfileResponse(profile), nil
}

func (s *UserService) CreateDriverLicense(ctx context.Context, userID uint, req dto.CreateDriverLicenseRequest) (*dto.DriverLicenseResponse, error) {
	if userID == 0 {
		return nil, validationError("user_id must be greater than zero")
	}

	issuedAt, err := parseRequiredDate(req.IssuedAt, "issued_at")
	if err != nil {
		return nil, err
	}
	expiresAt, err := parseRequiredDate(req.ExpiresAt, "expires_at")
	if err != nil {
		return nil, err
	}
	if expiresAt.Before(issuedAt) {
		return nil, validationError("expires_at cannot be before issued_at")
	}

	categories := make([]domains.DriverLicenseCategory, 0, len(req.Categories))
	for _, categoryReq := range req.Categories {
		categoryCode, err := validateCategoryCode(categoryReq.CategoryCode)
		if err != nil {
			return nil, err
		}

		category := domains.DriverLicenseCategory{
			CategoryCode: categoryCode,
		}

		if strings.TrimSpace(categoryReq.IssuedAt) != "" {
			categoryIssuedAt, err := parseOptionalDate(categoryReq.IssuedAt, "categories.issued_at")
			if err != nil {
				return nil, err
			}
			category.IssuedAt = &categoryIssuedAt
		}

		if strings.TrimSpace(categoryReq.ExpiresAt) != "" {
			categoryExpiresAt, err := parseOptionalDate(categoryReq.ExpiresAt, "categories.expires_at")
			if err != nil {
				return nil, err
			}
			category.ExpiresAt = &categoryExpiresAt
			if category.IssuedAt != nil && category.ExpiresAt.Before(*category.IssuedAt) {
				return nil, validationError("categories.expires_at cannot be before categories.issued_at")
			}
		}

		categories = append(categories, category)
	}

	licenseNumber := strings.TrimSpace(req.LicenseNumber)
	if licenseNumber == "" {
		return nil, validationError("license_number is required")
	}
	if len([]rune(licenseNumber)) > 100 {
		return nil, validationError("license_number must be at most 100 characters")
	}

	license := &domains.DriverLicense{
		UserID:        userID,
		LicenseNumber: licenseNumber,
		IssuedAt:      issuedAt,
		ExpiresAt:     expiresAt,
		Categories:    categories,
	}

	if err := s.userRepository.CreateDriverLicense(ctx, license); err != nil {
		return nil, err
	}

	createdLicense, err := s.userRepository.GetDriverLicense(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toDriverLicenseResponse(createdLicense), nil
}

func (s *UserService) issueAndPersistTokens(ctx context.Context, userID uint) (*dto.TokenPairResponse, error) {
	now := s.now()
	accessToken, err := s.generateToken(userID, "access", now.Add(s.accessTTL), s.accessSecret)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := now.Add(s.refreshTTL)
	refreshToken, err := s.generateToken(userID, "refresh", refreshExpiresAt, s.refreshSecret)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.SaveRefreshToken(ctx, userID, refreshToken, refreshExpiresAt); err != nil {
		return nil, err
	}

	return &dto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    refreshExpiresAt,
	}, nil
}

func (s *UserService) generateToken(userID uint, tokenType string, expiresAt time.Time, secret []byte) (string, error) {
	claims := tokenClaims{
		TokenType: tokenType,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(s.now()),
			NotBefore: jwt.NewNumericDate(s.now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign %s token: %w", tokenType, err)
	}
	return signedToken, nil
}

func (s *UserService) parseToken(tokenString string, secret []byte, expectedType string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrInvalidToken
		}
		return secret, nil
	})
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}
	if claims.TokenType != expectedType {
		return nil, apperrors.ErrInvalidToken
	}
	if claims.UserID == 0 {
		return nil, apperrors.ErrInvalidToken
	}
	if claims.Subject != "" && claims.Subject != fmt.Sprintf("%d", claims.UserID) {
		return nil, apperrors.ErrInvalidToken
	}
	return claims, nil
}

func normalizeAndValidateEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", validationError("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", validationError("email format is invalid")
	}
	return email, nil
}

func validatePassword(password string) error {
	switch {
	case strings.TrimSpace(password) == "":
		return validationError("password is required")
	case len(password) < minPasswordLength:
		return validationError("password must be at least 8 characters long")
	default:
		return nil
	}
}

func validateName(value string, field string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", validationError(field + " is required")
	}
	if len([]rune(value)) > 100 {
		return "", validationError(field + " must be at most 100 characters")
	}
	return value, nil
}

func parseBirthDate(value string) (time.Time, error) {
	birthDate, err := parseRequiredDate(value, "birth_date")
	if err != nil {
		return time.Time{}, err
	}
	if birthDate.After(time.Now()) {
		return time.Time{}, validationError("birth_date cannot be in the future")
	}
	return birthDate, nil
}

func parseRequiredDate(value string, field string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, validationError(field + " is required")
	}

	parsed, err := time.Parse(birthDateLayout, value)
	if err != nil {
		return time.Time{}, validationError(field + " must match YYYY-MM-DD")
	}
	return parsed, nil
}

func parseOptionalDate(value string, field string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse(birthDateLayout, value)
	if err != nil {
		return time.Time{}, validationError(field + " must match YYYY-MM-DD")
	}
	return parsed, nil
}

func validateCategoryCode(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "", validationError("category_code is required")
	}
	if len([]rune(value)) > 10 {
		return "", validationError("category_code must be at most 10 characters")
	}
	return value, nil
}

func validationError(message string) error {
	return fmt.Errorf("%w: %s", apperrors.ErrValidation, message)
}

func toUserResponse(user *domains.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func toProfileResponse(profile *domains.UserProfile) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		BirthDate: profile.BirthDate,
	}
}

func toDriverLicenseResponse(license *domains.DriverLicense) *dto.DriverLicenseResponse {
	categories := make([]dto.DriverLicenseCategoryResponse, 0, len(license.Categories))
	for _, category := range license.Categories {
		categories = append(categories, dto.DriverLicenseCategoryResponse{
			CategoryCode: category.CategoryCode,
			IssuedAt:     category.IssuedAt,
			ExpiresAt:    category.ExpiresAt,
		})
	}

	return &dto.DriverLicenseResponse{
		UserID:        license.UserID,
		LicenseNumber: license.LicenseNumber,
		IssuedAt:      license.IssuedAt,
		ExpiresAt:     license.ExpiresAt,
		Categories:    categories,
	}
}
