import { User } from '@src/app/models/user';
import { BackendErrors } from '@src/app/store/types/backend-errors.interface';

export interface AuthState {
  isSubmitting: boolean;
  currentUser: User | null;
  isLoggedIn: boolean | null;
  validationErrors: BackendErrors | null;
  isLoading: boolean;
}
