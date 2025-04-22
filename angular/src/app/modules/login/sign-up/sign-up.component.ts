import { Component, OnInit } from '@angular/core';
import {
  AbstractControl,
  UntypedFormBuilder,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { User } from 'src/app/model/user';
import { MustMatch } from 'src/app/shared/validators/must-match.validator';
import { registerAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import {
  selectIsSubmitting,
  selectValidationErrors,
} from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';

@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.scss'],
})
export class SignUpComponent implements OnInit {
  form: UntypedFormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  hidePassword = true;
  hideConfirmPassword = true;

  constructor(private store: Store<AppState>, private fb: UntypedFormBuilder) {
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
