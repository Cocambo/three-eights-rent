package handler

import (
	"car-service/internal/dto"
	"car-service/internal/service"
)

func toCarsCatalogResponse(result service.CatalogResult) dto.CarsCatalogResponse {
	items := make([]dto.CarCatalogItemResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toCarCatalogItemResponse(item))
	}

	return dto.CarsCatalogResponse{
		Items: items,
		Pagination: dto.PaginationMeta{
			Total:  result.Total,
			Limit:  result.Limit,
			Offset: result.Offset,
		},
	}
}

func toCarCatalogItemResponse(item service.CatalogCar) dto.CarCatalogItemResponse {
	return dto.CarCatalogItemResponse{
		ID:           item.ID,
		Brand:        item.Brand,
		Model:        item.Model,
		Year:         item.Year,
		FuelType:     item.FuelType,
		Transmission: item.Transmission,
		BodyType:     item.BodyType,
		SeatsCount:   item.SeatsCount,
		PricePerDay:  item.PricePerDay,
		Purpose:      item.Purpose,
		MainImageURL: item.MainImageURL,
	}
}

func toCarDetailsResponse(result service.CarDetailsResult) dto.CarDetailsResponse {
	images := make([]dto.CarImageResponse, 0, len(result.Images))
	for _, image := range result.Images {
		images = append(images, dto.CarImageResponse{
			ID:        image.ID,
			URL:       image.URL,
			IsMain:    image.IsMain,
			SortOrder: image.SortOrder,
		})
	}

	return dto.CarDetailsResponse{
		ID:           result.ID,
		Brand:        result.Brand,
		Model:        result.Model,
		Year:         result.Year,
		FuelType:     result.FuelType,
		Transmission: result.Transmission,
		BodyType:     result.BodyType,
		Color:        result.Color,
		SeatsCount:   result.SeatsCount,
		PricePerDay:  result.PricePerDay,
		Purpose:      result.Purpose,
		Description:  result.Description,
		Images:       images,
	}
}

func toUploadedCarImageResponse(result service.UploadedCarImageResult) dto.UploadedCarImageResponse {
	return dto.UploadedCarImageResponse{
		ID:          result.ID,
		CarID:       result.CarID,
		FileName:    result.FileName,
		ContentType: result.ContentType,
		FileSize:    result.FileSize,
		IsMain:      result.IsMain,
		SortOrder:   result.SortOrder,
		CreatedAt:   result.CreatedAt,
	}
}

func toListFavoritesResponse(items []service.FavoriteItem) dto.ListFavoritesResponse {
	responseItems := make([]dto.FavoriteResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, dto.FavoriteResponse{
			CarID:   item.CarID,
			AddedAt: item.AddedAt,
			Car:     toCarCatalogItemResponse(item.Car),
		})
	}

	return dto.ListFavoritesResponse{
		Items: responseItems,
	}
}

func toBookingResponse(item service.BookingRecord) dto.BookingResponse {
	return dto.BookingResponse{
		ID:          item.ID,
		CarID:       item.CarID,
		StartDate:   item.StartDate,
		EndDate:     item.EndDate,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		CancelledAt: item.CancelledAt,
	}
}

func toListBookingsResponse(items []service.BookingHistoryItem) dto.ListBookingsResponse {
	responseItems := make([]dto.BookingHistoryItemResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, dto.BookingHistoryItemResponse{
			ID:          item.ID,
			CarID:       item.CarID,
			StartDate:   item.StartDate,
			EndDate:     item.EndDate,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			CancelledAt: item.CancelledAt,
			Car:         toCarCatalogItemResponse(item.Car),
		})
	}

	return dto.ListBookingsResponse{
		Items: responseItems,
	}
}
