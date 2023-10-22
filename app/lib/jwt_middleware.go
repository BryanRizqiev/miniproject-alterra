package lib

import (
	"miniproject-alterra/app/config"
	global_response "miniproject-alterra/module/global/controller/response"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {

	return echoJwt.WithConfig(echoJwt.Config{
		SigningKey: []byte(config.JWT_SECRET),
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusUnauthorized, global_response.StandartResponse{
				Message: "Unauthorized.",
			})
		},
	})

}

func CreateToken(userId string, email string) (string, error) {

	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["email"] = email
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWT_SECRET))

}

func ExtractToken(e echo.Context) (string, string) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["user_id"].(string)
		email := claims["email"].(string)
		return userId, email
	}
	return "", ""
}
