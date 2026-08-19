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

func setupAuthRepoTest(t *testing.T) (repository.AuthRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:             db,
		WithoutReturning: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	return repository.NewAuthRepository(gormDB), mock
}

func newAuthUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "username", "email", "password", "role",
		"created_at", "updated_at", "deleted_at",
	})
}

func TestAuthRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupAuthRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`INSERT INTO "users"`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(&model.User{
			Username: "john",
			Email:    "john@example.com",
			Password: "hashedpassword",
			Role:     model.RoleUser,
		})
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupAuthRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`INSERT INTO "users"`).
			WillReturnError(fmt.Errorf("duplicate key"))

		err := repo.Create(&model.User{
			Username: "john",
			Email:    "john@example.com",
			Password: "hashedpassword",
			Role:     model.RoleUser,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
	})
}

func TestAuthRepository_FindByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupAuthRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newAuthUserRows().AddRow(
			uint(1), "john", "john@example.com", "hashedpassword", "user",
			fixedTime, fixedTime, nil,
		)
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		user, err := repo.FindByEmail("john@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "john", user.Username)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock := setupAuthRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.FindByEmail("unknown@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
