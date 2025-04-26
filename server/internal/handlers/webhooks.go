package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func PostStripeWebhook(c *fiber.Ctx) error {
	// Retrieve the raw body
	payload, err := io.ReadAll(bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error reading request body"})
	}

	// Get the Stripe-Signature header
	signatureHeader := c.Get("Stripe-Signature")

	// Get your webhook secret from an environment variable
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid signature"})
	}

	// Handle the event
	switch event.Type {
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		// Handle subscription created
		handleSubscriptionCreated(subscription)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		// Handle subscription updated
		handleSubscriptionUpdated(subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		// Handle subscription deleted
		handleSubscriptionDeleted(subscription)

	// ... handle other event types as needed

	default:
		return c.Status(http.StatusOK).SendString("Unhandled event type")
	}

	return c.SendStatus(http.StatusOK)
}

func handleSubscriptionCreated(subscription stripe.Subscription) {
	// Implement your logic for handling a created subscription
}

func handleSubscriptionUpdated(subscription stripe.Subscription) {
	// Implement your logic for handling an updated subscription
}

func handleSubscriptionDeleted(subscription stripe.Subscription) {
	// Implement your logic for handling a deleted subscription
}
