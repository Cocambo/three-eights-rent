package repository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sortColumn struct {
	Table string
	Name  string
}

var catalogSortColumns = map[string]sortColumn{
	"id":            {Table: "cars", Name: "id"},
	"year":          {Table: "cars", Name: "year"},
	"price_per_day": {Table: "cars", Name: "price_per_day"},
	"created_at":    {Table: "cars", Name: "created_at"},
}

func carCatalogFilterScopes(filter CarFilter) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		scopeCatalogSearch(filter.Search),
		scopeCatalogBrand(filter.Brand),
		scopeCatalogModel(filter.Model),
		scopeCatalogYearFrom(filter.YearFrom),
		scopeCatalogYearTo(filter.YearTo),
		scopeCatalogFuelType(filter.FuelType),
		scopeCatalogTransmission(filter.Transmission),
		scopeCatalogBodyType(filter.BodyType),
		scopeCatalogSeatsMin(filter.SeatsMin),
		scopeCatalogPriceMin(filter.PriceMin),
		scopeCatalogPriceMax(filter.PriceMax),
		scopeCatalogPurpose(filter.Purpose),
	}
}

func scopeCatalogSearch(query string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query = strings.TrimSpace(query)
		if query == "" {
			return db
		}

		pattern := "%" + escapeLike(query) + "%"
		return db.Where(
			"(cars.brand ILIKE ? ESCAPE '\\' OR cars.model ILIKE ? ESCAPE '\\')",
			pattern,
			pattern,
		)
	}
}

func scopeCatalogBrand(brand string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		brand = strings.TrimSpace(brand)
		if brand == "" {
			return db
		}

		return db.Where("cars.brand ILIKE ? ESCAPE '\\'", escapeLike(brand))
	}
}

func scopeCatalogModel(model string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		model = strings.TrimSpace(model)
		if model == "" {
			return db
		}

		return db.Where("cars.model ILIKE ? ESCAPE '\\'", escapeLike(model))
	}
}

func scopeCatalogYearFrom(yearFrom *int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if yearFrom == nil {
			return db
		}

		return db.Where("cars.year >= ?", *yearFrom)
	}
}

func scopeCatalogYearTo(yearTo *int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if yearTo == nil {
			return db
		}

		return db.Where("cars.year <= ?", *yearTo)
	}
}

func scopeCatalogFuelType(fuelType string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fuelType = strings.TrimSpace(fuelType)
		if fuelType == "" {
			return db
		}

		return db.Where("cars.fuel_type = ?", fuelType)
	}
}

func scopeCatalogTransmission(transmission string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		transmission = strings.TrimSpace(transmission)
		if transmission == "" {
			return db
		}

		return db.Where("cars.transmission = ?", transmission)
	}
}

func scopeCatalogBodyType(bodyType string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		bodyType = strings.TrimSpace(bodyType)
		if bodyType == "" {
			return db
		}

		return db.Where("cars.body_type = ?", bodyType)
	}
}

func scopeCatalogSeatsMin(seatsMin *int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if seatsMin == nil {
			return db
		}

		return db.Where("cars.seats_count >= ?", *seatsMin)
	}
}

func scopeCatalogPriceMin(priceMin *int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if priceMin == nil {
			return db
		}

		return db.Where("cars.price_per_day >= ?", *priceMin)
	}
}

func scopeCatalogPriceMax(priceMax *int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if priceMax == nil {
			return db
		}

		return db.Where("cars.price_per_day <= ?", *priceMax)
	}
}

func scopeCatalogPurpose(purpose string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		purpose = strings.TrimSpace(purpose)
		if purpose == "" {
			return db
		}

		return db.Where("cars.purpose = ?", purpose)
	}
}

func scopeCatalogSort(filter CarFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		column, ok := catalogSortColumns[filter.SortBy]
		if !ok {
			column = catalogSortColumns["id"]
		}

		desc := strings.EqualFold(filter.SortOrder, "desc")

		return db.Order(clause.OrderByColumn{
			Column: clause.Column{
				Table: column.Table,
				Name:  column.Name,
			},
			Desc: desc,
		})
	}
}

func scopeCatalogPagination(filter CarFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(filter.Limit).Offset(filter.Offset)
	}
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	)

	return replacer.Replace(value)
}
