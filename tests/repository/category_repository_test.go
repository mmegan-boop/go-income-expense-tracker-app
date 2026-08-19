package repository_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
)

func setupCategoryRepoTest(t *testing.T) (repository.CategoryRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:             db,
		WithoutReturning: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	return repository.NewCategoryRepository(gormDB), mock
}

func newCategoryRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "name", "created_at", "updated_at", "deleted_at",
	})
}

func TestCategoryRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`INSERT INTO "categories"`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&model.Category{Name: "Food"})
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`INSERT INTO "categories"`).
			WillReturnError(fmt.Errorf("connection lost"))

		err := repo.Create(&model.Category{Name: "Food"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection lost")
	})
}

func TestCategoryRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newCategoryRows().AddRow(
			uint(1), "Food", fixedTime, fixedTime, nil,
		)
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnRows(rows)

		cat, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, uint(1), cat.ID)
		assert.Equal(t, "Food", cat.Name)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "categories"`).
			WillReturnError(sql.ErrNoRows)

		cat, err := repo.FindByID(999)
		assert.Error(t, err)
		assert.Nil(t, cat)
	})
}

func TestCategoryRepository_FindByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newCategoryRows().AddRow(
			uint(1), "Food", fixedTime, fixedTime, nil,
		)
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnRows(rows)

		cat, err := repo.FindByName("Food")
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, "Food", cat.Name)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "categories"`).
			WillReturnError(sql.ErrNoRows)

		cat, err := repo.FindByName("Nonexistent")
		assert.Error(t, err)
		assert.Nil(t, cat)
	})
}

func TestCategoryRepository_FindAll(t *testing.T) {
	t.Run("success with multiple categories", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newCategoryRows().
			AddRow(uint(1), "Food", fixedTime, fixedTime, nil).
			AddRow(uint(2), "Transport", fixedTime, fixedTime, nil).
			AddRow(uint(3), "Salary", fixedTime, fixedTime, nil)
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnRows(rows)

		categories, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Len(t, categories, 3)
		assert.Equal(t, "Food", categories[0].Name)
		assert.Equal(t, "Transport", categories[1].Name)
		assert.Equal(t, "Salary", categories[2].Name)
	})

	t.Run("empty result", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newCategoryRows()
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnRows(rows)

		categories, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "categories"`).
			WillReturnError(fmt.Errorf("db error"))

		categories, err := repo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, categories)
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "categories"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&model.Category{
			ID:   1,
			Name: "Groceries",
		})
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "categories"`).
			WillReturnError(fmt.Errorf("constraint violation"))

		err := repo.Update(&model.Category{
			ID:   1,
			Name: "Groceries",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "constraint violation")
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "categories" SET "deleted_at"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupCategoryRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "categories" SET "deleted_at"`).
			WillReturnError(fmt.Errorf("delete failed"))

		err := repo.Delete(999)
		assert.Error(t, err)
	})
}
