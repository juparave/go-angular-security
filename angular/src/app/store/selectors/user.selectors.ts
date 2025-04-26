import { createSelector } from '@ngrx/store';
import { AppState } from '../interfaces/app-state';

// Base selector for user state
export const selectUserState = (state: AppState) => state.user;

// Select the current user
export const selectUser = createSelector(
  selectUserState,
  (state): { user: User | null; isLoading: boolean } => ({
    user: state.currentUser,
    isLoading: state.isLoading
  })
);

// Select the user's subscription
export const selectUserSubscription = createSelector(
  selectUserState,
  (state): { subscription: Subscription | null; isLoading: boolean; url?: string } => ({
    subscription: state.currentUser?.subscription || null,
    isLoading: state.isLoading,
    url: state.redirectUrl
  })
);

// Select if the user has a specific subscription plan
export const selectHasSubscriptionPlan = (plans: string[]) =>
  createSelector(
    selectUserSubscription,
    (state): { hasPlan: boolean; isLoading: boolean } => ({
      hasPlan: !!state.subscription && (
        plans.length === 0 ||
        plans.includes(state.subscription.plan)
      ),
      isLoading: state.isLoading
    })
  );

// Select if the user has an active subscription
export const selectHasActiveSubscription = createSelector(
  selectUserSubscription,
  (state): { isActive: boolean; isLoading: boolean } => {
    const subscription = state.subscription;

    const isActive = !!subscription && (
      subscription.status === 'active' ||
      subscription.status === 'trialing' ||
      (subscription.status === 'canceled' && !subscription.cancelAtPeriodEnd) ||
      (subscription.status === 'canceled' &&
        subscription.cancelAtPeriodEnd &&
        subscription.currentPeriodEnd &&
        new Date() < subscription.currentPeriodEnd)
    );

    return {
      isActive,
      isLoading: state.isLoading
    };
  }
);

// Convenience selectors for specific subscription types
export const selectHasTrialSubscription = selectHasSubscriptionPlan(['trial']);
export const selectHasBasicSubscription = selectHasSubscriptionPlan(['basic']);
export const selectHasProSubscription = selectHasSubscriptionPlan(['pro']);
export const selectHasPaidSubscription = selectHasSubscriptionPlan(['basic', 'pro']);
