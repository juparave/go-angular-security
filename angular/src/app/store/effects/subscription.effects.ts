import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap } from 'rxjs';
import { Subscription } from 'src/app/models/subscription';
import {
  getSubscriptionAction,
  getSubscriptionSuccessAction,
  getSubscriptionFailureAction,
} from '../actions/subscription.actions';
import { SubscriptionService } from 'src/app/services/subscription.service';

@Injectable()
export class GetSubscriptionEffect {
  constructor(
    private actions$: Actions,
    private subscriptionService: SubscriptionService
  ) { }

  getSubscription$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(getSubscriptionAction),
      switchMap(() => {
        return this.subscriptionService.getSubscription().pipe(
          map((subscription: Subscription) => {
            return getSubscriptionSuccessAction({ subscription });
          }),
          catchError((errorResponse: HttpErrorResponse) => {
            return of(
              getSubscriptionFailureAction({
                errors: { subscription: [errorResponse.error.message] },
              })
            );
          })
        );
      })
    );
  });
}
