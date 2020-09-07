package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"pb-dev-be/config"
	"pb-dev-be/helpers"
	"pb-dev-be/models"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

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

	res, err, user := models.CheckLogin(user)

	if !res {
		// c.Response().WriteHeader(http.StatusUnauthorized)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"messages": "Wrong Email Or Password",
		})
	}

	if err != nil {
		// c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"messages": err.Error(),
		})
	}

	expTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		UserId: user.Id,
		Name:   user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	//generate token
	// return c.String(http.StatusOK, "Success To Login")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret")) //uuid if u want to deploy it to server the sample was "secret"
	if err != nil {
		return err
	}

	result, err := models.UpdateLastLogin(user)

	if err != nil {
		// c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, result)
	}

	// http.SetCookie(c.Response().Writer, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   t,
	// 	Expires: expTime,
	// 	Domain: ,
	// })

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":  user.Id,
		"name":     user.Name,
		"role":     user.UserRole,
		"verified": user.IsVerified,
		"token":    t,
	})
}

func RefreshToken(c echo.Context) error {
	cookie, err := c.Request().Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// c.Response().WriteHeader(http.StatusUnauthorized)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"messages": err.Error(),
			})
		}
		// c.Response().WriteHeader(http.StatusBadRequest)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"messages": err.Error(),
		})
	}

	tknStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// c.Response().WriteHeader(http.StatusUnauthorized)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"messages": err.Error(),
			})
		}
		// c.Response().WriteHeader(http.StatusBadRequest)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"messages": err.Error(),
		})
	}
	if !tkn.Valid {
		// c.Response().WriteHeader(http.StatusUnauthorized)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"messages": err.Error(),
		})
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		// c.Response().WriteHeader(http.StatusUnauthorized)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"messages": "Unauthorized",
		})
	}

	expTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		// c.Response().WriteHeader(http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"messages": err.Error(),
		})
	}

	// http.SetCookie(c.Response().Writer, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   t,
	// 	Expires: expTime,
	// })

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": claims.Id,
		"name":    claims.Name,
		"token":   t,
	})
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

func GetAuthForMidtrans(c echo.Context) error {
	conf := config.GetConfig()
	SERVER_KEY := conf.MIDTRANS_SERVER_KEY + ":"
	auth := base64.StdEncoding.EncodeToString([]byte(SERVER_KEY))

	if auth == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"messages": "Error While Encoding Server Key " + SERVER_KEY})
	}

	return c.JSON(http.StatusOK, map[string]string{"key": auth})
}
