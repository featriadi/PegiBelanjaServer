package controller

import (
	"fmt"
	"net/http"
	"pb-dev-be/helpers"
	"pb-dev-be/models"

	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
)

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	// userRole := claims["user_role"].(string)
	// verified := claims["verified"].(string)
	// rememberMe := claims["remember_me"].(string)

	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func IsMitra(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isMitra := claims["user_role"].(string)
		if isMitra != "MIT" {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isMitra := claims["user_role"].(string)
		if isMitra != "ADM" {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

func CheckLogin(c echo.Context) error {
	var user models.User
	// var ares models.AuthResponse

	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	user.RememberMe = c.FormValue("remember_me")

	res, err, user := models.CheckLogin(user)

	if !res {
		return echo.ErrUnauthorized
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"messages": err.Error(),
		})
	}

	//generate token
	// return c.String(http.StatusOK, "Success To Login")

	token := jwt.New(jwt.SigningMethodHS256)

	//Set Claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["user_role"] = user.UserRole
	claims["verified"] = user.IsVerified
	claims["remember_me"] = user.RememberMe
	claims["expired"] = time.Now().Add(time.Hour * 168).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret")) //uuid if u want to deploy it to server the sample was "secret"
	if err != nil {
		return err
	}

	_, err = models.UpdateLastLogin(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": t,
	})

	// refreshToken := jwt.New(jwt.SigningMethodHS256)
	// rtClaims := refreshToken.Claims.(jwt.MapClaims)
	// claims["user_id"] = user.Id
	// claims["name"] = user.Name
	// claims["email"] = user.Email
	// claims["user_role"] = user.UserRole
	// claims["verified"] = user.IsVerified
	// claims["remember_me"] = user.RememberMe
	// rtClaims["expired"] = time.Now().Add(time.Hour * 24).Unix()
	// rt, err := refreshToken.SignedString([]byte("secret"))
	// if err != nil {
	// 	return err
	// }
	// return c.JSON(http.StatusOK, map[string]string{
	// 	"access_token":  t,
	// 	"refresh_token": rt,
	// })

}

func GenerateHashPassword(c echo.Context) error {
	password := c.Param("password")

	hash, _ := helpers.HashPassword(password)

	return c.JSON(http.StatusOK, hash)
}

func generateTokenPairMitra() (map[string]string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["user_role"] = "MIT"
	claims["expired"] = time.Now().Add(time.Minute * 15).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["expired"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}

func Token(c echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	err := c.Bind(&tokenReq)

	if err != nil {
		return err
	}

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {

			newTokenPair, err := generateTokenPairMitra()
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, newTokenPair)
		}

		return echo.ErrUnauthorized
	}

	return err
}
