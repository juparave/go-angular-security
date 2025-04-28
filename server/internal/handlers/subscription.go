package handlers

import (
	"net/http"
	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"

	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer" // Corrected import path
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

// CreateCheckoutSession creates a new Stripe Checkout session for a subscription.
// Body: JSON containing { priceId: string }
// Response: 200 with { sessionId: string }, or error.
func CreateCheckoutSession(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get User from database
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Use the global app config for Stripe key early for customer creation if needed
	// TODO: Refactor to use proper dependency injection instead of global 'app'
	if app == nil || app.Stripe.SecretKey == "" {
		fmt.Println("Error: Stripe configuration not initialized in handlers")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server configuration error"})
	}
	stripe.Key = app.Stripe.SecretKey

	// Ensure user has a Stripe Customer ID, create one if not
	stripeCustomerID := user.StripeCustomerID
	if stripeCustomerID == "" {
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Name:  stripe.String(user.FirstName + " " + user.LastName), // Optional: Add user's name
			// Add any other relevant metadata
			Metadata: map[string]string{
				"app_user_id": user.ID, // Removed .String() as user.ID is already a string
			},
		}
		newCustomer, err := customer.New(customerParams) // Use the customer package here
		if err != nil {
			fmt.Printf("Error creating Stripe customer: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create customer record"})
		}
		stripeCustomerID = newCustomer.ID
		user.StripeCustomerID = stripeCustomerID
		// Save the updated user record with the new Stripe Customer ID
		if err := database.DB.Save(&user).Error; err != nil {
			fmt.Printf("Error saving Stripe Customer ID to user %s: %v\n", userID, err)
			// Return an error to be safe.
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user record"})
		}
		fmt.Printf("Created and saved Stripe Customer ID %s for user %s\n", stripeCustomerID, userID)
	}

	// Parse request body for Price ID
	var reqBody struct {
		PriceID string `json:"priceId"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if reqBody.PriceID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing priceId in request body"})
	}

	// Use the global app config for Stripe key
	// TODO: Refactor to use proper dependency injection instead of global 'app'
	if app == nil || app.Stripe.SecretKey == "" {
		// Log the error for debugging
		fmt.Println("Error: Stripe configuration not initialized in handlers")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Server configuration error"})
	}
	stripe.Key = app.Stripe.SecretKey

	// Define success and cancel URLs (replace with your actual frontend URLs)
	// Consider making these configurable
	domain := os.Getenv("DOMAIN") // Or get from config: app.Server.Domain
	if domain == "" {
		domain = "http://localhost:4200" // Default fallback
	}
	successURL := domain + "/subscription/success?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := domain + "/subscription/cancel"

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerID), // Use the potentially newly created customer ID
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(reqBody.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		// Optionally allow promotion codes
		AllowPromotionCodes: stripe.Bool(true),
	}

	// Create the session
	s, err := session.New(params)
	if err != nil {
		fmt.Printf("Error creating Stripe checkout session: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create checkout session"})
	}

	// Return the session ID
	return c.JSON(fiber.Map{"sessionId": s.ID})
}
