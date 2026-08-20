package router

import (
	"go-income-expense-tracker-app/internal/controller"

	"github.com/labstack/echo/v5"
)

func RegisterRecordRoutes(protectedGroup *echo.Group, recordController *controller.RecordController) {
	protectedGroup.POST("/records", recordController.Create)
	protectedGroup.GET("/records", recordController.GetAll)
	protectedGroup.GET("/records/:id", recordController.GetByID)
	protectedGroup.PUT("/records/:id", recordController.Update)
	protectedGroup.DELETE("/records/:id", recordController.Delete)
	protectedGroup.GET("/records/report", recordController.ExportReport)
	protectedGroup.GET("/records/summary", recordController.GetSummary)
}
