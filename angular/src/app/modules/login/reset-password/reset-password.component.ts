import { Component, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { requestResetPasswordAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubmiting, selectValidationErrors } from 'src/app/store/selectors/auth.selectors';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { RequestResetPassword } from 'src/app/store/types/request-reset.interface';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnInit {
  form: UntypedFormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  constructor(private store: Store<AppState>, private fb: UntypedFormBuilder) {
    this.form = this.fb.group({
      email: ['', Validators.required],
    });
  }

  ngOnInit(): void {
    this.initializeValues();
  }

  private initializeValues() {
    this.isSubmitting$ = this.store.select(selectIsSubmiting);
    this.backendErrors$ = this.store.select(selectValidationErrors);
  }

  doRequestResetPassword() {
    if (this.form.valid) {
      const request: RequestResetPassword = this.form.value;
      this.store.dispatch(requestResetPasswordAction({ request }));
    }
  }
}
