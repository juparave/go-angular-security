import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { map, filter } from 'rxjs/operators';
import {
  loginSuccessAction,
  getCurrentUserSuccessAction,
  registerAccountSuccessAction,
} from '../actions/auth.actions';
import {
  getSubscriptionAction,
  getSubscriptionSuccessAction,
} from '../actions/subscription.actions';
import { Subscription } from 'src/app/models/subscription';

@Injectable()
export class SyncSubscriptionEffect {
  constructor(private actions$: Actions) {}

  // When user data includes inline subscription, use it directly
  syncSubscription$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(loginSuccessAction, getCurrentUserSuccessAction, registerAccountSuccessAction),
      filter((action) => !!action.currentUser?.subscription),
      map((action) => {
        const subscription = action.currentUser?.subscription as Subscription;
        return getSubscriptionSuccessAction({ subscription });
      })
    );
  });

  // When user data has no inline subscription, fetch it from the API
  fetchSubscription$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(loginSuccessAction, getCurrentUserSuccessAction, registerAccountSuccessAction),
      filter((action) => !action.currentUser?.subscription),
      map(() => getSubscriptionAction())
    );
  });
}
