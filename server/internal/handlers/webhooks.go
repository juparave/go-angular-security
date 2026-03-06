package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"server/internal/database"
	"server/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
	"github.com/stripe/stripe-go/v82/webhook"
)

func PostStripeWebhook(c *fiber.Ctx) error {
	// Retrieve the raw body
	payload, err := io.ReadAll(bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error reading request body"})
	}

	// log payload
	app.Log.Debug("Webhook payload:", string(payload))

	// Get the Stripe-Signature header
	signatureHeader := c.Get("Stripe-Signature")

	// Get your webhook secret from an environment variable
	endpointSecret := app.Stripe.WebhookSecret

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		app.Log.Error("Webhook signature verification failed:", err)
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
	app.Log.Info("Subscription created event received", "subscription_id", subscription.ID)
	// Get customerID from stripe's subscription object
	customerID := subscription.Customer.ID
	if customerID == "" {
		app.Log.Error("Customer ID not found in subscription")
		return
	}

	// Get user linked with this `customerID`
	var user models.User
	if err := database.DB.
		Preload("Subscription").
		Where("stripe_customer_id = ?", customerID).
		First(&user).Error; err != nil {
		app.Log.Error("User not found with customer ID")
		return
	}

	// Get the Stripe client
	stripeClient := &client.API{}
	stripeClient.Init(app.Stripe.SecretKey, nil)

	// Retrieve the full subscription object from Stripe
	fullSubscription, err := stripeClient.Subscriptions.Get(subscription.ID, nil)
	if err != nil {
		app.Log.Error("Error retrieving full subscription from Stripe:", "err", err)
		return
	}

	var currentPeriodStart, currentPeriodEnd int64
	if len(fullSubscription.Items.Data) > 0 {
		currentPeriodStart = fullSubscription.Items.Data[0].CurrentPeriodStart
		currentPeriodEnd = fullSubscription.Items.Data[0].CurrentPeriodEnd
	}

	// Create user subscription
	user.Subscription = models.Subscription{
		ID:                 subscription.ID,
		CustomerID:         customerID,
		PlanID:             fullSubscription.Items.Data[0].Plan.ID,
		Status:             string(subscription.Status),
		CurrentPeriodStart: time.Unix(currentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(currentPeriodEnd, 0),
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
	}

	// Persist subscription in database
	if err := database.DB.Save(&user.Subscription).Error; err != nil {
		app.Log.Error("Error persisting subscription in database:", "err", err)
		return
	}

	// Persist user in database
	// if err := database.DB.Save(&user); err != nil {
	// 	app.Log.Error("Error persisting user in database:", "err", err)
	// 	return
	// }

	app.Log.Info("Subscription created successfully")

}

func handleSubscriptionUpdated(subscription stripe.Subscription) {
	customerID := subscription.Customer.ID
	if customerID == "" {
		app.Log.Error("Customer ID not found in subscription")
		return
	}

	// Find the user by Stripe Customer ID and preload their subscription
	var user models.User
	if err := database.DB.
		Preload("Subscription").
		Where("stripe_customer_id = ?", customerID).
		First(&user).Error; err != nil {
		app.Log.Error("User not found with customer ID")
		return
	}

	// Get the Stripe client
	stripeClient := &client.API{}
	stripeClient.Init(app.Stripe.SecretKey, nil)

	// Retrieve the full subscription object from Stripe
	fullSubscription, err := stripeClient.Subscriptions.Get(subscription.ID, nil)
	if err != nil {
		app.Log.Error("Error retrieving full subscription from Stripe:", "err", err)
		return
	}

	var currentPeriodStart, currentPeriodEnd int64
	if len(fullSubscription.Items.Data) > 0 {
		currentPeriodStart = fullSubscription.Items.Data[0].CurrentPeriodStart
		currentPeriodEnd = fullSubscription.Items.Data[0].CurrentPeriodEnd
	}

	// Update the user's subscription fields
	sub := &user.Subscription
	sub.PlanID = fullSubscription.Items.Data[0].Plan.ID
	sub.Status = string(subscription.Status)
	sub.CurrentPeriodStart = time.Unix(currentPeriodStart, 0)
	sub.CurrentPeriodEnd = time.Unix(currentPeriodEnd, 0)
	sub.CancelAtPeriodEnd = subscription.CancelAtPeriodEnd

	// Persist the updated subscription in the database
	if err := database.DB.Save(sub).Error; err != nil {
		app.Log.Error("Error updating subscription in database:", "err", err)
		return
	}

	app.Log.Info("Subscription updated successfully")
}

func handleSubscriptionDeleted(subscription stripe.Subscription) {
	customerID := subscription.Customer.ID
	if customerID == "" {
		app.Log.Error("Customer ID not found in subscription")
		return
	}

	// Find the user by Stripe Customer ID and preload their subscription
	var user models.User
	if err := database.DB.
		Preload("Subscription").
		Where("stripe_customer_id = ?", customerID).
		First(&user).Error; err != nil {
		app.Log.Error("User not found with customer ID")
		return
	}

	// Mark the subscription as canceled in the database
	sub := &user.Subscription
	sub.Status = string(subscription.Status)
	sub.CancelAtPeriodEnd = subscription.CancelAtPeriodEnd
	if subscription.CanceledAt != 0 {
		canceledAt := time.Unix(subscription.CanceledAt, 0)
		sub.CanceledAt = &canceledAt
	}

	// Persist the updated subscription in the database
	if err := database.DB.Save(sub).Error; err != nil {
		app.Log.Error("Error updating subscription in database:", "err", err)
		return
	}

	app.Log.Info("Subscription deleted/canceled successfully")
}
