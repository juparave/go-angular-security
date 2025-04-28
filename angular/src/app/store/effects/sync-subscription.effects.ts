import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { map, filter } from 'rxjs/operators';
import { loginSuccessAction, getCurrentUserSuccessAction } from '../actions/auth.actions';
import { getSubscriptionSuccessAction } from '../actions/subscription.actions';
import { Subscription } from 'src/app/models/subscription';

@Injectable()
export class SyncSubscriptionEffect {
  constructor(private actions$: Actions) { }

  syncSubscription$ = createEffect(() => {
    return this.actions$.pipe(
      // Listen for successful login or user fetch actions
      ofType(loginSuccessAction, getCurrentUserSuccessAction),
      // Ensure currentUser and subscription data exist in the action payload
      filter(action => !!action.currentUser?.subscription),
      // Map the action to dispatch getSubscriptionSuccessAction with the subscription data
      map(action => {
        // Explicitly cast to ensure type safety, though filter should guarantee existence
        const subscription = action.currentUser?.subscription as Subscription;
        return getSubscriptionSuccessAction({ subscription });
      })
    );
  });
}
