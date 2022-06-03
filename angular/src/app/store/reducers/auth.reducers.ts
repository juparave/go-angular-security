import { Action, createReducer, on } from '@ngrx/store';

import { AuthState } from 'src/app/store/interfaces/auth-state';
import {
  logoutAction,
  logoutFailureAction,
  logoutSuccessAction,
  refreshTokenAction,
  refreshTokenFailureAction,
  refreshTokenSuccessAction,
  registerAction,
  registerFailureAction,
  registerSuccessAction,
} from 'src/app/store/actions/auth.actions';
import {
  loginAction,
  loginFailureAction,
  loginSuccessAction,
} from 'src/app/store/actions/auth.actions';
import {
  getCurrentUserAction,
  getCurrentUserFailureAction,
  getCurrentUserSuccessAction,
} from 'src/app/store/actions/auth.actions';

const initialState: AuthState = {
  isSubmitting: false,
  isLoading: false,
  currentUser: null,
  isLoggedIn: null,
  validationErrors: null,
};

export const authReducer = createReducer(
  initialState,
  on(
    registerAction,
    (state): AuthState => ({
      ...state,
      isSubmitting: true,
      validationErrors: null,
    })
  ),
  on(
    registerSuccessAction,
    (state, action): AuthState => ({
      ...state,
      isSubmitting: false,
      isLoggedIn: true,
      currentUser: action.currentUser,
    })
  ),
  on(
    registerFailureAction,
    (state, action): AuthState => ({
      ...state,
      isSubmitting: false,
      isLoggedIn: false,
      currentUser: null,
      validationErrors: action.errors,
    })
  ),
  on(
    loginAction,
    (state): AuthState => ({
      ...state,
      isSubmitting: true,
      validationErrors: null,
    })
  ),
  on(
    loginSuccessAction,
    (state, action): AuthState => ({
      ...state,
      isSubmitting: false,
      isLoggedIn: true,
      currentUser: action.currentUser,
    })
  ),
  on(
    loginFailureAction,
    (state, action): AuthState => ({
      ...state,
      isSubmitting: false,
      isLoggedIn: false,
      currentUser: null,
      validationErrors: action.errors,
    })
  ),
  on(
    logoutAction,
    (state): AuthState => ({
      ...state,
      isSubmitting: true,
    })
  ),
  on(
    logoutSuccessAction,
    (state, action): AuthState => ({
      ...state,
      isSubmitting: false,
      isLoggedIn: false,
      currentUser: null,
    })
  ),
  on(
    logoutFailureAction,
    (state, action): AuthState => ({
      ...state,
      isSubmitting: false,
    })
  ),
  on(
    getCurrentUserAction,
    (state): AuthState => ({
      ...state,
      isLoading: true,
    })
  ),
  on(
    getCurrentUserSuccessAction,
    (state, action): AuthState => ({
      ...state,
      isLoading: false,
      isLoggedIn: true,
      currentUser: action.currentUser,
    })
  ),
  on(
    getCurrentUserFailureAction,
    (state): AuthState => ({
      ...state,
      isLoading: false,
      isLoggedIn: false,
      currentUser: null,
    })
  ),
  on(
    refreshTokenAction,
    (state): AuthState => ({
      ...state,
      isLoading: true,
    })
  ),
  on(
    refreshTokenSuccessAction,
    (state): AuthState => ({
      ...state,
      isLoading: false,
    })
  ),
  on(
    refreshTokenFailureAction,
    (state): AuthState => ({
      ...state,
      isLoading: false,
      isLoggedIn: false,
      currentUser: null,
    })
  )
);
