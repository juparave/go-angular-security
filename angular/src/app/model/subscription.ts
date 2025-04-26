export interface Subscription {
  id: string;
  plan: string; // e.g., "trial", "basic", "pro"
  status: string; // Current status, e.g., "trialing", "active", "canceled"
  trialStart: Date | null; // Start date of the trial period
  trialEnd: Date | null; // End date of the trial period
  stripeCustomerId: string; // Link to Stripe Customer ID
  stripeSubscriptionId: string; // Link to Stripe Subscription ID
  currentPeriodStart: Date | null; // Start date of the current billing period
  currentPeriodEnd: Date | null; // End date of the current billing period
  cancelAtPeriodEnd: boolean; // If true, subscription cancels at current period end
  canceledAt: Date | null; // Timestamp when cancellation occurred
}
