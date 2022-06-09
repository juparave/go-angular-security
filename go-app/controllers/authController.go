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
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	user.SetPassword(data["password"])

	res := database.DB.Create(&user)
	// verify if user was created
	if res.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

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
		c.Status(404)
		return c.JSON(fiber.Map{
			"errors": fiber.Map{
				"user": []string{"not found"},
			},
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"errors": fiber.Map{
				"password": []string{"incorrect"},
			},
		})
	}

	token, err := util.GenerateJwt(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	refreshToken, err := util.GenerateRefreshJwt(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(util.AccessTokenDuration),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:    "refreshjwt",
		Value:   refreshToken,
		Expires: time.Now().Add(util.RefreshTokenDuration),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&refreshCookie)

	return c.JSON(fiber.Map{
		"message": "success",
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
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "token invalid or expired",
		})
	}

	user := models.User{
		ID: issuer,
	}

	if err := database.DB.First(&user).Error; err != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"errors": fiber.Map{
				"user": []string{"not found"},
			},
		})
	}

	token, err := util.GenerateJwt(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	newRefreshToken, err := util.GenerateRefreshJwt(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user.Token = token
	user.RefreshToken = newRefreshToken

	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24 * 7),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	refreshCookie := fiber.Cookie{
		Name:    "refreshjwt",
		Value:   newRefreshToken,
		Expires: time.Now().Add(time.Hour * 24 * 7),
		// only accesible by backend
		HTTPOnly: true,
	}

	c.Cookie(&refreshCookie)

	return c.JSON(fiber.Map{
		"message": "success",
		// return user object with token and refreshToken
		"user": user,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

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

	cookie := c.Cookies("jwt")

	userID, _ := util.ParseJwt(cookie)

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

	cookie := c.Cookies("jwt")

	userID, _ := util.ParseJwt(cookie)

	user := models.User{
		ID: userID,
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}
