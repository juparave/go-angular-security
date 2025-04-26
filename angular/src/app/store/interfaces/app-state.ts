import { AuthState } from './auth-state';
import { SubscriptionState } from './subscription-state';

export interface AppState {
  readonly auth: AuthState;
  readonly subscription: SubscriptionState;
}
