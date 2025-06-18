package controllers

import (
	"go-app/database"
	"go-app/models"
	"go-app/util"
	"strings" // Added strings import
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": err,
			})
	}

	// Validate required fields
	if data["firstName"] == "" || data["lastName"] == "" || data["email"] == "" || data["password"] == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "all fields (firstName, lastName, email, password) are required",
			})
	}

	if data["password"] != data["confirmPassword"] {
		c.Status(400) // This line is redundant as Status is set by the return statement
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "passwords do not match",
			})
	}

	user := models.User{
		FirstName: data["firstName"],
		LastName:  data["lastName"],
		Email:     data["email"],
	}

	user.SetPassword(data["password"])

	res := database.DB.Create(&user)
	// verify if user was created
	if res.Error != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": res.Error.Error(),
			})
	}

	errToken := util.GenerateUserTokens(&user)
	if errToken != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"internal server error": errToken.Error(),
			})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body not readable or not JSON"})
	}

	email, emailExists := data["email"]
	password, passwordExists := data["password"]

	if !emailExists || email == "" || !passwordExists || password == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "email and password are required",
			})
	}

	var user models.User
	database.DB.Where("email = ?", email).First(&user)

	if user.ID == "" {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"user": []string{"not found"},
				},
			})
	}

	if err := user.ComparePassword(password); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"password": []string{"incorrect"},
				},
			})
	}

	tokenErr := util.GenerateUserTokens(&user)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Login successful, but could not generate tokens: " + tokenErr.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body not readable or not JSON"})
	}

	refreshToken, tokenExists := data["refreshToken"]
	if !tokenExists || refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "refreshToken is required",
			})
	}

	issuer, err := util.ParseJwt(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "token invalid or expired",
			})
	}

	var user models.User // Declare user here
	user.ID = issuer // Assign issuer to ID

	if dbErr := database.DB.First(&user).Error; dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{
					"errors": fiber.Map{
						"user": []string{"not found"},
					},
				})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error: " + dbErr.Error()})
	}

	tokenErr := util.GenerateUserTokens(&user)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Token generation failed: " + tokenErr.Error(),
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:     "refreshjwt",
		Value:    user.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	}
	c.Cookie(&refreshCookie)

	return c.JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

func User(c *fiber.Ctx) error {
	jwtToken := util.GetJWT(c)
	if jwtToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated: Missing token",
		})
	}

	id, err := util.ParseJwt(jwtToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated: Invalid token",
		})
	}

	if id == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated: Invalid token resulted in empty user identifier",
		})
	}

	var user models.User
	dbResult := database.DB.Where("id = ?", id).First(&user)

	if dbResult.Error != nil {
		if dbResult.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{
					"errors": fiber.Map{
						"user": []string{"token valid but user not found"},
					},
				})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error: " + dbResult.Error.Error()})
	}

	if user.ID == "" {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"user": []string{"user not found despite successful query"},
				},
			})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:     "refreshjwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&refreshCookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body not readable or not JSON"})
	}

	jwtToken := util.GetJWT(c)
	if jwtToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated: Missing token"})
	}

	userID, err := util.ParseJwt(jwtToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated: Invalid token"})
	}
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated: Invalid token resulted in empty user identifier"})
	}

	// Fetch the existing user to update
	var user models.User // Single declaration of 'user' for this function scope
	if res := database.DB.First(&user, "id = ?", userID); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error: " + res.Error.Error()})
	}

	updates := make(map[string]interface{})
	if val, ok := data["first_name"]; ok {
		updates["first_name"] = val
	}
	if val, ok := data["last_name"]; ok {
		updates["last_name"] = val
	}
	if val, ok := data["email"]; ok {
		updates["email"] = val
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusOK).JSON(user)
	}

	if res := database.DB.Model(&user).Updates(updates); res.Error != nil {
		if strings.Contains(res.Error.Error(), "UNIQUE") || strings.Contains(res.Error.Error(), "duplicate") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": res.Error.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating user: " + res.Error.Error()})
	}

	database.DB.First(&user, "id = ?", userID)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body not readable or not JSON"})
	}

	jwtToken := util.GetJWT(c)
	if jwtToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated: Missing token"})
	}

	userID, err := util.ParseJwt(jwtToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated: Invalid token"})
	}
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated: Invalid token resulted in empty user identifier"})
	}

	// Validate presence of password and password_confirm
	newPassword, newPassExists := data["password"]
	confirmPassword, confirmPassExists := data["password_confirm"]

	if !newPassExists || !confirmPassExists {
		// This error message might be too generic if only one is missing.
		// The original code just checks newPassword != confirmPassword.
		// If one is missing, it will be empty string, so "passwords do not match" might still be hit.
		// For more clarity:
		if !newPassExists && !confirmPassExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password and password_confirm are required"})
		}
		if !newPassExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password is required"})
		}
		// Note: The original code would hit "passwords do not match" if confirmPassExists is false and newPassExists is true.
		// To maintain that behavior or make it more specific for missing confirm:
		if !confirmPassExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password_confirm is required"})
		}
	}

	// Validate password is not empty
	if newPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "password cannot be empty"})
	}

	if newPassword != confirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "passwords do not match"})
	}

	// var user models.User // Not strictly needed to fetch full user before this kind of update
	// if res := database.DB.First(&user, "id = ?", userID); res.Error != nil {
	// 	if res.Error == gorm.ErrRecordNotFound {
	// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	// 	}
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error: " + res.Error.Error()})
	// }

	userToUpdate := models.User{ID: userID}
	userToUpdate.SetPassword(newPassword)

	if res := database.DB.Model(&models.User{ID: userID}).Select("password").Updates(userToUpdate); res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating password: " + res.Error.Error()})
	}

	var updatedUser models.User
	database.DB.First(&updatedUser, "id = ?", userID)

	return c.JSON(updatedUser)
}

func RequestResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body not readable or not JSON"})
	}

	email, emailExists := data["email"]
	if !emailExists || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "email is required"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Option 1: Return 404 (current behavior being tested)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": fiber.Map{"user": []string{"not found"}}})
			// Option 2: Return 200 OK to prevent email enumeration (common practice)
			// return c.JSON(fiber.Map{"message": "If your email address exists in our system, you will receive a password reset link."})
		}
		// Other DB errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error: " + err.Error()})
	}

	// If user.ID is empty even if no DB error (should not happen with GORM's First)
	if user.ID == "" { // This check is somewhat redundant if First(&user) above succeeded without error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": fiber.Map{"user": []string{"not found"}}})
	}

	encToken, err := util.GenerateResetPasswordToken(&user)
	if err != nil {
		// Log this server-side for diagnostics
		// log.Printf("Error generating reset password token for user %s: %v", user.Email, err)
		return c.SendStatus(fiber.StatusInternalServerError) // No body sent
	}

	err = util.SendResetPasswordEmail(&user, encToken)
	if err != nil {
		// Log this server-side
		// log.Printf("Error sending reset password email for user %s: %v", user.Email, err)
		return c.SendStatus(fiber.StatusInternalServerError) // No body sent
	}

	return c.JSON(fiber.Map{
		"message": "success", // Generic success message
	})
}
