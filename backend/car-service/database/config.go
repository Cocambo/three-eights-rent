package database

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (c Config) Validate() error {
	switch {
	case strings.TrimSpace(c.Host) == "":
		return fmt.Errorf("host is required")
	case strings.TrimSpace(c.Port) == "":
		return fmt.Errorf("port is required")
	case strings.TrimSpace(c.User) == "":
		return fmt.Errorf("user is required")
	case strings.TrimSpace(c.Name) == "":
		return fmt.Errorf("dbname is required")
	case strings.TrimSpace(c.SSLMode) == "":
		return fmt.Errorf("sslmode is required")
	case c.MaxIdleConns < 0:
		return fmt.Errorf("max idle conns must be >= 0")
	case c.MaxOpenConns <= 0:
		return fmt.Errorf("max open conns must be > 0")
	case c.MaxIdleConns > c.MaxOpenConns:
		return fmt.Errorf("max idle conns must be <= max open conns")
	case c.ConnMaxLifetime < 0:
		return fmt.Errorf("conn max lifetime must be >= 0")
	case c.ConnMaxIdleTime < 0:
		return fmt.Errorf("conn max idle time must be >= 0")
	default:
		return nil
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}
