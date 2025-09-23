package jwtutil

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func JWTMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")

	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
		ContextKey: "user",
	})
}

func ExtractUser(c echo.Context) (*CustomClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		fmt.Println("user is not *jwt.Token")
		return nil, fmt.Errorf("invalid token format")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}
	
	return claims, nil
}
