package service

import (
	"bytes"
	"errors"
	"fmt"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrRecordForbidden = errors.New("record does not belong to this user")
	ErrInvalidRecord   = errors.New("invalid record type or amount")
)

type RecordService interface {
	Create(userID uint, req dto.RecordRequest) (*model.Record, error)
	GetByID(userID uint, id int) (*model.Record, error)
	GetAllByUser(userID uint) ([]model.Record, error)
	Update(userID uint, id int, req dto.RecordRequest) (*model.Record, error)
	Delete(userID uint, id int) error
	ExportReport(userID uint, req dto.ExportReportRequest) (string, error)
	GetSummary(userID uint, req dto.SummaryRequest) (*dto.SummaryResponse, error)
}

type recordService struct {
	recordRepository   repository.RecordRepository
	categoryRepository repository.CategoryRepository
	cloudinaryService  CloudinaryService
}

func NewRecordService(recordRepository repository.RecordRepository, categoryRepository repository.CategoryRepository, cloudinaryService CloudinaryService) RecordService {
	return &recordService{recordRepository: recordRepository, categoryRepository: categoryRepository, cloudinaryService: cloudinaryService}
}

func (s *recordService) Create(userID uint, req dto.RecordRequest) (*model.Record, error) {
	if !isValidRecordType(req.RecordType) || req.Amount <= 0 {
		return nil, ErrInvalidRecord
	}

	if _, err := s.categoryRepository.FindByID(int(req.CategoryID)); err != nil {
		return nil, ErrCategoryNotFound
	}

	record := &model.Record{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		RecordType:  model.RecordType(req.RecordType),
		Amount:      req.Amount,
		Description: req.Description,
		RecordDate:  parseRecordDate(req.RecordDate),
	}

	if err := s.recordRepository.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *recordService) GetByID(userID uint, id int) (*model.Record, error) {
	record, err := s.recordRepository.FindByID(id)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	if record.UserID != userID {
		return nil, ErrRecordForbidden
	}

	return record, nil
}

func (s *recordService) GetAllByUser(userID uint) ([]model.Record, error) {
	return s.recordRepository.FindAllByUserID(userID)
}

func (s *recordService) Update(userID uint, id int, req dto.RecordRequest) (*model.Record, error) {
	record, err := s.recordRepository.FindByID(id)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	if record.UserID != userID {
		return nil, ErrRecordForbidden
	}

	if !isValidRecordType(req.RecordType) || req.Amount <= 0 {
		return nil, ErrInvalidRecord
	}

	if _, err := s.categoryRepository.FindByID(int(req.CategoryID)); err != nil {
		return nil, ErrCategoryNotFound
	}

	record.CategoryID = req.CategoryID
	record.RecordType = model.RecordType(req.RecordType)
	record.Amount = req.Amount
	record.Description = req.Description
	record.RecordDate = parseRecordDate(req.RecordDate)

	if err := s.recordRepository.Update(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *recordService) Delete(userID uint, id int) error {
	record, err := s.recordRepository.FindByID(id)
	if err != nil {
		return ErrRecordNotFound
	}

	if record.UserID != userID {
		return ErrRecordForbidden
	}

	return s.recordRepository.Delete(id)
}

func (s *recordService) ExportReport(userID uint, req dto.ExportReportRequest) (string, error) {
	startDate, err := time.Parse("02-01-2006", req.StartDate)
	if err != nil {
		return "", ErrInvalidRecord
	}

	endDate, err := time.Parse("02-01-2006 15:04:05", req.EndDate+" 23:59:59")
	if err != nil {
		return "", ErrInvalidRecord
	}

	records, err := s.recordRepository.FindAllByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		return "", ErrRecordNotFound
	}
	f := excelize.NewFile()
	defer f.Close()
	f.SetSheetName("Sheet1", "Records")
	sheet := "Records"
	// Set column widths
	f.SetColWidth(sheet, "A", "A", 20) // Date
	f.SetColWidth(sheet, "B", "B", 20) // Type
	f.SetColWidth(sheet, "C", "C", 20) // Category
	f.SetColWidth(sheet, "D", "D", 20) // Amount
	f.SetColWidth(sheet, "E", "E", 30) // Description
	f.SetCellValue(sheet, "A1", "Date")
	f.SetCellValue(sheet, "B1", "Type")
	f.SetCellValue(sheet, "C1", "Category")
	f.SetCellValue(sheet, "D1", "Amount")
	f.SetCellValue(sheet, "E1", "Description")

	for i, record := range records {
		row := i + 2

		categoryName := ""
		if category, err := s.categoryRepository.FindByID(int(record.CategoryID)); err == nil {
			categoryName = category.Name
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), record.RecordDate.Format("02-01-2006"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), string(record.RecordType))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), categoryName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), record.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), record.Description)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf(
		"%s_income-expense-report_%s_to_%s",
		uuid.New().String(),
		req.StartDate,
		req.EndDate,
	)

	url, err := s.cloudinaryService.UploadRaw(buf.Bytes(), fileName)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *recordService) GetSummary(userID uint, req dto.SummaryRequest) (*dto.SummaryResponse, error) {
	startDate, err := time.Parse("01-2006", req.Month)
	if err != nil {
		return nil, ErrInvalidRecord
	}

	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	records, err := s.recordRepository.FindAllByUserIDAndDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	var totalIncome float64
	var totalExpense float64
	incomeByCategory := make(map[uint]float64)
	expenseByCategory := make(map[uint]float64)

	// separates totals into totalIncome and totalExpense based on record_type
	for _, record := range records {
		switch record.RecordType {
		case model.RecordTypeIncome:
			totalIncome += record.Amount
			incomeByCategory[record.CategoryID] += record.Amount
		case model.RecordTypeExpense:
			totalExpense += record.Amount
			expenseByCategory[record.CategoryID] += record.Amount
		}
	}

	incomeCategories := buildCategorySummaries(incomeByCategory, totalIncome, s.categoryRepository)
	expenseCategories := buildCategorySummaries(expenseByCategory, totalExpense, s.categoryRepository)

	return &dto.SummaryResponse{
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		IncomeCategories:  incomeCategories,
		ExpenseCategories: expenseCategories,
	}, nil
}

func buildCategorySummaries(byCategory map[uint]float64, total float64, categoryRepo repository.CategoryRepository) []dto.CategorySummary {
	var summaries []dto.CategorySummary

	// iterates over each categoryID → amount pair
	for categoryID, amount := range byCategory {
		// queries the repository to get the name from the ID
		categoryName := ""
		if category, err := categoryRepo.FindByID(int(categoryID)); err == nil {
			categoryName = category.Name
		}

		// amount / total * 100, rounded to 2 decimal places
		percentage := 0.0
		if total > 0 {
			percentage = math.Round(amount/total*10000) / 100
		}

		// creates a CategorySummary struct with name, amount, and percentage for each category
		summaries = append(summaries, dto.CategorySummary{
			CategoryName: categoryName,
			Amount:       amount,
			Percentage:   percentage,
		})
	}

	return summaries
}

func isValidRecordType(recordType string) bool {
	return model.RecordType(recordType) == model.RecordTypeIncome ||
		model.RecordType(recordType) == model.RecordTypeExpense
}

// converts a date string into a time.Time value
func parseRecordDate(value string) time.Time {
	if value == "" {
		return time.Now()
	}

	for _, layout := range []string{"02-01-2006", time.RFC3339} {
		if t, err := time.Parse(layout, value); err == nil {
			return t
		}
	}

	return time.Now()
}
