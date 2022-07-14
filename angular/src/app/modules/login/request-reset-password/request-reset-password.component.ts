import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { AppState } from 'src/app/store/interfaces/app-state';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';

@Component({
  selector: 'app-request-reset-password',
  templateUrl: './request-reset-password.component.html',
  styleUrls: ['./request-reset-password.component.scss']
})
export class RequestResetPasswordComponent implements OnInit {
  form: FormGroup;
  isSubmitting$!: Observable<boolean>;
  backendErrors$!: Observable<BackendErrors | null>;

  constructor(private store: Store<AppState>, private fb: FormBuilder) {
    this.form = this.fb.group({
      email: ['', Validators.required],
    });
  }

  ngOnInit(): void {
  }

  doRequestResetPassword() {
  }
}
