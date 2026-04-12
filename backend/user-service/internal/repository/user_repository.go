package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"user-service/internal/domains"
	apperrors "user-service/internal/errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domains.User) error
	GetUserByEmail(ctx context.Context, email string) (*domains.User, error)
	GetUserByID(ctx context.Context, id uint) (*domains.User, error)
	GetUserByRefreshToken(ctx context.Context, refreshToken string) (*domains.User, error)
	SaveRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	GetProfile(ctx context.Context, userID uint) (*domains.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *domains.UserProfile) error
	CreateDriverLicense(ctx context.Context, license *domains.DriverLicense) error
	GetDriverLicense(ctx context.Context, userID uint) (*domains.DriverLicense, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) CreateUser(ctx context.Context, user *domains.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (r *gormUserRepository) GetUserByEmail(ctx context.Context, email string) (*domains.User, error) {
	var u domains.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &u, nil
}

func (r *gormUserRepository) GetUserByID(ctx context.Context, id uint) (*domains.User, error) {
	var u domains.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &u, nil
}

func (r *gormUserRepository) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*domains.User, error) {
	var u domains.User
	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&u).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &u, nil
}

func (r *gormUserRepository) SaveRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error {
	res := r.db.WithContext(ctx).Model(&domains.User{}).Where("id = ?", userID).Updates(map[string]any{
		"refresh_token":            refreshToken,
		"refresh_token_expires_at": expiresAt,
	})
	if res.Error != nil {
		return mapDBError(res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *gormUserRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	res := r.db.WithContext(ctx).Model(&domains.User{}).Where("refresh_token = ?", refreshToken).Updates(map[string]any{
		"refresh_token":            nil,
		"refresh_token_expires_at": nil,
	})
	if res.Error != nil {
		return mapDBError(res.Error)
	}
	return nil
}

func (r *gormUserRepository) GetProfile(ctx context.Context, userID uint) (*domains.UserProfile, error) {
	var p domains.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &p, nil
}

func (r *gormUserRepository) UpdateProfile(ctx context.Context, profile *domains.UserProfile) error {
	if err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(profile).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (r *gormUserRepository) CreateDriverLicense(ctx context.Context, license *domains.DriverLicense) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domains.DriverLicense
		err := tx.Where("user_id = ?", license.UserID).First(&existing).Error
		switch {
		case err == nil:
			if err := tx.Model(&existing).Updates(map[string]any{
				"license_number": license.LicenseNumber,
				"issued_at":      license.IssuedAt,
				"expires_at":     license.ExpiresAt,
			}).Error; err != nil {
				return mapDBError(err)
			}

			if err := tx.Where("driver_license_id = ?", existing.ID).Delete(&domains.DriverLicenseCategory{}).Error; err != nil {
				return mapDBError(err)
			}

			for i := range license.Categories {
				license.Categories[i].DriverLicenseID = existing.ID
			}

			if len(license.Categories) > 0 {
				if err := tx.Create(&license.Categories).Error; err != nil {
					return mapDBError(err)
				}
			}

			license.ID = existing.ID
			return nil
		case errors.Is(err, gorm.ErrRecordNotFound):
			if err := tx.Create(license).Error; err != nil {
				return mapDBError(err)
			}
			return nil
		default:
			return mapDBError(err)
		}
	})
}

func (r *gormUserRepository) GetDriverLicense(ctx context.Context, userID uint) (*domains.DriverLicense, error) {
	var license domains.DriverLicense
	err := r.db.WithContext(ctx).
		Preload("Categories").
		Where("user_id = ?", userID).
		First(&license).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &license, nil
}

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if isUniqueEmailViolation(pgErr) {
			return apperrors.ErrDuplicateEmail
		}
		return apperrors.ErrConflict
	}
	return err
}

func isUniqueEmailViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	name := strings.ToLower(pgErr.ConstraintName)
	detail := strings.ToLower(pgErr.Detail)
	return strings.Contains(name, "email") || strings.Contains(detail, "email")
}
