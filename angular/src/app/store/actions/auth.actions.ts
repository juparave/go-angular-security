import { createAction, props } from '@ngrx/store';
import { User } from 'src/app/model/user';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { LoginRequest } from 'src/app/store/types/login-request.interface';
import { RequestResetPassword } from 'src/app/store/types/request-reset.interface';

export enum ActionTypes {
  REGISTER = '[Auth] Register',
  REGISTER_SUCCESS = '[Auth] Register success',
  REGISTER_FAILURE = '[Auth] Register failure',

  LOGIN = '[Auth] Login',
  LOGIN_SUCCESS = '[Auth] Login success',
  LOGIN_FAILURE = '[Auth] Login failure',

  LOGOUT = '[Auth] Logout',
  LOGOUT_SUCCESS = '[Auth] Logout succes',
  LOGOUT_FAILURE = '[Auth] Logout failure',

  REFRESH_TOKEN = '[Auth] Refresh token',
  REFRESH_TOKEN_SUCCESS = '[Auth] Refresh token succes',
  REFRESH_TOKEN_FAILURE = '[Auth] Refresh token failure',

  GET_CURRENT_USER = '[Auth] Get current user',
  GET_CURRENT_USER_SUCCESS = '[Auth] Get current user success',
  GET_CURRENT_USER_FAILURE = '[Auth] Get current user failure',

  REQUEST_RESET_PASSWORD = '[Auth] Request reset password',
  VALIDATE_RESET_PASSWORD = '[Auth] Validate reset password',

  RESET_PASSWORD = '[Auth] Reset password',
  RESET_PASSWORD_SUCCESS = '[Auth] Reset password sucess',
  RESET_PASSWORD_FAILURE = '[Auth] Reset password failure',
}

export const registerAction = createAction(
  ActionTypes.REGISTER,
  props<{ request: User }>()
);

export const registerSuccessAction = createAction(
  ActionTypes.REGISTER_SUCCESS,
  props<{ currentUser: User }>()
);

export const registerFailureAction = createAction(
  ActionTypes.REGISTER_FAILURE,
  props<{ errors: BackendErrors }>()
);

export const loginAction = createAction(
  ActionTypes.LOGIN,
  props<{ request: LoginRequest }>()
);

export const loginSuccessAction = createAction(
  ActionTypes.LOGIN_SUCCESS,
  props<{ currentUser: User }>()
);

export const loginFailureAction = createAction(
  ActionTypes.LOGIN_FAILURE,
  props<{ errors: BackendErrors }>()
);

export const logoutAction = createAction(ActionTypes.LOGOUT);

export const logoutSuccessAction = createAction(ActionTypes.LOGOUT_SUCCESS);

export const logoutFailureAction = createAction(
  ActionTypes.LOGOUT_FAILURE,
  props<{ errors: BackendErrors }>()
);

export const refreshTokenAction = createAction(ActionTypes.REFRESH_TOKEN);

export const refreshTokenSuccessAction = createAction(
  ActionTypes.REFRESH_TOKEN_SUCCESS,
  props<{ currentUser: User }>()
);

export const refreshTokenFailureAction = createAction(
  ActionTypes.REFRESH_TOKEN_FAILURE,
  props<{ errors: BackendErrors }>()
);

export const getCurrentUserAction = createAction(ActionTypes.GET_CURRENT_USER);

export const getCurrentUserSuccessAction = createAction(
  ActionTypes.GET_CURRENT_USER_SUCCESS,
  props<{ currentUser: User }>()
);

export const getCurrentUserFailureAction = createAction(
  ActionTypes.GET_CURRENT_USER_FAILURE
);

export const requestResetPasswordAction = createAction(
  ActionTypes.REQUEST_RESET_PASSWORD,
  props<{ request: RequestResetPassword }>()
);
