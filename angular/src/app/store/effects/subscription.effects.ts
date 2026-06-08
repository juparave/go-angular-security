import { HttpErrorResponse } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { Observable, of } from 'rxjs';
import { catchError, map, switchMap } from 'rxjs/operators';

interface SubscriptionApiResponse {
  id: string;
  plan: string;
  status: string;
  trial_start?: string | null;
  trial_end?: string | null;
  customer: string;
  cancel_at_period_end: boolean;
  canceled_at?: number | null;
  items?: { data: { current_period_end: number; current_period_start: number }[] };
}

import { Subscription } from 'src/app/models/subscription';
import { SubscriptionService } from 'src/app/services/subscription.service';
import {
  getSubscriptionAction,
  getSubscriptionSuccessAction,
  getSubscriptionFailureAction,
} from '../actions/subscription.actions';

/**
 * Effect for handling subscription-related actions
 */
@Injectable()
export class GetSubscriptionEffect {
  private actions$ = inject(Actions);
  private subscriptionService = inject(SubscriptionService);

  /**
   * Main effect to fetch subscription data
   */
  getSubscription$ = createEffect(() =>
    this.actions$.pipe(
      ofType(getSubscriptionAction),
      switchMap(() =>
        this.subscriptionService.getSubscription().pipe(
          map((response) =>
            this.mapToSuccessAction(response as unknown as SubscriptionApiResponse)
          ),
          catchError((error) => this.handleError(error))
        )
      )
    )
  );

  /**
   * Maps API response to success action with properly formatted subscription
   */
  private mapToSuccessAction(
    response: SubscriptionApiResponse
  ): ReturnType<typeof getSubscriptionSuccessAction> {
    return getSubscriptionSuccessAction({
      subscription: this.formatSubscription(response),
    });
  }

  /**
   * Transforms raw subscription response to domain model
   */
  private formatSubscription(response: SubscriptionApiResponse): Subscription {
    const subscription: Subscription = {
      id: response.id,
      plan: response.plan,
      status: response.status,
      trialStart: this.toDateOrNull(response.trial_start),
      trialEnd: this.toDateOrNull(response.trial_end),
      stripeCustomerId: response.customer,
      stripeSubscriptionId: response.id,
      cancelAtPeriodEnd: response.cancel_at_period_end,
      canceledAt: this.toDateFromTimestamp(response.canceled_at),
      currentPeriodStart: null,
      currentPeriodEnd: null,
    };

    // Process billing period data if available (not present for free tier)
    if (this.hasItemsData(response)) {
      const item = response.items!.data[0];
      subscription.currentPeriodEnd = this.toDateFromTimestamp(item.current_period_end);
      subscription.currentPeriodStart = this.toDateFromTimestamp(item.current_period_start);
    }

    return subscription;
  }

  /**
   * Handles HTTP errors and creates appropriate failure action
   */
  private handleError(
    error: HttpErrorResponse
  ): Observable<ReturnType<typeof getSubscriptionFailureAction>> {
    console.error('Error fetching subscription:', error);

    const errorMessage = error.error?.message || 'An error occurred while fetching subscription.';

    return of(
      getSubscriptionFailureAction({
        errors: { subscription: [errorMessage] },
      })
    );
  }

  /**
   * Helper method to check if response has items data
   */
  private hasItemsData(response: SubscriptionApiResponse): boolean {
    return (response.items?.data?.length ?? 0) > 0;
  }

  /**
   * Helper method to convert date string to Date or null
   */
  private toDateOrNull(dateString: string | null | undefined): Date | null {
    return dateString ? new Date(dateString) : null;
  }

  /**
   * Helper method to convert timestamp to Date or null
   */
  private toDateFromTimestamp(timestamp: number | null | undefined): Date | null {
    return timestamp ? new Date(timestamp * 1000) : null;
  }
}
