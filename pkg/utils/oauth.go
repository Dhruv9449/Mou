package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateJWTToken(id uint, email string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["email"] = email
	claims["role"] = role

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ValidateTokenAndGetUser(tokenString string) (models.User, error) {
	var err error
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return models.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		user := models.User{}
		if database.DB.Where("id = ?", claims["id"]).First(&user).RowsAffected == 0 {
			return models.User{}, errors.New("user not found")
		}
		database.DB.Where("id = ?", claims["id"]).First(&user)
		return user, nil
	}

	return models.User{}, err
}

func GetUserFromHeader(c *fiber.Ctx) (models.User, error) {
	headers := c.GetReqHeaders()
	authorizationString := strings.Split(headers["Authorization"], " ")
	if len(authorizationString) != 2 {
		return models.User{}, errors.New("invalid Authorization Header")
	}

	token := authorizationString[1]

	if token == "" {
		return models.User{}, errors.New("not logged in")
	}
	user, err := ValidateTokenAndGetUser(token)

	if err != nil {
		return models.User{}, errors.New("invalid credentials")
	}

	return user, nil
}
