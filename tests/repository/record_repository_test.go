package repository_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
)

func setupRecordRepoTest(t *testing.T) (repository.RecordRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:             db,
		WithoutReturning: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	return repository.NewRecordRepository(gormDB), mock
}

func newRecordRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "user_id", "category_id", "record_type", "amount",
		"description", "record_date", "created_at", "updated_at", "deleted_at",
	})
}

var fixedTime = time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

func TestRecordRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`INSERT INTO "records"`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&model.Record{
			UserID:      1,
			CategoryID:  1,
			RecordType:  model.RecordTypeIncome,
			Amount:      100.0,
			Description: "salary",
			RecordDate:  fixedTime,
		})
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`INSERT INTO "records"`).
			WillReturnError(fmt.Errorf("connection lost"))

		err := repo.Create(&model.Record{
			UserID:      1,
			CategoryID:  1,
			RecordType:  model.RecordTypeIncome,
			Amount:      100.0,
			Description: "salary",
			RecordDate:  fixedTime,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection lost")
	})
}

func TestRecordRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newRecordRows().AddRow(
			uint(1), uint(1), uint(1), "income", 100.0,
			"salary", fixedTime, fixedTime, fixedTime, nil,
		)
		mock.ExpectQuery(`SELECT \* FROM "records"`).WillReturnRows(rows)

		record, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, uint(1), record.ID)
		assert.Equal(t, uint(1), record.UserID)
		assert.Equal(t, model.RecordTypeIncome, record.RecordType)
		assert.Equal(t, 100.0, record.Amount)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "records"`).
			WillReturnError(sql.ErrNoRows)

		record, err := repo.FindByID(999)
		assert.Error(t, err)
		assert.Nil(t, record)
	})
}

func TestRecordRepository_FindAllByUserID(t *testing.T) {
	t.Run("success with multiple records", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newRecordRows().
			AddRow(uint(1), uint(1), uint(1), "income", 100.0, "salary", fixedTime, fixedTime, fixedTime, nil).
			AddRow(uint(2), uint(1), uint(2), "expense", 50.0, "food", fixedTime, fixedTime, fixedTime, nil)
		mock.ExpectQuery(`SELECT \* FROM "records"`).WillReturnRows(rows)

		records, err := repo.FindAllByUserID(1)
		assert.NoError(t, err)
		assert.Len(t, records, 2)
		assert.Equal(t, model.RecordTypeIncome, records[0].RecordType)
		assert.Equal(t, model.RecordTypeExpense, records[1].RecordType)
	})

	t.Run("empty result", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newRecordRows()
		mock.ExpectQuery(`SELECT \* FROM "records"`).WillReturnRows(rows)

		records, err := repo.FindAllByUserID(999)
		assert.NoError(t, err)
		assert.Empty(t, records)
	})
}

func TestRecordRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "records"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&model.Record{
			ID:          1,
			UserID:      1,
			CategoryID:  2,
			RecordType:  model.RecordTypeExpense,
			Amount:      200.0,
			Description: "updated grocery",
			RecordDate:  fixedTime,
		})
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "records"`).
			WillReturnError(fmt.Errorf("constraint violation"))

		err := repo.Update(&model.Record{
			ID:         1,
			UserID:     1,
			CategoryID: 2,
			RecordType: model.RecordTypeExpense,
			Amount:     200.0,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "constraint violation")
	})
}

func TestRecordRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "records" SET "deleted_at"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "records" SET "deleted_at"`).
			WillReturnError(fmt.Errorf("delete failed"))

		err := repo.Delete(999)
		assert.Error(t, err)
	})
}

func TestRecordRepository_FindAllByUserIDAndDateRange(t *testing.T) {
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("success with records", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newRecordRows().
			AddRow(uint(1), uint(1), uint(1), "income", 5000.0, "salary", fixedTime, fixedTime, fixedTime, nil).
			AddRow(uint(2), uint(1), uint(2), "expense", 150.0, "groceries", fixedTime, fixedTime, fixedTime, nil)
		mock.ExpectQuery(`SELECT \* FROM "records"`).WillReturnRows(rows)

		records, err := repo.FindAllByUserIDAndDateRange(1, startDate, endDate)
		assert.NoError(t, err)
		assert.Len(t, records, 2)
	})

	t.Run("empty range", func(t *testing.T) {
		repo, mock := setupRecordRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newRecordRows()
		mock.ExpectQuery(`SELECT \* FROM "records"`).WillReturnRows(rows)

		records, err := repo.FindAllByUserIDAndDateRange(1, startDate, endDate)
		assert.NoError(t, err)
		assert.Empty(t, records)
	})
}
