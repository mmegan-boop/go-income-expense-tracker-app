package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"

	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"
	appvalidator "go-income-expense-tracker-app/internal/validator"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req dto.RegisterRequest) (*model.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) Login(req dto.LoginRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(id int) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserService) UpdateProfile(id int, req dto.UpdateUserRequest) (*model.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) Create(req dto.CategoryRequest) (*model.Category, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryService) GetByID(id int) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryService) GetAll() ([]model.Category, error) {
	args := m.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryService) Update(id int, req dto.CategoryRequest) (*model.Category, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockRecordService struct {
	mock.Mock
}

func (m *MockRecordService) Create(userID uint, req dto.RecordRequest) (*model.Record, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Record), args.Error(1)
}

func (m *MockRecordService) GetByID(userID uint, id int) (*model.Record, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Record), args.Error(1)
}

func (m *MockRecordService) GetAllByUser(userID uint) ([]model.Record, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Record), args.Error(1)
}

func (m *MockRecordService) Update(userID uint, id int, req dto.RecordRequest) (*model.Record, error) {
	args := m.Called(userID, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Record), args.Error(1)
}

func (m *MockRecordService) Delete(userID uint, id int) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

func (m *MockRecordService) ExportReport(userID uint, req dto.ExportReportRequest) (string, error) {
	args := m.Called(userID, req)
	return args.String(0), args.Error(1)
}

func (m *MockRecordService) GetSummary(userID uint, req dto.SummaryRequest) (*dto.SummaryResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SummaryResponse), args.Error(1)
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Validator = appvalidator.New()
	return e
}

func createRequest(method, path string, body any) *http.Request {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return req
}

func createRequestWithQuery(method, path string, params map[string]string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func createContextWithJWT(e *echo.Echo, req *http.Request, userID int, role model.Role) (*echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	token := &jwt.Token{
		Claims: &middleware.JWTCustomClaims{
			ID:   userID,
			Role: role,
		},
	}

	ctx.Set("user", token)

	middleware.VerifyToken(func(c *echo.Context) error {
		return nil
	})(ctx)

	return ctx, rec
}

func createContextWithPathParamsAndJWT(e *echo.Echo, req *http.Request, pathValues echo.PathValues, userID int, role model.Role) (*echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	routeInfo := &echo.RouteInfo{
		Method: req.Method,
		Path:   req.URL.Path,
	}
	for _, pv := range pathValues {
		routeInfo.Parameters = append(routeInfo.Parameters, pv.Name)
	}
	ctx.InitializeRoute(routeInfo, &pathValues)

	token := &jwt.Token{
		Claims: &middleware.JWTCustomClaims{
			ID:   userID,
			Role: role,
		},
	}

	ctx.Set("user", token)

	middleware.VerifyToken(func(c *echo.Context) error {
		return nil
	})(ctx)

	return ctx, rec
}
