import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Subscription } from '../models/subscription';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class SubscriptionService {
  constructor(private http: HttpClient) { }

  /**
   * Get the current user's subscription details
   */
  getSubscription(): Observable<Subscription> {
    return this.http.get<Subscription>(`${environment.apiUrl}/subscriptions/current`);
  }

  /**
   * Update subscription information
   */
  updateSubscription(subscription: Partial<Subscription>): Observable<Subscription> {
    return this.http.patch<Subscription>(
      `${environment.apiUrl}/subscriptions/${subscription.id}`,
      { subscription }
    );
  }

  /**
   * Cancel a subscription (marks it to be canceled at period end)
   */
  cancelSubscription(id: string): Observable<Subscription> {
    return this.http.post<Subscription>(
      `${environment.apiUrl}/subscriptions/${id}/cancel`,
      {}
    );
  }

  /**
   * Reactivate a canceled subscription
   */
  reactivateSubscription(id: string): Observable<Subscription> {
    return this.http.post<Subscription>(
      `${environment.apiUrl}/subscriptions/${id}/reactivate`,
      {}
    );
  }

  /**
   * Change subscription plan
   */
  changePlan(id: string, plan: string): Observable<Subscription> {
    return this.http.post<Subscription>(
      `${environment.apiUrl}/subscriptions/${id}/change-plan`,
      { plan }
    );
  }

  /**
   * Create a Stripe Checkout session for the given price ID
   * Expects the backend to return an object with a 'url' property for redirection.
   */
  createCheckoutSession(priceId: string): Observable<{ url: string }> {
    return this.http.post<{ url: string }>(
      `${environment.apiUrl}/subscriptions/create-checkout-session`,
      { priceId } // Send the selected price ID to the backend
    );
  }
}
