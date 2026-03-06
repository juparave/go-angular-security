package handlers

import (
	"fmt"
	"net/http"
	"server/internal/database"
	"server/internal/middleware"
	"server/internal/models"
	"server/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
)

// GetCurrentSubscription returns the current active subscription for the user's account.
func GetCurrentSubscription(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account required"})
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.First(&account, "id = ?", accountID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Account not found"})
	}

	customerID := account.StripeCustomerID
	if customerID == "" {
		return c.Status(http.StatusPaymentRequired).JSON(fiber.Map{"error": "No active subscription found"})
	}

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

// CreateCheckoutSession creates a Stripe checkout session for subscription.
func CreateCheckoutSession(c *fiber.Ctx) error {
	_, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account required"})
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.First(&account, "id = ?", accountID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Account not found"})
	}

	stripe.Key = app.Stripe.SecretKey

	stripeCustomerID := account.StripeCustomerID
	if stripeCustomerID == "" {
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(account.ContactEmail),
			Metadata: map[string]string{
				"account_id": account.ID,
			},
		}
		newCustomer, err := customer.New(customerParams)
		if err != nil {
			fmt.Printf("Error creating Stripe customer: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create customer"})
		}
		stripeCustomerID = newCustomer.ID
		account.StripeCustomerID = stripeCustomerID
		if err := masterDB.Save(&account).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update account"})
		}
		fmt.Printf("Created and saved Stripe Customer ID %s for account %s\n", stripeCustomerID, accountID)
	}

	var reqBody struct {
		PriceID string `json:"priceId"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if reqBody.PriceID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing priceId in request body"})
	}

	domain := app.Domain
	if domain == "" {
		domain = "http://localhost:4200"
	}
	successURL := domain + "/subscription/success?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := domain + "/subscription/cancel"

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerID),
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
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		AllowPromotionCodes: stripe.Bool(true),
	}

	s, err := session.New(params)
	if err != nil {
		fmt.Printf("Error creating Stripe checkout session: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create checkout session"})
	}

	return c.JSON(fiber.Map{"sessionId": s.ID})
}

// PatchSubscription updates a subscription with the provided parameters.
func PatchSubscription(c *fiber.Ctx) error {
	// TODO: Implement subscription patch logic
	return c.Status(http.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented"})
}

// PostCancelSubscription cancels a subscription.
func PostCancelSubscription(c *fiber.Ctx) error {
	// TODO: Implement subscription cancellation logic
	return c.Status(http.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented"})
}

// PostReactivateSubscription reactivates a canceled subscription.
func PostReactivateSubscription(c *fiber.Ctx) error {
	// TODO: Implement subscription reactivation logic
	return c.Status(http.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented"})
}

// PostChangeSubscription changes the plan of a subscription.
func PostChangeSubscription(c *fiber.Ctx) error {
	// TODO: Implement subscription change logic
	return c.Status(http.StatusNotImplemented).JSON(fiber.Map{"error": "Not implemented"})
}
