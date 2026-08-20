package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
)

type MockRecordRepository struct {
	mock.Mock
}

func (m *MockRecordRepository) Create(record *model.Record) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockRecordRepository) FindByID(id int) (*model.Record, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Record), args.Error(1)
}

func (m *MockRecordRepository) FindAllByUserID(userID uint) ([]model.Record, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Record), args.Error(1)
}

func (m *MockRecordRepository) Update(record *model.Record) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockRecordRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRecordRepository) FindAllByUserIDAndDateRange(userID uint, startDate time.Time, endDate time.Time) ([]model.Record, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).([]model.Record), args.Error(1)
}

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(id int) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByName(name string) (*model.Category, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindAll() ([]model.Category, error) {
	args := m.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

var testCategory = &model.Category{ID: 1, Name: "Food"}

type MockCloudinaryService struct {
	mock.Mock
}

func (m *MockCloudinaryService) UploadRaw(data []byte, fileName string) (string, error) {
	args := m.Called(data, fileName)
	return args.String(0), args.Error(1)
}

func newRecordService() (*MockRecordRepository, *MockCategoryRepository, service.RecordService) {
	mockRecordRepo := new(MockRecordRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockCloudinary := new(MockCloudinaryService)
	svc := service.NewRecordService(mockRecordRepo, mockCategoryRepo, mockCloudinary)
	return mockRecordRepo, mockCategoryRepo, svc
}

func TestRecordService_Create(t *testing.T) {
	baseReq := dto.RecordRequest{
		CategoryID:  1,
		RecordType:  "income",
		Amount:      100.0,
		Description: "salary",
		RecordDate:  "15-01-2026",
	}

	t.Run("success income", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		mockCategoryRepo.On("FindByID", 1).Return(testCategory, nil)
		mockRecordRepo.On("Create", mock.AnythingOfType("*model.Record")).Return(nil)

		record, err := svc.Create(1, baseReq)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, uint(1), record.UserID)
		assert.Equal(t, model.RecordTypeIncome, record.RecordType)
		assert.Equal(t, 100.0, record.Amount)
		mockRecordRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("success expense", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		req := baseReq
		req.RecordType = "expense"

		mockCategoryRepo.On("FindByID", 1).Return(testCategory, nil)
		mockRecordRepo.On("Create", mock.AnythingOfType("*model.Record")).Return(nil)

		record, err := svc.Create(1, req)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, model.RecordTypeExpense, record.RecordType)
	})

	t.Run("invalid record type", func(t *testing.T) {
		_, _, svc := newRecordService()

		req := baseReq
		req.RecordType = "invalid"

		record, err := svc.Create(1, req)
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Nil(t, record)
	})

	t.Run("invalid amount zero", func(t *testing.T) {
		_, _, svc := newRecordService()

		req := baseReq
		req.Amount = 0

		record, err := svc.Create(1, req)
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Nil(t, record)
	})

	t.Run("invalid amount negative", func(t *testing.T) {
		_, _, svc := newRecordService()

		req := baseReq
		req.Amount = -50

		record, err := svc.Create(1, req)
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Nil(t, record)
	})

	t.Run("category not found", func(t *testing.T) {
		_, mockCategoryRepo, svc := newRecordService()

		mockCategoryRepo.On("FindByID", 1).Return(nil, fmt.Errorf("not found"))

		record, err := svc.Create(1, baseReq)
		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Nil(t, record)
	})

	t.Run("repository create error", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		mockCategoryRepo.On("FindByID", 1).Return(testCategory, nil)
		mockRecordRepo.On("Create", mock.AnythingOfType("*model.Record")).
			Return(fmt.Errorf("db error"))

		record, err := svc.Create(1, baseReq)
		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestRecordService_GetByID(t *testing.T) {
	ownerRecord := &model.Record{
		ID:     1,
		UserID: 1,
		Amount: 100.0,
	}

	t.Run("success", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindByID", 1).Return(ownerRecord, nil)

		record, err := svc.GetByID(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, uint(1), record.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		record, err := svc.GetByID(1, 999)
		assert.ErrorIs(t, err, service.ErrRecordNotFound)
		assert.Nil(t, record)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		otherUserRecord := &model.Record{ID: 1, UserID: 2}
		mockRecordRepo.On("FindByID", 1).Return(otherUserRecord, nil)

		record, err := svc.GetByID(1, 1)
		assert.ErrorIs(t, err, service.ErrRecordForbidden)
		assert.Nil(t, record)
	})
}

func TestRecordService_GetAllByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		records := []model.Record{
			{ID: 1, UserID: 1, Amount: 100.0},
			{ID: 2, UserID: 1, Amount: 200.0},
		}
		mockRecordRepo.On("FindAllByUserID", uint(1)).Return(records, nil)

		result, err := svc.GetAllByUser(1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("empty result", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindAllByUserID", uint(999)).Return([]model.Record{}, nil)

		result, err := svc.GetAllByUser(999)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestRecordService_Update(t *testing.T) {
	ownerRecord := &model.Record{
		ID:         1,
		UserID:     1,
		CategoryID: 1,
		RecordType: model.RecordTypeIncome,
		Amount:     100.0,
	}

	validReq := dto.RecordRequest{
		CategoryID:  1,
		RecordType:  "expense",
		Amount:      250.0,
		Description: "updated",
		RecordDate:  "15-01-2026",
	}

	t.Run("success", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		mockRecordRepo.On("FindByID", 1).Return(ownerRecord, nil)
		mockCategoryRepo.On("FindByID", 1).Return(testCategory, nil)
		mockRecordRepo.On("Update", mock.AnythingOfType("*model.Record")).Return(nil)

		record, err := svc.Update(1, 1, validReq)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, model.RecordTypeExpense, record.RecordType)
		assert.Equal(t, 250.0, record.Amount)
	})

	t.Run("not found", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		record, err := svc.Update(1, 999, validReq)
		assert.ErrorIs(t, err, service.ErrRecordNotFound)
		assert.Nil(t, record)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		otherRecord := &model.Record{ID: 1, UserID: 2}
		mockRecordRepo.On("FindByID", 1).Return(otherRecord, nil)

		record, err := svc.Update(1, 1, validReq)
		assert.ErrorIs(t, err, service.ErrRecordForbidden)
		assert.Nil(t, record)
	})

	t.Run("invalid record type", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindByID", 1).Return(ownerRecord, nil)

		req := validReq
		req.RecordType = "bad"
		record, err := svc.Update(1, 1, req)
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Nil(t, record)
	})

	t.Run("invalid amount", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindByID", 1).Return(ownerRecord, nil)

		req := validReq
		req.Amount = -10
		record, err := svc.Update(1, 1, req)
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Nil(t, record)
	})

	t.Run("category not found", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		mockRecordRepo.On("FindByID", 1).Return(ownerRecord, nil)
		mockCategoryRepo.On("FindByID", 1).Return(nil, fmt.Errorf("not found"))

		record, err := svc.Update(1, 1, validReq)
		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Nil(t, record)
	})

	t.Run("repository update error", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		mockRecordRepo.On("FindByID", 1).Return(ownerRecord, nil)
		mockCategoryRepo.On("FindByID", 1).Return(testCategory, nil)
		mockRecordRepo.On("Update", mock.AnythingOfType("*model.Record")).
			Return(fmt.Errorf("update failed"))

		record, err := svc.Update(1, 1, validReq)
		assert.Error(t, err)
		assert.Nil(t, record)
	})
}

func TestRecordService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		record := &model.Record{ID: 1, UserID: 1}
		mockRecordRepo.On("FindByID", 1).Return(record, nil)
		mockRecordRepo.On("Delete", 1).Return(nil)

		err := svc.Delete(1, 1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		err := svc.Delete(1, 999)
		assert.ErrorIs(t, err, service.ErrRecordNotFound)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		record := &model.Record{ID: 1, UserID: 2}
		mockRecordRepo.On("FindByID", 1).Return(record, nil)

		err := svc.Delete(1, 1)
		assert.ErrorIs(t, err, service.ErrRecordForbidden)
	})

	t.Run("repository delete error", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		record := &model.Record{ID: 1, UserID: 1}
		mockRecordRepo.On("FindByID", 1).Return(record, nil)
		mockRecordRepo.On("Delete", 1).Return(fmt.Errorf("delete failed"))

		err := svc.Delete(1, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}

func TestRecordService_GetSummary(t *testing.T) {
	jan2026 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	janEnd := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("success with records", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, svc := newRecordService()

		records := []model.Record{
			{ID: 1, UserID: 1, CategoryID: 1, RecordType: model.RecordTypeIncome, Amount: 5000},
			{ID: 2, UserID: 1, CategoryID: 2, RecordType: model.RecordTypeExpense, Amount: 1000},
			{ID: 3, UserID: 1, CategoryID: 3, RecordType: model.RecordTypeExpense, Amount: 500},
		}
		mockRecordRepo.On("FindAllByUserIDAndDateRange", uint(1), jan2026, janEnd).Return(records, nil)

		mockCategoryRepo.On("FindByID", 1).Return(&model.Category{ID: 1, Name: "Salary"}, nil)
		mockCategoryRepo.On("FindByID", 2).Return(&model.Category{ID: 2, Name: "Rent"}, nil)
		mockCategoryRepo.On("FindByID", 3).Return(&model.Category{ID: 3, Name: "Food"}, nil)

		summary, err := svc.GetSummary(1, dto.SummaryRequest{Month: "01-2026"})
		assert.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, 5000.0, summary.TotalIncome)
		assert.Equal(t, 1500.0, summary.TotalExpense)
		assert.Len(t, summary.IncomeCategories, 1)
		assert.Len(t, summary.ExpenseCategories, 2)
		assert.Equal(t, "Salary", summary.IncomeCategories[0].CategoryName)
	})

	t.Run("no records", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindAllByUserIDAndDateRange", uint(1), jan2026, janEnd).
			Return([]model.Record{}, nil)

		summary, err := svc.GetSummary(1, dto.SummaryRequest{Month: "01-2026"})
		assert.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, 0.0, summary.TotalIncome)
		assert.Equal(t, 0.0, summary.TotalExpense)
	})

	t.Run("invalid month format", func(t *testing.T) {
		_, _, svc := newRecordService()

		summary, err := svc.GetSummary(1, dto.SummaryRequest{Month: "invalid"})
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Nil(t, summary)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindAllByUserIDAndDateRange", uint(1), jan2026, janEnd).
			Return([]model.Record{}, fmt.Errorf("db error"))

		summary, err := svc.GetSummary(1, dto.SummaryRequest{Month: "01-2026"})
		assert.ErrorIs(t, err, service.ErrRecordNotFound)
		assert.Nil(t, summary)
	})
}

func TestRecordService_ExportReport(t *testing.T) {
	jan2026 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	janEnd := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("success with records", func(t *testing.T) {
		mockRecordRepo, mockCategoryRepo, _ := newRecordService()

		records := []model.Record{
			{ID: 1, UserID: 1, CategoryID: 1, RecordType: model.RecordTypeIncome, Amount: 5000, Description: "salary", RecordDate: jan2026},
		}
		mockRecordRepo.On("FindAllByUserIDAndDateRange", uint(1), jan2026, janEnd).Return(records, nil)
		mockCategoryRepo.On("FindByID", 1).Return(&model.Category{ID: 1, Name: "Salary"}, nil)

		fakeURL := "https://res.cloudinary.com/demo/raw/upload/v1/report.xlsx"
		mockCloudinary := new(MockCloudinaryService)
		mockCloudinary.On("UploadRaw", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string")).Return(fakeURL, nil)

		svcWithCloudinary := service.NewRecordService(mockRecordRepo, mockCategoryRepo, mockCloudinary)
		url, err := svcWithCloudinary.ExportReport(1, dto.ExportReportRequest{StartDate: "01-01-2026", EndDate: "31-01-2026"})
		assert.NoError(t, err)
		assert.Equal(t, fakeURL, url)
	})

	t.Run("invalid start date", func(t *testing.T) {
		_, _, svc := newRecordService()

		data, err := svc.ExportReport(1, dto.ExportReportRequest{StartDate: "bad-date", EndDate: "31-01-2026"})
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Empty(t, data)
	})

	t.Run("invalid end date", func(t *testing.T) {
		_, _, svc := newRecordService()

		data, err := svc.ExportReport(1, dto.ExportReportRequest{StartDate: "01-01-2026", EndDate: "bad-date"})
		assert.ErrorIs(t, err, service.ErrInvalidRecord)
		assert.Empty(t, data)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRecordRepo, _, svc := newRecordService()

		mockRecordRepo.On("FindAllByUserIDAndDateRange", uint(1), jan2026, janEnd).
			Return([]model.Record{}, fmt.Errorf("db error"))

		data, err := svc.ExportReport(1, dto.ExportReportRequest{StartDate: "01-01-2026", EndDate: "31-01-2026"})
		assert.ErrorIs(t, err, service.ErrRecordNotFound)
		assert.Empty(t, data)
	})
}
