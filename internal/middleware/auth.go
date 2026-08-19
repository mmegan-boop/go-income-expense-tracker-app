package middleware

import (
	"context"
	"errors"
	"fmt"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

type JWTCustomClaims struct {
	ID   int        `json:"id"`
	Role model.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey       string
	ExpiresDuration int
}

type contextKey string

const userContextKey = contextKey("user")

func (jwtCfg *JWTConfig) Init() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey: []byte(jwtCfg.SecretKey),
	}
}

func (jwtCfg *JWTConfig) GenerateToken(userID int, role model.Role) (string, error) {
	expire := jwt.NewNumericDate(time.Now().Local().Add(time.Minute * time.Duration(int64(jwtCfg.ExpiresDuration))))

	claims := &JWTCustomClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expire,
		},
		Role: role,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := rawToken.SignedString([]byte(jwtCfg.SecretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

// Retrieves the authenticated user's JWT claims
func GetUser(ctx context.Context) (*JWTCustomClaims, error) {
	user, ok := ctx.Value(userContextKey).(*jwt.Token)
	fmt.Println("get user user", user, "isi userContextKey", userContextKey)
	if !ok || user == nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := user.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

// Retrieves the authenticated user's ID from the JWT
func GetUserID(ctx context.Context) (int, error) {
	claim, err := GetUser(ctx)

	if err != nil {
		return 0, errors.New("invalid token")
	}

	return claim.ID, nil
}

// RequireRole returns middleware that restricts access to users with one of the specified roles.
func RequireRole(allowedRoles ...model.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claims, err := GetUser(c.Request().Context())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.Response[string]{
					Status:  http.StatusUnauthorized,
					Message: "unauthorized",
				})
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, dto.Response[string]{
				Status:  http.StatusForbidden,
				Message: "forbidden: insufficient permissions",
			})
		}
	}
}

// Middleware that validates/accesses the authenticated JWT and makes it available to downstream handlers
func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Retrieve JWT token from the Echo context using the "user" key
		user := c.Get("user").(*jwt.Token)
		fmt.Println("isi user di verifyToken", user)
		if user == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "invalid token",
			})
		}

		// After successful authentication, downstream handlers can access:
		// - c.Request().Context() to retrieve the JWT through GetUser or GetUserID
		ctx := context.WithValue(c.Request().Context(), userContextKey, user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
