package handlers

import (
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// PostContact handles contact form submissions and sends a notification email.
func PostContact(c *fiber.Ctx) error {
	var data struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if data.Name == "" || data.Email == "" || data.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "name, email, and message are required",
		})
	}

	if !emailRegex.MatchString(data.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid email address",
		})
	}

	contactTo := app.Email.ContactTo
	if contactTo == "" {
		log.Printf("Contact form from %s <%s>: %s", data.Name, data.Email, data.Message)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "message received",
		})
	}

	if emailService != nil {
		if err := emailService.Send(contactTo, "New contact form message from "+data.Name, "contact", map[string]interface{}{
			"Name":    data.Name,
			"Email":   data.Email,
			"Message": data.Message,
		}); err != nil {
			log.Printf("Failed to queue contact email: %v", err)
		}
	} else {
		log.Printf("Contact form from %s <%s>: %s", data.Name, data.Email, data.Message)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "message received",
	})
}
