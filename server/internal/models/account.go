package models

import (
	"time"

	"gorm.io/gorm"
)

// AccountStatus represents the status of an account/tenant
type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusSuspended AccountStatus = "suspended"
	AccountStatusCanceled  AccountStatus = "canceled"
)

// SubscriptionTier represents the subscription plan tier
type SubscriptionTier string

const (
	TierFree SubscriptionTier = "free"
	TierPro  SubscriptionTier = "pro"
)

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete"
)

// Account represents a tenant in the multitenant system.
// Each account has its own SQLite database for data isolation.
type Account struct {
	ID        string         `json:"id" gorm:"primaryKey;size:11"`
	Name      string         `json:"name" gorm:"size:128;not null"`
	ContactEmail string      `json:"contactEmail" gorm:"size:128;not null"`
	Status    AccountStatus  `json:"status" gorm:"size:32;default:'active'"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Stripe integration
	StripeCustomerID     string `json:"stripeCustomerId" gorm:"size:255;index"`
	StripeSubscriptionID string `json:"stripeSubscriptionId" gorm:"size:255"`

	// Subscription tracking
	PlanTier           SubscriptionTier   `json:"planTier" gorm:"size:32;default:'free'"`
	SubscriptionStatus SubscriptionStatus `json:"subscriptionStatus" gorm:"size:32"`
	CurrentPeriodEnd   *time.Time         `json:"currentPeriodEnd"`
	CancelAtPeriodEnd  bool               `json:"cancelAtPeriodEnd"`

	// User limits based on plan
	MaxUsers int `json:"maxUsers" gorm:"default:1"`
}

// BeforeCreate generates a short UUID for the account ID
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = SortableShortUUID()
	}
	return nil
}

// IsActive returns true if the account is active
func (a *Account) IsActive() bool {
	return a.Status == AccountStatusActive
}

// HasActiveSubscription returns true if the account has an active paid subscription
func (a *Account) HasActiveSubscription() bool {
	return a.SubscriptionStatus == SubscriptionStatusActive && a.PlanTier == TierPro
}

// CanAddUser checks if the account can add more users based on the plan limit
func (a *Account) CanAddUser(currentCount int) bool {
	return currentCount < a.MaxUsers
}

// TableName returns the table name for the Account model
func (Account) TableName() string {
	return "accounts"
}
