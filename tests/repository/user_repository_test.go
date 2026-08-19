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

func setupUserRepoTest(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:             db,
		WithoutReturning: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	return repository.NewUserRepository(gormDB), mock
}

func newUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "username", "email", "password", "role",
		"created_at", "updated_at", "deleted_at",
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newUserRows().AddRow(
			uint(1), "john", "john@example.com", "hashedpassword", "user",
			fixedTime, fixedTime, nil,
		)
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		user, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "john", user.Username)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, "hashedpassword", user.Password)
		assert.Equal(t, model.Role("user"), user.Role)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.FindByID(999)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	t.Run("success with multiple users", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newUserRows().
			AddRow(uint(1), "john", "john@example.com", "hashedpw1", "user", fixedTime, fixedTime, nil).
			AddRow(uint(2), "jane", "jane@example.com", "hashedpw2", "admin", fixedTime, fixedTime, nil)
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		users, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "john", users[0].Username)
		assert.Equal(t, "jane", users[1].Username)
		assert.Equal(t, model.RoleAdmin, users[1].Role)
	})

	t.Run("empty result", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		rows := newUserRows()
		mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

		users, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnError(fmt.Errorf("db error"))

		users, err := repo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "users"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(&model.User{
			ID:       1,
			Username: "johnny",
			Email:    "john@example.com",
			Password: "hashedpassword",
			Role:     model.RoleUser,
		})
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "users"`).
			WillReturnError(fmt.Errorf("constraint violation"))

		err := repo.Update(&model.User{
			ID:       1,
			Username: "johnny",
			Email:    "john@example.com",
			Role:     model.RoleUser,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "constraint violation")
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "users" SET "deleted_at"`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		repo, mock := setupUserRepoTest(t)
		defer assert.NoError(t, mock.ExpectationsWereMet())

		mock.ExpectExec(`UPDATE "users" SET "deleted_at"`).
			WillReturnError(fmt.Errorf("delete failed"))

		err := repo.Delete(999)
		assert.Error(t, err)
	})
}
