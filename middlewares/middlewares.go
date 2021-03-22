package middlewares

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var IsAuthtenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("secret"),
})

func ValidateKey(username, password string, c echo.Context) (bool, error) {
	// Be careful to use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(username), []byte("PB-xPvRXeUjSI")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("vrNWDnTCwlfKjaKUeywt")) == 1 {
		return true, nil
	}
	return false, nil
}
