import { Component, OnInit } from '@angular/core';
import { SubscriptionService } from 'src/app/services/subscription.service'; // Assuming service exists

interface PlanFeature {
  name: string;
  included: boolean;
}

interface Plan {
  id: string; // Corresponds to Stripe Price ID (e.g., price_basic_monthly, price_basic_yearly)
  name: string;
  description: string;
  priceMonthly: number;
  priceYearly: number;
  features: PlanFeature[];
}

@Component({
  selector: 'app-subscription-page',
  templateUrl: './subscription-page.component.html',
  styleUrls: ['./subscription-page.component.scss']
})
export class SubscriptionPageComponent implements OnInit {

  isYearly: boolean = false; // Default to monthly
  plans: Plan[] = [
    {
      id: 'basic', // Base ID, will append _monthly or _yearly later
      name: 'Basic',
      description: 'Essential features for individuals.',
      priceMonthly: 10,
      priceYearly: 100, // e.g., 2 months free
      features: [
        { name: 'Feature A', included: true },
        { name: 'Feature B', included: true },
        { name: 'Feature C', included: false },
        { name: 'Feature D', included: false },
      ]
    },
    {
      id: 'pro',
      name: 'Pro',
      description: 'Advanced features for professionals.',
      priceMonthly: 25,
      priceYearly: 250,
      features: [
        { name: 'Feature A', included: true },
        { name: 'Feature B', included: true },
        { name: 'Feature C', included: true },
        { name: 'Feature D', included: false },
      ]
    },
    {
      id: 'enterprise',
      name: 'Enterprise',
      description: 'Comprehensive solutions for teams.',
      priceMonthly: 50,
      priceYearly: 500,
      features: [
        { name: 'Feature A', included: true },
        { name: 'Feature B', included: true },
        { name: 'Feature C', included: true },
        { name: 'Feature D', included: true },
      ]
    }
  ];

  // Inject SubscriptionService if needed for checkout logic
  constructor(private subscriptionService: SubscriptionService) { }

  ngOnInit(): void {
    // Potentially fetch dynamic plan details from backend if needed
  }

  getPriceId(planIdBase: string): string {
    // Construct the Stripe Price ID based on the toggle state
    // IMPORTANT: These IDs must match the Price IDs created in your Stripe account
    return `price_${planIdBase}_${this.isYearly ? 'yearly' : 'monthly'}`;
  }

  selectPlan(planIdBase: string): void {
    const priceId = this.getPriceId(planIdBase);
    console.log(`Selected plan with Price ID: ${priceId}`);
    // Call the service method to initiate checkout
    // This service method should handle the backend call to create a Stripe Checkout session
    // and redirect the user to Stripe.
    this.subscriptionService.createCheckoutSession(priceId).subscribe({
      next: (response) => {
        // Assuming the backend returns a { url: string } for Stripe redirect
        if (response && response.url) {
          window.location.href = response.url;
        } else {
          console.error('Failed to get checkout URL from backend.');
          // Handle error display to user
        }
      },
      error: (err) => {
        console.error('Error creating checkout session:', err);
        // Handle error display to user
      }
    });
  }
}
