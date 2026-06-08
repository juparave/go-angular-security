import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Store } from '@ngrx/store';
import { AuthService } from '../../../../services/auth/auth.service';
import { AppState } from '../../../../store/interfaces/app-state';

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <div class="card bg-base-100 shadow-xl max-w-md mx-auto">
      <div class="card-body">
        <h2 class="card-title text-2xl mb-4">Change Password</h2>

        <form [formGroup]="changePasswordForm" (ngSubmit)="onSubmit()" class="space-y-4">
          <!-- Current Password -->
          <div class="form-control">
            <label class="label" for="currentPassword">
              <span class="label-text">Current Password</span>
            </label>
            <input
              id="currentPassword"
              type="password"
              formControlName="currentPassword"
              placeholder="Enter current password"
              class="input input-bordered w-full"
              [class.input-error]="isFieldInvalid('currentPassword')"
            />
            @if (isFieldInvalid('currentPassword')) {
              <span class="label-text-alt text-error">Current password is required</span>
            }
          </div>

          <!-- New Password -->
          <div class="form-control">
            <label class="label" for="newPassword">
              <span class="label-text">New Password</span>
            </label>
            <input
              id="newPassword"
              type="password"
              formControlName="newPassword"
              placeholder="Enter new password"
              class="input input-bordered w-full"
              [class.input-error]="isFieldInvalid('newPassword')"
            />
            @if (isFieldInvalid('newPassword')) {
              <span class="label-text-alt text-error">
                Password must be at least 8 characters
              </span>
            }
          </div>

          <!-- Confirm Password -->
          <div class="form-control">
            <label class="label" for="confirmPassword">
              <span class="label-text">Confirm New Password</span>
            </label>
            <input
              id="confirmPassword"
              type="password"
              formControlName="confirmPassword"
              placeholder="Confirm new password"
              class="input input-bordered w-full"
              [class.input-error]="isFieldInvalid('confirmPassword') || passwordMismatch()"
            />
            @if (passwordMismatch()) {
              <span class="label-text-alt text-error">Passwords do not match</span>
            }
          </div>

          <!-- Error Message -->
          @if (errorMessage) {
            <div class="alert alert-error">
              <span>{{ errorMessage }}</span>
            </div>
          }

          <!-- Success Message -->
          @if (successMessage) {
            <div class="alert alert-success">
              <span>{{ successMessage }}</span>
            </div>
          }

          <!-- Submit Button -->
          <div class="form-control mt-6">
            <button
              type="submit"
              [disabled]="changePasswordForm.invalid || isSubmitting"
              class="btn btn-primary"
            >
              @if (isSubmitting) {
                <span class="loading loading-spinner"></span>
              }
              Change Password
            </button>
          </div>
        </form>
      </div>
    </div>
  `,
  styles: [
    `
      :host {
        display: block;
      }
    `,
  ],
})
export class ChangePasswordComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private store = inject(Store<AppState>);

  changePasswordForm: FormGroup;
  isSubmitting = false;
  errorMessage = '';
  successMessage = '';

  constructor() {
    this.changePasswordForm = this.fb.group(
      {
        currentPassword: ['', [Validators.required]],
        newPassword: ['', [Validators.required, Validators.minLength(8)]],
        confirmPassword: ['', [Validators.required]],
      },
      { validators: this.passwordMatchValidator }
    );
  }

  passwordMatchValidator(form: FormGroup) {
    const newPassword = form.get('newPassword')?.value;
    const confirmPassword = form.get('confirmPassword')?.value;
    return newPassword === confirmPassword ? null : { mismatch: true };
  }

  isFieldInvalid(fieldName: string): boolean {
    const field = this.changePasswordForm.get(fieldName);
    return !!(field?.invalid && (field?.dirty || field?.touched));
  }

  passwordMismatch(): boolean {
    return (
      this.changePasswordForm.hasError('mismatch') &&
      this.changePasswordForm.get('confirmPassword')?.touched
    );
  }

  onSubmit(): void {
    if (this.changePasswordForm.invalid || this.isSubmitting) {
      return;
    }

    this.isSubmitting = true;
    this.errorMessage = '';
    this.successMessage = '';

    const { currentPassword, newPassword, confirmPassword } = this.changePasswordForm.value;

    this.authService
      .changePassword({
        currentPassword,
        newPassword,
        confirmPassword,
      })
      .subscribe({
        next: () => {
          this.isSubmitting = false;
          this.successMessage = 'Password changed successfully!';
          this.changePasswordForm.reset();
        },
        error: (error) => {
          this.isSubmitting = false;
          if (error.error?.errors?.currentPassword) {
            this.errorMessage = 'Current password is incorrect';
          } else {
            this.errorMessage = error.error?.message || 'Failed to change password';
          }
        },
      });
  }
}
