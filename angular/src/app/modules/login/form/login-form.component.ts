import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  UntypedFormBuilder,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { loginAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubmitting, selectValidationErrors } from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { LoginRequest } from 'src/app/store/types/login-request.interface';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { BackendErrorMessagesModule } from 'src/app/shared/modules/backend-error-messages/backend-error-messages.module';

@Component({
  selector: 'app-login-form',
  templateUrl: './login-form.component.html',
  styleUrls: ['./login-form.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MatCardModule,
    MatProgressBarModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    BackendErrorMessagesModule,
  ],
})
export class LoginFormComponent implements OnInit {
  loginForm: UntypedFormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;
  route = inject(ActivatedRoute);

  hidePassword = true;

  constructor(
    private store: Store<AppState>,
    private fb: UntypedFormBuilder
  ) {
    this.loginForm = this.fb.group({
      email: ['', Validators.required],
      password: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.initializeValues();
  }

  private initializeValues() {
    this.isSubmitting$ = this.store.select(selectIsSubmitting);
    this.backendErrors$ = this.store.select(selectValidationErrors);
  }

  doLogin() {
    if (this.loginForm.valid) {
      const request: LoginRequest = {
        user: this.loginForm.value,
      };
      // get return url from query parameters or default to home page
      const returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/app';
      this.store.dispatch(loginAction({ request, returnUrl }));
    }
  }
}
