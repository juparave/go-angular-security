import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ReactiveFormsModule,
  UntypedFormBuilder,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { requestResetPasswordAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubmitting, selectValidationErrors } from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { RequestResetPassword } from 'src/app/store/types/request-reset.interface';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { BackendErrorMessagesModule } from 'src/app/shared/modules/backend-error-messages/backend-error-messages.module';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    BackendErrorMessagesModule,
  ],
})
export class ResetPasswordComponent implements OnInit {
  private store = inject<Store<AppState>>(Store);
  private fb = inject(UntypedFormBuilder);

  form: UntypedFormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  constructor() {
    this.form = this.fb.group({
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

  doRequestResetPassword() {
    if (this.form.valid) {
      const request: RequestResetPassword = this.form.value;
      this.store.dispatch(requestResetPasswordAction({ request }));
    }
  }
}
