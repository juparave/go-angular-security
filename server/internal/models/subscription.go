// Package models contains data structures for application models.
package models

import "time"

// Subscription represents a customer's subscription to a plan.
// It includes Stripe subscription fields and metadata for tracking status and periods.
type Subscription struct {
	ID                 string     `json:"id" gorm:"primaryKey;size:255"` // Unique subscription identifier
	CustomerID         string     `json:"customerId" gorm:"size:255"`    // Associated customer ID
	PlanID             string     `json:"planId" gorm:"size:255"`        // Associated plan ID
	Status             string     `json:"status" gorm:"size:50"`         // Subscription status (e.g., active, canceled)
	CurrentPeriodStart time.Time  `json:"currentPeriodStart"`            // Start of current billing period
	CurrentPeriodEnd   time.Time  `json:"currentPeriodEnd"`              // End of current billing period
	CancelAtPeriodEnd  bool       `json:"cancelAtPeriodEnd"`             // Whether subscription will cancel at period end
	CanceledAt         *time.Time `json:"canceledAt"`                    // When the subscription was canceled (if applicable)
	CreatedAt          time.Time  `json:"createdAt"`                     // Creation timestamp
	UpdatedAt          time.Time  `json:"updatedAt"`                     // Last update timestamp
}
