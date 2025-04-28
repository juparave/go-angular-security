package handlers

import (
	"net/http"
	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/subscription"
)

// GetCurrentSubscription returns the current active subscription for a given customer.
// Route param: customerId
// Response: 200 with subscription JSON, or 404 if not found.
func GetCurrentSubscription(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// get User from database
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	customerID := user.StripeCustomerID
	if customerID == "" {
		// This user does not have a Stripe CustomerID, so we asume
		// no subscription either
		return c.Status(http.StatusPaymentRequired).JSON(fiber.Map{"error": "No active subscription found"})
	}

	// TODO: find a better place to do this initiation
	stripe.Key = app.Stripe.SecretKey

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
		Status:   stripe.String(string(stripe.SubscriptionStatusActive)),
	}

	iter := subscription.List(params)
	for iter.Next() {
		sub := iter.Subscription()
		return c.JSON(sub)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No active subscription found"})
}

// PatchSubscription updates a subscription with the provided parameters.
// Route param: subscriptionId
// Body: JSON matching stripe.SubscriptionParams
// Response: 200 with updated subscription JSON, or error.
func PatchSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscriptionId")
	var updateParams stripe.SubscriptionParams
	if err := c.BodyParser(&updateParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedSub, err := subscription.Update(subscriptionID, &updateParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedSub)
}

// PostCancelSubscription sets CancelAtPeriodEnd to true for a subscription, scheduling cancellation.
// Route param: subscriptionId
// Response: 200 with updated subscription JSON, or error.
func PostCancelSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscriptionId")

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	canceledSub, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(canceledSub)
}

// PostReactivateSubscription reactivates a subscription by setting CancelAtPeriodEnd to false.
// Route param: subscriptionId
// Response: 200 with updated subscription JSON, or error.
func PostReactivateSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscriptionId")

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}

	reactivatedSub, err := subscription.Update(subscriptionID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reactivatedSub)
}

// PostChangeSubscription updates a subscription with custom parameters.
// Route param: subscriptionId
// Body: JSON matching stripe.SubscriptionParams
// Response: 200 with updated subscription JSON, or error.
func PostChangeSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscriptionId")
	var changeParams stripe.SubscriptionParams
	if err := c.BodyParser(&changeParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	changedSub, err := subscription.Update(subscriptionID, &changeParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(changedSub)
}
