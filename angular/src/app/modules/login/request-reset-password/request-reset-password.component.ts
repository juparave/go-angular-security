import { Component, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { requestResetPasswordAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';
import { RequestResetPassword } from 'src/app/store/types/request-reset.interface';
import {
  selectIsSubmiting,
  selectValidationErrors,
} from 'src/app/store/selectors/auth.selectors';

@Component({
  selector: 'app-request-reset-password',
  templateUrl: './request-reset-password.component.html',
  styleUrls: ['./request-reset-password.component.scss']
})
export class RequestResetPasswordComponent implements OnInit {
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
      const request: RequestResetPassword = {
        email: this.form.value.email,
      };
      console.log(request);
      this.store.dispatch(requestResetPasswordAction({ request }));
    }
  }
}
