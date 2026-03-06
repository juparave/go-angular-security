import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { Observable, of } from 'rxjs';
import { catchError, map, switchMap } from 'rxjs/operators';

import { Subscription } from 'src/app/models/subscription';
import { SubscriptionService } from 'src/app/services/subscription.service';
import {
  getSubscriptionAction,
  getSubscriptionSuccessAction,
  getSubscriptionFailureAction
} from '../actions/subscription.actions';

/**
 * Effect for handling subscription-related actions
 */
@Injectable()
export class GetSubscriptionEffect {
  constructor(
    private actions$: Actions,
    private subscriptionService: SubscriptionService
  ) { }

  /**
   * Main effect to fetch subscription data
   */
  getSubscription$ = createEffect(() =>
    this.actions$.pipe(
      ofType(getSubscriptionAction),
      switchMap(() =>
        this.subscriptionService.getSubscription().pipe(
          map((response) => this.mapToSuccessAction(response)),
          catchError((error) => this.handleError(error))
        )
      )
    )
  );

  /**
   * Maps API response to success action with properly formatted subscription
   */
  private mapToSuccessAction(response: any): ReturnType<typeof getSubscriptionSuccessAction> {
    return getSubscriptionSuccessAction({
      subscription: this.formatSubscription(response)
    });
  }

  /**
   * Transforms raw subscription response to domain model
   */
  private formatSubscription(response: any): Subscription {
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

    // Process billing period data if available
    if (this.hasItemsData(response)) {
      const item = response.items.data[0];
      subscription.currentPeriodEnd = this.toDateFromTimestamp(item.current_period_end);
      subscription.currentPeriodStart = this.toDateFromTimestamp(item.current_period_start);
    } else {
      console.warn('Subscription data is missing or incomplete:', response);
    }

    return subscription;
  }

  /**
   * Handles HTTP errors and creates appropriate failure action
   */
  private handleError(error: HttpErrorResponse): Observable<ReturnType<typeof getSubscriptionFailureAction>> {
    console.error('Error fetching subscription:', error);

    const errorMessage = error.error?.message || 'An error occurred while fetching subscription.';

    return of(getSubscriptionFailureAction({
      errors: { subscription: [errorMessage] }
    }));
  }

  /**
   * Helper method to check if response has items data
   */
  private hasItemsData(response: any): boolean {
    return response.items?.data?.length > 0;
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