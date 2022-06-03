import { createSelector } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';
import { AuthState } from 'src/app/store/interfaces/auth-state';

export const authFeatureSelector = (state: AppState): AuthState => state.auth;

export const isSubmittingSelector = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.isSubmitting
);

export const validationErrorsSelector = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.validationErrors
);

export const isLoggedInSelector = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.isLoggedIn
);

export const isAnonymousSelector = createSelector(
  authFeatureSelector,
  (authState: AuthState) => !(authState.isLoggedIn === true)
);

export const currentUserSelector = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.currentUser
);
