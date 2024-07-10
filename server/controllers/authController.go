package controllers

import (
	"fmt"
	"go-app/database"
	"go-app/models"
	"go-app/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": err,
			})
	}

	if data["password"] != data["confirmPassword"] {
		c.Status(400)
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

	err := util.GenerateUserTokens(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"internal server error": err,
			})
	}

	// set cookie with jwt
	setCookies(c, user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == "" {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"errors": fiber.Map{
					"user": []string{"not found"},
				},
			})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
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

	return c.JSON(fiber.Map{
		"message": "success",
		"user":    user,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("refreshToken error:", err)
		return err
	}

	refreshToken := data["refreshToken"]

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

	return c.JSON(fiber.Map{
		"message": "success",
		// return user object with token and refreshToken
		"user": user,
	})
}

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

	return c.JSON(user)
}

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

	return c.JSON(fiber.Map{
		"message": "success",
	})

}

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
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Where("id = ?", userID).Updates(user)
	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
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

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func RequestResetPassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

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

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// setCookies sets jwt and refreshjwt cookies
// must be called after setTokens
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
