package controller

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type RecordController struct {
	recordService service.RecordService
}

func NewRecordController(recordService service.RecordService) *RecordController {
	return &RecordController{recordService: recordService}
}

func (c *RecordController) Create(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	var req dto.RecordRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	record, err := c.recordService.Create(uint(userID), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRecord) {
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		if errors.Is(err, service.ErrCategoryNotFound) {
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to create record",
		})
	}

	return ctx.JSON(http.StatusCreated, dto.Response[*model.Record]{
		Status:  http.StatusCreated,
		Message: "record created successfully",
		Data:    record,
	})
}

func (c *RecordController) GetByID(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid record id",
		})
	}

	record, err := c.recordService.GetByID(uint(userID), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, service.ErrRecordForbidden):
			return ctx.JSON(http.StatusForbidden, dto.Response[any]{
				Status:  http.StatusForbidden,
				Message: err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
				Status:  http.StatusInternalServerError,
				Message: "failed to get record",
			})
		}
	}

	return ctx.JSON(http.StatusOK, dto.Response[*model.Record]{
		Status:  http.StatusOK,
		Message: "record retrieved successfully",
		Data:    record,
	})
}

func (c *RecordController) GetAll(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	records, err := c.recordService.GetAllByUser(uint(userID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to get records",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[[]model.Record]{
		Status:  http.StatusOK,
		Message: "records retrieved successfully",
		Data:    records,
	})
}

func (c *RecordController) Update(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid record id",
		})
	}

	var req dto.RecordRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	record, err := c.recordService.Update(uint(userID), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, service.ErrRecordForbidden):
			return ctx.JSON(http.StatusForbidden, dto.Response[any]{
				Status:  http.StatusForbidden,
				Message: err.Error(),
			})
		case errors.Is(err, service.ErrInvalidRecord):
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		case errors.Is(err, service.ErrCategoryNotFound):
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
				Status:  http.StatusInternalServerError,
				Message: "failed to update record",
			})
		}
	}

	return ctx.JSON(http.StatusOK, dto.Response[*model.Record]{
		Status:  http.StatusOK,
		Message: "record updated successfully",
		Data:    record,
	})
}

func (c *RecordController) Delete(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid record id",
		})
	}

	if err := c.recordService.Delete(uint(userID), id); err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, service.ErrRecordForbidden):
			return ctx.JSON(http.StatusForbidden, dto.Response[any]{
				Status:  http.StatusForbidden,
				Message: err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
				Status:  http.StatusInternalServerError,
				Message: "failed to delete record",
			})
		}
	}

	return ctx.JSON(http.StatusOK, dto.Response[any]{
		Status:  http.StatusOK,
		Message: "record deleted successfully",
	})
}

func (c *RecordController) ExportReport(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")

	if startDate == "" || endDate == "" {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "start_date and end_date are required",
		})
	}

	url, err := c.recordService.ExportReport(uint(userID), dto.ExportReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		if errors.Is(err, service.ErrInvalidRecord) {
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: "invalid date format, use DD-MM-YYYY",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to generate report",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[string]{
		Status:  http.StatusOK,
		Message: "report exported successfully",
		Data:    url,
	})
}

func (c *RecordController) GetSummary(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	month := ctx.QueryParam("month")
	if month == "" {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "month is required",
		})
	}

	summary, err := c.recordService.GetSummary(uint(userID), dto.SummaryRequest{Month: month})
	if err != nil {
		if errors.Is(err, service.ErrInvalidRecord) {
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: "invalid month format, use MM-YYYY",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to get summary",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[*dto.SummaryResponse]{
		Status:  http.StatusOK,
		Message: "summary retrieved successfully",
		Data:    summary,
	})
}
