package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/oauth"
	"github.com/Dhruv9449/mou/pkg/serializers"
	"github.com/Dhruv9449/mou/pkg/utils"
	"github.com/Dhruv9449/mou/pkg/utils/dbutils"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

func handleGoogleLogin(c *fiber.Ctx) error {
	URL, err := url.Parse(oauth.OauthConf.Endpoint.AuthURL)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not parse the URL",
		})
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauth.OauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauth.OauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauth.OauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauth.OauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()

	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func callback(c *fiber.Ctx) error {
	state := c.FormValue("state")

	if state != oauth.OauthStateString {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid oauth state",
		})
	}

	code := c.FormValue("code")
	token, err := oauth.OauthConf.Exchange(context.Background(), code)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Code exchange failed",
		})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed getting user info",
		})
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed reading response body",
		})
	}

	var user models.User
	json.Unmarshal(response, &user)

	// Check if user exists in database
	if !dbutils.CheckUserExists(user.Email) {
		database.DB.Create(&user)
	} else {
		database.DB.Where("email = ?", user.Email).First(&user)
	}

	tokenString, err := utils.CreateJWTToken(user.ID, user.Email, user.Role)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create JWT token",
		})
	}

	return c.Status(http.StatusOK).JSON(serializers.UserLoginSerializer(user, tokenString))
}

func googleLogin(c *fiber.Ctx) error {
	type RequestBody struct {
		Id_token string `json:"id_token"`
	}

	var body RequestBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify the token
	token, err := idtoken.Validate(context.Background(), body.Id_token, oauth.OauthConf.ClientID)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token",
			"message": err.Error(),
		})
	}

	var user models.User
	user.Email = token.Claims["email"].(string)
	user.Name = token.Claims["name"].(string)
	user.Picture = token.Claims["picture"].(string)

	// Check if user exists in database
	if !dbutils.CheckUserExists(user.Email) {
		database.DB.Create(&user)
	} else {
		database.DB.Where("email = ?", user.Email).First(&user)
		if user.Picture != token.Claims["picture"].(string) {
			database.DB.Model(&user).Update("picture", token.Claims["picture"].(string))
		}
	}

	tokenString, err := utils.CreateJWTToken(user.ID, user.Email, user.Role)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create JWT token",
		})
	}

	return c.Status(http.StatusOK).JSON(serializers.UserLoginSerializer(user, tokenString))
}

func AuthRouter(app *fiber.App) {
	group := app.Group("/oauth")
	group.Get("/google", handleGoogleLogin)
	group.Get("/callback", callback)
	group.Post("/login", googleLogin)
}
