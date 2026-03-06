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
	stripesub "github.com/stripe/stripe-go/v82/subscription"
	stripesubitem "github.com/stripe/stripe-go/v82/subscriptionitem"
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

	// If no Stripe customer, return a virtual free-tier subscription
	if customerID == "" {
		return c.JSON(fiber.Map{
			"id":                   "free_" + account.ID,
			"plan":                 string(account.PlanTier),
			"status":               string(account.SubscriptionStatus),
			"cancel_at_period_end": false,
			"customer":             "",
		})
	}

	stripe.Key = app.Stripe.SecretKey

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
		Status:   stripe.String(string(stripe.SubscriptionStatusActive)),
	}

	iter := stripesub.List(params)
	for iter.Next() {
		sub := iter.Subscription()
		return c.JSON(sub)
	}

	// Fallback: account exists with Stripe customer but no active Stripe subscription
	return c.JSON(fiber.Map{
		"id":                   "free_" + account.ID,
		"plan":                 string(account.PlanTier),
		"status":               string(account.SubscriptionStatus),
		"cancel_at_period_end": false,
		"customer":             customerID,
	})
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
		SuccessURL:          stripe.String(successURL),
		CancelURL:           stripe.String(cancelURL),
		AllowPromotionCodes: stripe.Bool(true),
	}

	s, err := session.New(params)
	if err != nil {
		fmt.Printf("Error creating Stripe checkout session: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create checkout session"})
	}

	return c.JSON(fiber.Map{"sessionId": s.ID})
}

// PatchSubscription updates a subscription (e.g. quantity, metadata).
func PatchSubscription(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account required"})
	}

	subID := c.Params("id")
	if subID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Subscription ID required"})
	}

	var body struct {
		Metadata map[string]string `json:"metadata"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	stripe.Key = app.Stripe.SecretKey

	params := &stripe.SubscriptionParams{}
	if len(body.Metadata) > 0 {
		params.Metadata = body.Metadata
	}

	sub, err := stripesub.Update(subID, params)
	if err != nil {
		fmt.Printf("Error updating subscription %s: %v\n", subID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update subscription"})
	}

	return c.JSON(sub)
}

// PostCancelSubscription sets a subscription to cancel at period end.
func PostCancelSubscription(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account required"})
	}

	subID := c.Params("id")
	if subID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Subscription ID required"})
	}

	stripe.Key = app.Stripe.SecretKey

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	sub, err := stripesub.Update(subID, params)
	if err != nil {
		fmt.Printf("Error canceling subscription %s: %v\n", subID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cancel subscription"})
	}

	// Update account plan tier if we have one linked
	masterDB := database.Manager.GetMasterDB()
	masterDB.Model(&models.Account{}).
		Where("stripe_customer_id = ?", sub.Customer.ID).
		Update("cancel_at_period_end", true)

	return c.JSON(fiber.Map{"message": "Subscription will be canceled at period end", "subscription": sub})
}

// PostReactivateSubscription removes the cancel-at-period-end flag from a subscription.
func PostReactivateSubscription(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account required"})
	}

	subID := c.Params("id")
	if subID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Subscription ID required"})
	}

	stripe.Key = app.Stripe.SecretKey

	// Verify it's currently set to cancel
	existing, err := stripesub.Get(subID, nil)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Subscription not found"})
	}
	if !existing.CancelAtPeriodEnd {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Subscription is not set to cancel"})
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}

	sub, err := stripesub.Update(subID, params)
	if err != nil {
		fmt.Printf("Error reactivating subscription %s: %v\n", subID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reactivate subscription"})
	}

	masterDB := database.Manager.GetMasterDB()
	masterDB.Model(&models.Account{}).
		Where("stripe_customer_id = ?", sub.Customer.ID).
		Update("cancel_at_period_end", false)

	return c.JSON(fiber.Map{"message": "Subscription reactivated", "subscription": sub})
}

// PostChangeSubscription changes the price/plan of an active subscription.
func PostChangeSubscription(c *fiber.Ctx) error {
	accountID := middleware.GetAccountID(c)
	if accountID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Account required"})
	}

	subID := c.Params("id")
	if subID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Subscription ID required"})
	}

	var body struct {
		PriceID string `json:"priceId"`
	}
	if err := c.BodyParser(&body); err != nil || body.PriceID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "priceId is required"})
	}

	stripe.Key = app.Stripe.SecretKey

	// Get current subscription to find the item ID
	existing, err := stripesub.Get(subID, nil)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Subscription not found"})
	}

	if len(existing.Items.Data) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Subscription has no items"})
	}

	// Update the first subscription item's price
	itemID := existing.Items.Data[0].ID
	itemParams := &stripe.SubscriptionItemParams{
		Price: stripe.String(body.PriceID),
	}
	_, err = stripesubitem.Update(itemID, itemParams)
	if err != nil {
		fmt.Printf("Error changing subscription plan %s: %v\n", subID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to change plan"})
	}

	// Return the updated subscription
	updated, err := stripesub.Get(subID, nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Plan changed but failed to retrieve updated subscription"})
	}

	return c.JSON(fiber.Map{"message": "Subscription plan changed", "subscription": updated})
}
