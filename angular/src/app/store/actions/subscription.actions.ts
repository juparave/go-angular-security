import { createAction, props } from '@ngrx/store';
import { Subscription } from 'src/app/models/subscription';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';

export enum ActionTypes {
  GET_SUBSCRIPTION = '[Subscription] Get subscription',
  GET_SUBSCRIPTION_SUCCESS = '[Subscription] Get subscription success',
  GET_SUBSCRIPTION_FAILURE = '[Subscription] Get subscription failure',

  UPDATE_SUBSCRIPTION = '[Subscription] Update subscription',
  UPDATE_SUBSCRIPTION_SUCCESS = '[Subscription] Update subscription success',
  UPDATE_SUBSCRIPTION_FAILURE = '[Subscription] Update subscription failure',

  CANCEL_SUBSCRIPTION = '[Subscription] Cancel subscription',
  CANCEL_SUBSCRIPTION_SUCCESS = '[Subscription] Cancel subscription success',
  CANCEL_SUBSCRIPTION_FAILURE = '[Subscription] Cancel subscription failure',

  REACTIVATE_SUBSCRIPTION = '[Subscription] Reactivate subscription',
  REACTIVATE_SUBSCRIPTION_SUCCESS = '[Subscription] Reactivate subscription success',
  REACTIVATE_SUBSCRIPTION_FAILURE = '[Subscription] Reactivate subscription failure',

  CHANGE_PLAN = '[Subscription] Change plan',
  CHANGE_PLAN_SUCCESS = '[Subscription] Change plan success',
  CHANGE_PLAN_FAILURE = '[Subscription] Change plan failure',
}

// Get subscription
export const getSubscriptionAction = createAction(
  ActionTypes.GET_SUBSCRIPTION
);

export const getSubscriptionSuccessAction = createAction(
  ActionTypes.GET_SUBSCRIPTION_SUCCESS,
  props<{ subscription: Subscription }>()
);

export const getSubscriptionFailureAction = createAction(
  ActionTypes.GET_SUBSCRIPTION_FAILURE,
  props<{ errors: BackendErrors }>()
);

// Update subscription
export const updateSubscriptionAction = createAction(
  ActionTypes.UPDATE_SUBSCRIPTION,
  props<{ subscription: Partial<Subscription> }>()
);

export const updateSubscriptionSuccessAction = createAction(
  ActionTypes.UPDATE_SUBSCRIPTION_SUCCESS,
  props<{ subscription: Subscription }>()
);

export const updateSubscriptionFailureAction = createAction(
  ActionTypes.UPDATE_SUBSCRIPTION_FAILURE,
  props<{ errors: BackendErrors }>()
);

// Cancel subscription
export const cancelSubscriptionAction = createAction(
  ActionTypes.CANCEL_SUBSCRIPTION,
  props<{ id: string }>()
);

export const cancelSubscriptionSuccessAction = createAction(
  ActionTypes.CANCEL_SUBSCRIPTION_SUCCESS,
  props<{ subscription: Subscription }>()
);

export const cancelSubscriptionFailureAction = createAction(
  ActionTypes.CANCEL_SUBSCRIPTION_FAILURE,
  props<{ errors: BackendErrors }>()
);

// Reactivate subscription
export const reactivateSubscriptionAction = createAction(
  ActionTypes.REACTIVATE_SUBSCRIPTION,
  props<{ id: string }>()
);

export const reactivateSubscriptionSuccessAction = createAction(
  ActionTypes.REACTIVATE_SUBSCRIPTION_SUCCESS,
  props<{ subscription: Subscription }>()
);

export const reactivateSubscriptionFailureAction = createAction(
  ActionTypes.REACTIVATE_SUBSCRIPTION_FAILURE,
  props<{ errors: BackendErrors }>()
);

// Change plan
export const changePlanAction = createAction(
  ActionTypes.CHANGE_PLAN,
  props<{ id: string; plan: string }>()
);

export const changePlanSuccessAction = createAction(
  ActionTypes.CHANGE_PLAN_SUCCESS,
  props<{ subscription: Subscription }>()
);

export const changePlanFailureAction = createAction(
  ActionTypes.CHANGE_PLAN_FAILURE,
  props<{ errors: BackendErrors }>()
);
