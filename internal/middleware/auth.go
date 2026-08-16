package middleware

import (
	"context"
	"errors"
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

func GetUser(ctx context.Context) (*JWTCustomClaims, error) {
	user, ok := ctx.Value(userContextKey).(*jwt.Token)
	if !ok || user == nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := user.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func GetUserID(ctx context.Context) (int, error) {
	claim, err := GetUser(ctx)

	if err != nil {
		return 0, errors.New("invalid token")
	}

	return claim.ID, nil
}

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		user := c.Get("user").(*jwt.Token)

		if user == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "invalid token",
			})
		}

		ctx := context.WithValue(c.Request().Context(), userContextKey, user)
		c.SetRequest(c.Request().WithContext(ctx))

		userData, err := GetUser(ctx)
		if userData == nil || err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "invalid token",
			})
		}

		c.Set("userData", userData)

		return next(c)
	}
}
