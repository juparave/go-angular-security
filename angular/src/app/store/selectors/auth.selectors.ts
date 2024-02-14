import { createSelector } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';
import { AuthState } from 'src/app/store/interfaces/auth-state';

export const authFeatureSelector = (state: AppState): AuthState => state.auth;

export const selectIsSubmiting = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.isSubmitting
);

export const selectValidationErrors = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.validationErrors
);

export const selectIsLoggedIn = createSelector(
  authFeatureSelector,
  // (authState: AuthState) => authState.isLoggedIn
  (authState: AuthState) => { return { isLoading: authState.isLoading, isLoggedIn: authState.isLoggedIn } }
);

export const selectIsAnonymous = createSelector(
  authFeatureSelector,
  // (authState: AuthState) => !(authState.isLoggedIn === true)
  (authState: AuthState) => { return { isLoading: authState.isLoading, isAnon: !(authState.isLoggedIn === true) } }
);

export const selectCurrentUser = createSelector(
  authFeatureSelector,
  (authState: AuthState) => authState.currentUser
);
