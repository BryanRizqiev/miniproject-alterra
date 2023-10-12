package lib

import (
	"miniproject-alterra/app/config"

	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {

	return echoJwt.WithConfig(echoJwt.Config{
		SigningKey: []byte(config.JWT_SECRET),
	})

}
