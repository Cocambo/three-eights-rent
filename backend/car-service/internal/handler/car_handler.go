package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"car-service/internal/dto"
	apperrors "car-service/internal/errors"
	"car-service/internal/service"

	"github.com/gin-gonic/gin"
)

const maxCarImageMultipartRequestSize int64 = 12 << 20

type CarHandler struct {
	carService service.CarService
}

func NewCarHandler(carService service.CarService) *CarHandler {
	return &CarHandler{carService: carService}
}

func (h *CarHandler) Health(c *gin.Context) {
	writeSuccess(c, http.StatusOK, dto.HealthResponse{
		Status: "ok",
	})
}

func (h *CarHandler) List(c *gin.Context) {
	var req dto.ListCarsQuery
	if !bindQuery(c, &req) {
		return
	}

	catalog, err := h.carService.GetCatalog(c.Request.Context(), service.CatalogFilter{
		Search:       req.Q,
		Brand:        req.Brand,
		Model:        req.Model,
		YearFrom:     req.YearFrom,
		YearTo:       req.YearTo,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		SeatsMin:     req.SeatsMin,
		PriceMin:     req.PriceMin,
		PriceMax:     req.PriceMax,
		Purpose:      req.Purpose,
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
		Limit:        req.Limit,
		Offset:       req.Offset,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toCarsCatalogResponse(catalog))
}

func (h *CarHandler) GetByID(c *gin.Context) {
	var req dto.GetCarByIDURI
	if !bindURI(c, &req) {
		return
	}

	car, err := h.carService.GetCarDetails(c.Request.Context(), req.ID)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toCarDetailsResponse(car))
}

func (h *CarHandler) UploadImage(c *gin.Context) {
	var uri dto.GetCarByIDURI
	if !bindURI(c, &uri) {
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxCarImageMultipartRequestSize)

	var form dto.UploadCarImageForm
	if !bindForm(c, &form) {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		writeError(c, mapUploadFileError(err))
		return
	}

	uploadedImage, err := h.carService.UploadCarImage(c.Request.Context(), service.UploadCarImageCommand{
		CarID:      uri.ID,
		FileHeader: fileHeader,
		IsMain:     form.IsMain,
		SortOrder:  form.SortOrder,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusCreated, toUploadedCarImageResponse(uploadedImage))
}

func mapUploadFileError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "http: request body too large") {
		return apperrors.New(apperrors.ErrValidation, "request body is too large")
	}

	if errors.Is(err, multipart.ErrMessageTooLarge) {
		return apperrors.New(apperrors.ErrValidation, "request body is too large")
	}

	if errors.Is(err, http.ErrMissingFile) {
		return apperrors.New(apperrors.ErrValidation, "file is required")
	}

	return apperrors.New(apperrors.ErrValidation, err.Error())
}
