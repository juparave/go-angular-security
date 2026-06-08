import { Component, EventEmitter, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { TeamService } from '../../../services/team.service';

@Component({
  selector: 'app-team-invite',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <div class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Invite Team Member</h3>

        <form [formGroup]="inviteForm" (ngSubmit)="onSubmit()" class="space-y-4">
          <!-- Email -->
          <div class="form-control">
            <label class="label" for="email">
              <span class="label-text">Email *</span>
            </label>
            <input
              id="email"
              type="email"
              formControlName="email"
              placeholder="member@example.com"
              class="input input-bordered w-full"
              [class.input-error]="isFieldInvalid('email')"
            />
            @if (isFieldInvalid('email')) {
              <span class="label-text-alt text-error">Valid email is required</span>
            }
          </div>

          <!-- First Name -->
          <div class="form-control">
            <label class="label" for="firstName">
              <span class="label-text">First Name</span>
            </label>
            <input
              id="firstName"
              type="text"
              formControlName="firstName"
              placeholder="John"
              class="input input-bordered w-full"
            />
          </div>

          <!-- Last Name -->
          <div class="form-control">
            <label class="label" for="lastName">
              <span class="label-text">Last Name</span>
            </label>
            <input
              id="lastName"
              type="text"
              formControlName="lastName"
              placeholder="Doe"
              class="input input-bordered w-full"
            />
          </div>

          <!-- Role -->
          <div class="form-control">
            <label class="label" for="roleId">
              <span class="label-text">Role *</span>
            </label>
            <select
              id="roleId"
              formControlName="roleId"
              class="select select-bordered w-full"
              [class.select-error]="isFieldInvalid('roleId')"
            >
              <option [ngValue]="null" disabled>Select a role</option>
              @for (role of roles; track role.id) {
                <option [ngValue]="role.id">{{ role.name }} - {{ role.description }}</option>
              }
            </select>
            @if (isFieldInvalid('roleId')) {
              <span class="label-text-alt text-error">Role is required</span>
            }
          </div>

          <!-- Error Message -->
          @if (errorMessage) {
            <div class="alert alert-error">
              <span>{{ errorMessage }}</span>
            </div>
          }

          <!-- Actions -->
          <div class="modal-action">
            <button type="button" class="btn btn-ghost" (click)="closeModal.emit()">Cancel</button>
            <button
              type="submit"
              [disabled]="inviteForm.invalid || isSubmitting"
              class="btn btn-primary"
            >
              @if (isSubmitting) {
                <span class="loading loading-spinner loading-sm"></span>
              }
              Send Invitation
            </button>
          </div>
        </form>
      </div>
      <div
        class="modal-backdrop"
        (click)="closeModal.emit()"
        (keyup.escape)="closeModal.emit()"
        tabindex="0"
        role="button"
        aria-label="Close modal"
      ></div>
    </div>
  `,
})
export class TeamInviteComponent {
  @Output() closeModal = new EventEmitter<void>();
  @Output() invited = new EventEmitter<void>();

  private fb = inject(FormBuilder);
  private teamService = inject(TeamService);

  inviteForm: FormGroup;
  isSubmitting = false;
  errorMessage = '';
  roles = this.teamService.getRoles();

  constructor() {
    this.inviteForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      firstName: [''],
      lastName: [''],
      roleId: [null, Validators.required],
    });
  }

  isFieldInvalid(fieldName: string): boolean {
    const field = this.inviteForm.get(fieldName);
    return !!(field?.invalid && (field?.dirty || field?.touched));
  }

  onSubmit(): void {
    if (this.inviteForm.invalid || this.isSubmitting) {
      return;
    }

    this.isSubmitting = true;
    this.errorMessage = '';

    this.teamService.inviteMember(this.inviteForm.value).subscribe({
      next: () => {
        this.isSubmitting = false;
        this.invited.emit();
      },
      error: (error) => {
        this.isSubmitting = false;
        this.errorMessage = error.error?.message || 'Failed to send invitation';
      },
    });
  }
}
