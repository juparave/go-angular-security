import { Component, OnInit } from '@angular/core';
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
import { User } from 'src/app/models/user';
import { MustMatch } from 'src/app/shared/validators/must-match.validator';
import { registerAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubmitting, selectValidationErrors } from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { BackendErrorMessagesModule } from 'src/app/shared/modules/backend-error-messages/backend-error-messages.module';

@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    BackendErrorMessagesModule,
  ],
})
export class SignUpComponent implements OnInit {
  form: UntypedFormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  hidePassword = true;
  hideConfirmPassword = true;

  constructor(
    private store: Store<AppState>,
    private fb: UntypedFormBuilder
  ) {
    this.form = this.fb.group(
      {
        firstName: [''],
        lastName: [''],
        email: ['', [Validators.required, Validators.email]],
        password: ['', [Validators.required, Validators.minLength(8)]],
        confirmPassword: ['', Validators.required],
      },
      {
        validators: MustMatch('password', 'confirmPassword'),
      }
    );
  }

  // convenience getter to access form fields
  get f(): { [key: string]: AbstractControl } {
    return this.form.controls;
  }

  get errorMessages(): { [key: string]: string } {
    return {
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
        ? 'Confirm Password is required'
        : this.f['confirmPassword'].hasError('mustMatch')
          ? 'Passwords must match'
          : '',
    };
  }

  ngOnInit(): void {
    this.initializeValues();
  }

  private initializeValues() {
    this.isSubmitting$ = this.store.select(selectIsSubmitting);
    this.backendErrors$ = this.store.select(selectValidationErrors);
  }

  doRegister() {
    if (this.form.valid) {
      const request: User = {
        ...this.form.value,
      };
      this.store.dispatch(registerAction({ request }));
    }
  }
}
