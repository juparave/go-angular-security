import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  AbstractControl,
  ReactiveFormsModule,
  UntypedFormBuilder,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { RouterModule } from '@angular/router';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { MustMatch } from 'src/app/shared/validators/must-match.validator';
import { registerAccountAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubmitting, selectValidationErrors } from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { AccountRegistration } from 'src/app/store/types/account-registration.interface';
import { BackendErrorMessagesModule } from 'src/app/shared/modules/backend-error-messages/backend-error-messages.module';
import { AuthLayoutComponent } from 'src/app/shared/components/auth-layout/auth-layout.component';

@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    BackendErrorMessagesModule,
    AuthLayoutComponent,
  ],
})
export class SignUpComponent implements OnInit {
  private store = inject<Store<AppState>>(Store);
  private fb = inject(UntypedFormBuilder);

  form: UntypedFormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  constructor() {
    this.form = this.fb.group(
      {
        accountName: ['', Validators.required],
        contactEmail: ['', [Validators.required, Validators.email]],
        firstName: ['', Validators.required],
        lastName: ['', Validators.required],
        email: ['', [Validators.required, Validators.email]],
        password: ['', [Validators.required, Validators.minLength(8)]],
        confirmPassword: ['', Validators.required],
      },
      {
        validators: MustMatch('password', 'confirmPassword'),
      }
    );
  }

  get f(): { [key: string]: AbstractControl } {
    return this.form.controls;
  }

  get errorMessages(): { [key: string]: string } {
    return {
      accountName: this.f['accountName'].hasError('required') ? 'Account name is required' : '',
      contactEmail: this.f['contactEmail'].hasError('required')
        ? 'Contact email is required'
        : this.f['contactEmail'].hasError('email')
          ? 'Contact email is invalid'
          : '',
      firstName: this.f['firstName'].hasError('required') ? 'First name is required' : '',
      lastName: this.f['lastName'].hasError('required') ? 'Last name is required' : '',
      email: this.f['email'].hasError('required')
        ? 'Email is required'
        : this.f['email'].hasError('email')
          ? 'Email is invalid'
          : '',
      password: this.f['password'].hasError('required')
        ? 'Password is required'
        : this.f['password'].hasError('minlength')
          ? 'Password must be at least 8 characters'
          : '',
      confirmPassword: this.f['confirmPassword'].hasError('required')
        ? 'Please confirm your password'
        : this.f['confirmPassword'].hasError('mustMatch')
          ? 'Passwords do not match'
          : '',
    };
  }

  ngOnInit(): void {
    this.isSubmitting$ = this.store.select(selectIsSubmitting);
    this.backendErrors$ = this.store.select(selectValidationErrors);
  }

  doRegister() {
    if (this.form.valid) {
      const request: AccountRegistration = { ...this.form.value };
      this.store.dispatch(registerAccountAction({ request }));
    }
  }
}
