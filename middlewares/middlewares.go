package middlewares

import (
	"github.com/labstack/echo/v4/middleware"
)

var IsAuthtenticated = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("secret"),
})
