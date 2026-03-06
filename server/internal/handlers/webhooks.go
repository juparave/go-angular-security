package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"server/internal/database"
	"server/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
	"github.com/stripe/stripe-go/v82/webhook"
)

// HandleStripeWebhook handles incoming Stripe webhook events
func HandleStripeWebhook(c *fiber.Ctx) error {
	payload, err := io.ReadAll(bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error reading request body"})
	}

	signatureHeader := c.Get("Stripe-Signature")
	endpointSecret := app.Stripe.WebhookSecret

	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid signature"})
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		handleCheckoutSessionCompleted(session)

	case "customer.subscription.created":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		handleSubscriptionCreated(subscription)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		handleSubscriptionUpdated(subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		handleSubscriptionDeleted(subscription)

	case "invoice.paid":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		handleInvoicePaid(invoice)

	case "invoice.payment_failed":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing webhook payload"})
		}
		handleInvoicePaymentFailed(invoice)

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return c.SendStatus(http.StatusOK)
}

// PostStripeWebhook is the legacy handler for backward compatibility
func PostStripeWebhook(c *fiber.Ctx) error {
	return HandleStripeWebhook(c)
}

func handleCheckoutSessionCompleted(session stripe.CheckoutSession) {
	log.Printf("Checkout session completed: %s", session.ID)

	customerID := session.Customer.ID
	if customerID == "" {
		log.Printf("Customer ID not found in checkout session")
		return
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.Where("stripe_customer_id = ?", customerID).First(&account).Error; err != nil {
		log.Printf("Account not found with customer ID: %s", customerID)
		return
	}

	// Update account with subscription info
	subscriptionID := session.Subscription.ID
	if subscriptionID != "" {
		account.StripeSubscriptionID = subscriptionID
		account.SubscriptionStatus = models.SubscriptionStatusActive
		account.PlanTier = models.TierPro
		account.MaxUsers = 10 // Pro tier

		if err := masterDB.Save(&account).Error; err != nil {
			log.Printf("Error updating account: %v", err)
			return
		}
		log.Printf("Account %s upgraded to Pro", account.ID)
	}
}

func handleSubscriptionCreated(subscription stripe.Subscription) {
	log.Printf("Subscription created: %s", subscription.ID)

	customerID := subscription.Customer.ID
	if customerID == "" {
		log.Printf("Customer ID not found in subscription")
		return
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.Where("stripe_customer_id = ?", customerID).First(&account).Error; err != nil {
		log.Printf("Account not found with customer ID: %s", customerID)
		return
	}

	// Update account subscription details
	account.StripeSubscriptionID = subscription.ID
	account.SubscriptionStatus = models.SubscriptionStatus(subscription.Status)
	account.PlanTier = models.TierPro
	account.MaxUsers = 10

	// Get period end from subscription items
	if len(subscription.Items.Data) > 0 {
		periodEnd := time.Unix(subscription.Items.Data[0].CurrentPeriodEnd, 0)
		account.CurrentPeriodEnd = &periodEnd
	}

	if err := masterDB.Save(&account).Error; err != nil {
		log.Printf("Error updating account: %v", err)
		return
	}

	log.Printf("Subscription created for account %s", account.ID)
}

func handleSubscriptionUpdated(subscription stripe.Subscription) {
	log.Printf("Subscription updated: %s", subscription.ID)

	customerID := subscription.Customer.ID
	if customerID == "" {
		log.Printf("Customer ID not found in subscription")
		return
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.Where("stripe_customer_id = ?", customerID).First(&account).Error; err != nil {
		log.Printf("Account not found with customer ID: %s", customerID)
		return
	}

	account.SubscriptionStatus = models.SubscriptionStatus(subscription.Status)
	account.CancelAtPeriodEnd = subscription.CancelAtPeriodEnd

	// Get period end from subscription items
	if len(subscription.Items.Data) > 0 {
		endTime := time.Unix(subscription.Items.Data[0].CurrentPeriodEnd, 0)
		account.CurrentPeriodEnd = &endTime
	}

	if err := masterDB.Save(&account).Error; err != nil {
		log.Printf("Error updating account: %v", err)
		return
	}

	log.Printf("Subscription updated for account %s", account.ID)
}

func handleSubscriptionDeleted(subscription stripe.Subscription) {
	log.Printf("Subscription deleted: %s", subscription.ID)

	customerID := subscription.Customer.ID
	if customerID == "" {
		log.Printf("Customer ID not found in subscription")
		return
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.Where("stripe_customer_id = ?", customerID).First(&account).Error; err != nil {
		log.Printf("Account not found with customer ID: %s", customerID)
		return
	}

	// Downgrade to free tier
	account.StripeSubscriptionID = ""
	account.SubscriptionStatus = models.SubscriptionStatusCanceled
	account.PlanTier = models.TierFree
	account.MaxUsers = 1
	account.CancelAtPeriodEnd = false

	if err := masterDB.Save(&account).Error; err != nil {
		log.Printf("Error updating account: %v", err)
		return
	}

	log.Printf("Subscription canceled for account %s, downgraded to free", account.ID)
}

func handleInvoicePaid(invoice stripe.Invoice) {
	log.Printf("Invoice paid: %s", invoice.ID)

	customerID := invoice.Customer.ID
	if customerID == "" {
		return
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.Where("stripe_customer_id = ?", customerID).First(&account).Error; err != nil {
		return
	}

	// Update period end if available
	if invoice.PeriodEnd > 0 {
		endTime := time.Unix(invoice.PeriodEnd, 0)
		account.CurrentPeriodEnd = &endTime
		masterDB.Save(&account)
	}
}

func handleInvoicePaymentFailed(invoice stripe.Invoice) {
	log.Printf("Invoice payment failed: %s", invoice.ID)

	customerID := invoice.Customer.ID
	if customerID == "" {
		return
	}

	masterDB := database.Manager.GetMasterDB()
	var account models.Account
	if err := masterDB.Where("stripe_customer_id = ?", customerID).First(&account).Error; err != nil {
		return
	}

	// Mark subscription as past due
	account.SubscriptionStatus = models.SubscriptionStatusPastDue
	masterDB.Save(&account)
}

// getStripeClient returns a Stripe client initialized with the app's secret key
func getStripeClient() *client.API {
	stripeClient := &client.API{}
	stripeClient.Init(app.Stripe.SecretKey, nil)
	return stripeClient
}
