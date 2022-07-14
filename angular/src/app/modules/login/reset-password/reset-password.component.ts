import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { requestResetPasswordAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { isSubmittingSelector, validationErrorsSelector } from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { RequestResetPassword } from 'src/app/store/types/request-reset.interface';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnInit {
  form: FormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  constructor(private store: Store<AppState>, private fb: FormBuilder) {
    this.form = this.fb.group({
      email: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.initializeValues();
  }

  private initializeValues() {
    this.isSubmitting$ = this.store.select(isSubmittingSelector);
    this.backendErrors$ = this.store.select(validationErrorsSelector);
  }

  doRequestResetPassword() {
    if (this.form.valid) {
      const request: RequestResetPassword = this.form.value;
      this.store.dispatch(requestResetPasswordAction({ request }));
    }
  }
}
