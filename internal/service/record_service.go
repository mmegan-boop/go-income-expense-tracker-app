package service

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
	"time"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrRecordForbidden  = errors.New("record does not belong to this user")
	ErrInvalidRecord    = errors.New("invalid record type or amount")
)

type RecordService interface {
	Create(userID uint, req dto.RecordRequest) (*model.Record, error)
	GetByID(userID uint, id int) (*model.Record, error)
	GetAllByUser(userID uint) ([]model.Record, error)
	Update(userID uint, id int, req dto.RecordRequest) (*model.Record, error)
	Delete(userID uint, id int) error
}

type recordService struct {
	recordRepository repository.RecordRepository
}

func NewRecordService(recordRepository repository.RecordRepository) RecordService {
	return &recordService{recordRepository: recordRepository}
}

func (s *recordService) Create(userID uint, req dto.RecordRequest) (*model.Record, error) {
	if !isValidRecordType(req.RecordType) || req.Amount <= 0 {
		return nil, ErrInvalidRecord
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

func isValidRecordType(recordType string) bool {
	return model.RecordType(recordType) == model.RecordTypeIncome ||
		model.RecordType(recordType) == model.RecordTypeExpense
}

func parseRecordDate(value string) time.Time {
	if value == "" {
		return time.Now()
	}

	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if t, err := time.Parse(layout, value); err == nil {
			return t
		}
	}

	return time.Now()
}
