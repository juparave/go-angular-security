import { Component, OnInit } from '@angular/core';
import { SubscriptionService } from 'src/app/services/subscription.service';
import { loadStripe, Stripe } from '@stripe/stripe-js'; // Import Stripe types
import { environment } from 'src/environments/environment'; // For Stripe key

// Declare the stripe object if loaded globally (e.g., via index.html script)
// A better approach is using the @stripe/stripe-js library directly.
declare var Stripe: any; // Or use the imported Stripe type if using the library properly

// Ideally, load these from environment, but hardcoded here for clarity
const PRICE_IDS: Record<string, string> = {
  'basic_monthly': 'price_1RIZP44D300v3DFM7QJThARK',
  'pro_monthly': 'price_1RIZDG4D300v3DFMdmDnDQsS',
  'basic_yearly': 'price_1RIZDI4D300v3DFMQlQ8czSM',
  'pro_yearly': 'price_1RIZDJ4D300v3DFMFkyISsbD'
};

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

  stripePromise: Promise<Stripe | null>; // Hold the Stripe instance promise

  // Inject SubscriptionService
  constructor(private subscriptionService: SubscriptionService) {
    // Initialize Stripe.js asynchronously
    this.stripePromise = loadStripe(environment.stripePublishableKey);
  }

  ngOnInit(): void {
    // Potentially fetch dynamic plan details from backend if needed
  }

  async redirectToCheckout(sessionId: string) {
    const stripe = await this.stripePromise;
    if (!stripe) {
      console.error('Stripe.js has not loaded yet.');
      // Handle error display to user
      return;
    }

    const { error } = await stripe.redirectToCheckout({
      sessionId: sessionId
    });

    // If `redirectToCheckout` fails due to a browser or network
    // error, display the localized error message to your customer
    // using `error.message`.
    if (error) {
      console.error('Error redirecting to Stripe Checkout:', error);
      // Handle error display to user (e.g., show a notification)
    }
  }

  getPriceId(planIdBase: string): string {
    const period = this.isYearly ? 'yearly' : 'monthly';
    const key = `${planIdBase}_${period}`;
    return PRICE_IDS[key] || '';
  }

  selectPlan(planIdBase: string): void {
    const priceId = this.getPriceId(planIdBase);
    console.log(`Selected plan with Price ID: ${priceId}`);
    // Call the service method to initiate checkout
    // This service method should handle the backend call to create a Stripe Checkout session
    // Call the service method to initiate checkout
    this.subscriptionService.createCheckoutSession(priceId).subscribe({
      next: (response) => {
        // Backend returns { sessionId: string }
        if (response && response.sessionId) {
          // Use Stripe.js to redirect to Checkout
          this.redirectToCheckout(response.sessionId);
        } else {
          console.error('Failed to get sessionId from backend.');
          // Handle error display to user
        }
      },
      error: (err) => {
        console.error('Error creating checkout session:', err);
        // Handle error display to user (e.g., show a notification)
      }
    });
  }
}
