import { createReducer, on } from '@ngrx/store';
import { SubscriptionState } from '../interfaces/subscription-state';
import {
  getSubscriptionAction,
  getSubscriptionSuccessAction,
  getSubscriptionFailureAction,
  updateSubscriptionAction,
  updateSubscriptionSuccessAction,
  updateSubscriptionFailureAction,
  cancelSubscriptionAction,
  cancelSubscriptionSuccessAction,
  cancelSubscriptionFailureAction,
  reactivateSubscriptionAction,
  reactivateSubscriptionSuccessAction,
  reactivateSubscriptionFailureAction,
  changePlanAction,
  changePlanSuccessAction,
  changePlanFailureAction,
} from '../actions/subscription.actions';
import { getCurrentUserFailureAction, logoutSuccessAction } from '../actions/auth.actions';

const initialState: SubscriptionState = {
  subscription: null,
  isLoading: true,
  isUpdating: false,
  validationErrors: null,
  redirectUrl: null,
};

export const subscriptionReducer = createReducer(
  initialState,

  // Get subscription
  on(
    getSubscriptionAction,
    (state): SubscriptionState => ({
      ...state,
      isLoading: true,
      validationErrors: null,
    })
  ),
  on(
    getSubscriptionSuccessAction,
    (state, action): SubscriptionState => ({
      ...state,
      isLoading: false,
      subscription: action.subscription,
    })
  ),
  on(
    getSubscriptionFailureAction,
    (state, action): SubscriptionState => ({
      ...state,
      isLoading: false,
      validationErrors: action.errors,
    })
  ),

  // Update subscription
  on(
    updateSubscriptionAction,
    (state): SubscriptionState => ({
      ...state,
      isUpdating: true,
      validationErrors: null,
    })
  ),
  on(
    updateSubscriptionSuccessAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      subscription: action.subscription,
    })
  ),
  on(
    updateSubscriptionFailureAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      validationErrors: action.errors,
    })
  ),

  // Cancel subscription
  on(
    cancelSubscriptionAction,
    (state): SubscriptionState => ({
      ...state,
      isUpdating: true,
      validationErrors: null,
    })
  ),
  on(
    cancelSubscriptionSuccessAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      subscription: action.subscription,
    })
  ),
  on(
    cancelSubscriptionFailureAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      validationErrors: action.errors,
    })
  ),

  // Reactivate subscription
  on(
    reactivateSubscriptionAction,
    (state): SubscriptionState => ({
      ...state,
      isUpdating: true,
      validationErrors: null,
    })
  ),
  on(
    reactivateSubscriptionSuccessAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      subscription: action.subscription,
    })
  ),
  on(
    reactivateSubscriptionFailureAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      validationErrors: action.errors,
    })
  ),

  // Change plan
  on(
    changePlanAction,
    (state): SubscriptionState => ({
      ...state,
      isUpdating: true,
      validationErrors: null,
    })
  ),
  on(
    changePlanSuccessAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      subscription: action.subscription,
    })
  ),
  on(
    changePlanFailureAction,
    (state, action): SubscriptionState => ({
      ...state,
      isUpdating: false,
      validationErrors: action.errors,
    })
  ),

  // Reset subscription state when user is not authenticated
  on(
    getCurrentUserFailureAction,
    logoutSuccessAction,
    (state): SubscriptionState => ({
      ...state,
      subscription: null,
      isLoading: false,
    })
  )
);
