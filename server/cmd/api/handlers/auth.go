package handlers

import (
	"server/internal/database"
	"server/internal/models"
	"server/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	keyPassword        = "password"
	keyConfirmPassword = "confirmPassword" // Used in Register and UpdatePassword
	keyFirstName       = "firstName"       // Used in Register and UpdateInfo
	keyLastName        = "lastName"        // Used in Register and UpdateInfo
	keyEmail           = "email"           // Used in Register, Login, UpdateInfo, RequestResetPassword
	keyRefreshToken    = "refreshToken"    // Used in RefreshToken
)

// Register handles new user registration.
// It parses user details (firstName, lastName, email, password, confirmPassword) from the request body.
// It validates that the passwords match, creates a new user record in the database,
// generates JWT access and refresh tokens, sets them as HTTP-only cookies,
// and returns the newly created user object (excluding password).
func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": err,
			})
	}

	if data[keyPassword] != data[keyConfirmPassword] {
		c.Status(400)
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "passwords do not match",
			})
	}

	user := models.User{
		FirstName: data[keyFirstName],
		LastName:  data[keyLastName],
		Email:     data[keyEmail],
	}

	user.SetPassword(data[keyPassword])

	res := database.DB.Create(&user)
	// verify if user was created
	if res.Error != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": res.Error.Error(),
			})
	}

	err := util.GenerateUserTokens(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"internal server error": err,
			})
	}

	// set cookie with jwt
	setCookies(c, user)

	return c.Status(fiber.StatusOK).JSON(user)
}

// Login handles user authentication.
// It parses the email and password from the request body.
// It finds the user by email, verifies the password, generates new JWT access and refresh tokens,
// sets them as HTTP-only cookies, and returns a success message along with the authenticated user object.
// Returns 404 if the user is not found or 400 if the password is incorrect.
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data[keyEmail]).First(&user)

	if user.ID == "" {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"user": []string{"not found"},
				},
			})
	}

	if err := user.ComparePassword(data[keyPassword]); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"password": []string{"incorrect"},
				},
			})
	}

	err := util.GenerateUserTokens(&user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// set jwt and refreshjwt cookies
	setCookies(c, user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

// RefreshToken handles the renewal of JWT access tokens using a refresh token.
// It expects the refresh token in the request body.
// It parses and validates the refresh token, finds the associated user,
// generates new JWT access and refresh tokens, sets them as HTTP-only cookies,
// and returns a success message along with the user object containing the new tokens.
// Returns 400 if the token is invalid/expired or 404 if the user is not found.
func RefreshToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		// Consider using a proper logger instead of fmt.Println or just returning
		return err
	}
	refreshToken := data[keyRefreshToken]

	issuer, err := util.ParseJwt(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "token invalid or expired",
			})
	}

	user := models.User{
		ID: issuer,
	}

	if err := database.DB.First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"user": []string{"not found"},
				},
			})
	}

	err = util.GenerateUserTokens(&user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// set new jwt and refreshjwt cookies
	setCookies(c, user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		// return user object with token and refreshToken
		"user": user,
	})
}

// User retrieves the details of the currently authenticated user.
// It extracts the JWT from the request (header or cookie), parses the user ID from it,
// fetches the user details from the database, and returns the user object.
// Returns 404 if the user associated with the valid token is not found.
func User(c *fiber.Ctx) error {
	// get jwt from header or cookie
	jwt := util.GetJWT(c)
	id, _ := util.ParseJwt(jwt)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	if user.ID == "" {
		c.Status(404)
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"user": []string{"token valid but user not found"},
				},
			})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Logout handles user logout by invalidating the JWT cookie.
// It sets the 'jwt' cookie's expiration time to a past date, effectively removing it.
// Returns a success message.
func Logout(c *fiber.Ctx) error {
	// set expiration time to the past to remove cookie
	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})

}

// UpdateInfo handles updating the authenticated user's profile information (first name, last name, email).
// It parses the updated data from the request body.
// It identifies the user based on the JWT, updates the corresponding user record in the database,
// and returns the updated user object.
func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// get jwt from header or cookie
	jwt := util.GetJWT(c)
	userID, _ := util.ParseJwt(jwt)

	user := models.User{
		ID:        userID,
		FirstName: data[keyFirstName],
		LastName:  data[keyLastName],
		Email:     data[keyEmail],
	}

	database.DB.Model(&user).Where("id = ?", userID).Updates(user)
	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdatePassword handles changing the authenticated user's password.
// It parses the new password and confirmation from the request body.
// It validates that the passwords match, identifies the user via JWT,
// sets the new hashed password for the user, updates the database record,
// and returns the updated user object (excluding password).
// Returns 400 if the passwords do not match.
func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data[keyPassword] != data[keyConfirmPassword] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	// get jwt from header or cookie
	jwt := util.GetJWT(c)
	userID, _ := util.ParseJwt(jwt)

	user := models.User{
		ID: userID,
	}

	user.SetPassword(data[keyPassword])

	database.DB.Model(&user).Updates(user)

	return c.Status(fiber.StatusOK).JSON(user)
}

// RequestResetPassword handles the initiation of a password reset process.
// It parses the user's email from the request body.
// It finds the user by email, generates an encrypted password reset token,
// sends an email to the user containing the reset link/token (implementation pending),
// and returns a success message.
// Returns 404 if the user email is not found.
func RequestResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data[keyEmail]).First(&user)

	if user.ID == "" {
		c.Status(404)
		return c.JSON(fiber.Map{
			"errors": fiber.Map{
				"user": []string{"not found"},
			},
		})
	}

	encToken, err := util.GenerateResetPasswordToken(&user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// send email with token
	err = util.SendResetPasswordEmail(&user, encToken)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// setCookies is a helper function to set the JWT access and refresh tokens
// as HTTP-only cookies in the Fiber context.
// It takes the Fiber context and the user object (which should contain the generated tokens) as input.
func setCookies(c *fiber.Ctx, user models.User) {
	// set jwt cookie
	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   user.AccessToken,
		Expires: time.Now().Add(util.AccessTokenDuration),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:    "refreshjwt",
		Value:   user.RefreshToken,
		Expires: time.Now().Add(util.RefreshTokenDuration),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&refreshCookie)
}
