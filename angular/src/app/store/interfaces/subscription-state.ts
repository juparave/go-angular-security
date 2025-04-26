import { Subscription } from 'src/app/models/subscription';
import { BackendErrors } from '../types/backend-errors.interface';

export interface SubscriptionState {
  subscription: Subscription | null;
  isLoading: boolean;
  isUpdating: boolean;
  validationErrors: BackendErrors | null;
  redirectUrl: string | null;
}
