import { createSelector } from '@ngrx/store';
import { AppState } from '../interfaces/app-state';
import { SubscriptionState } from '../interfaces/subscription-state';
import { Subscription } from 'src/app/models/subscription';

export const selectSubscriptionState = (state: AppState): SubscriptionState => state.subscription;

export const selectSubscription = createSelector(
  selectSubscriptionState,
  (state: SubscriptionState): { subscription: Subscription | null; isLoading: boolean; url?: string | null } => ({
    subscription: state.subscription,
    isLoading: state.isLoading,
    url: state.redirectUrl
  })
);

export const selectIsSubscriptionLoading = createSelector(
  selectSubscriptionState,
  (state: SubscriptionState): boolean => state.isLoading
);

export const selectIsSubscriptionUpdating = createSelector(
  selectSubscriptionState,
  (state: SubscriptionState): boolean => state.isUpdating
);

export const selectSubscriptionErrors = createSelector(
  selectSubscriptionState,
  (state: SubscriptionState) => state.validationErrors
);

// Select if the subscription has a specific plan
export const selectHasSubscriptionPlan = (plans: string[]) =>
  createSelector(
    selectSubscription,
    (state): { hasPlan: boolean; isLoading: boolean } => ({
      hasPlan: !!state.subscription && (
        plans.length === 0 ||
        plans.includes(state.subscription.plan)
      ),
      isLoading: state.isLoading
    })
  );

// Select if the subscription is active
export const selectIsSubscriptionActive = createSelector(
  selectSubscription,
  (state): { isActive: boolean; isLoading: boolean } => {
    const subscription = state.subscription;

    const isActive = !!subscription && (
      subscription.status === 'active' ||
      subscription.status === 'trialing' ||
      (subscription.status === 'canceled' && !subscription.cancelAtPeriodEnd) ||
      (subscription.status === 'canceled' &&
        subscription.cancelAtPeriodEnd &&
        subscription.currentPeriodEnd &&
        new Date() < new Date(subscription.currentPeriodEnd))
    );

    return {
      isActive: isActive ?? false,
      isLoading: state.isLoading
    };
  }
);

// Convenience selectors for specific subscription types combined with active check
export const selectHasActiveTrial = createSelector(
  selectSubscription,
  (state): { hasActiveTrial: boolean; isLoading: boolean } => {
    const subscription = state.subscription;

    const hasActiveTrial = !!subscription &&
      subscription.plan === 'trial' && (
        subscription.status === 'active' ||
        subscription.status === 'trialing' ||
        (subscription.status === 'canceled' &&
          subscription.currentPeriodEnd &&
          new Date() < new Date(subscription.currentPeriodEnd))
      );

    return {
      hasActiveTrial: hasActiveTrial ?? false,
      isLoading: state.isLoading
    };
  }
);

export const selectHasActiveBasic = createSelector(
  selectSubscription,
  (state): { hasActiveBasic: boolean; isLoading: boolean } => {
    const subscription = state.subscription;

    const hasActiveBasic = !!subscription &&
      subscription.plan === 'basic' && (
        subscription.status === 'active' ||
        (subscription.status === 'canceled' &&
          subscription.currentPeriodEnd &&
          new Date() < new Date(subscription.currentPeriodEnd))
      );

    return {
      hasActiveBasic: hasActiveBasic ?? false,
      isLoading: state.isLoading
    };
  }
);

export const selectHasActivePro = createSelector(
  selectSubscription,
  (state): { hasActivePro: boolean; isLoading: boolean } => {
    const subscription = state.subscription;

    const hasActivePro = !!subscription &&
      subscription.plan === 'pro' && (
        subscription.status === 'active' ||
        (subscription.status === 'canceled' &&
          subscription.currentPeriodEnd &&
          new Date() < new Date(subscription.currentPeriodEnd))
      );

    return {
      hasActivePro: hasActivePro ?? false,
      isLoading: state.isLoading
    };
  }
);

export const selectHasActivePaid = createSelector(
  selectSubscription,
  (state): { hasActivePaid: boolean; isLoading: boolean } => {
    const subscription = state.subscription;

    const hasActivePaid = !!subscription &&
      (subscription.plan === 'basic' || subscription.plan === 'pro') && (
        subscription.status === 'active' ||
        (subscription.status === 'canceled' &&
          subscription.currentPeriodEnd &&
          new Date() < new Date(subscription.currentPeriodEnd))
      );

    return {
      hasActivePaid: hasActivePaid ?? false,
      isLoading: state.isLoading
    };
  }
);
