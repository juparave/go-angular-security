import { Component, EventEmitter, Input, OnChanges, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { TeamService, TeamMember } from '../../../services/team.service';

@Component({
  selector: 'app-team-member-edit',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <div class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Edit Team Member</h3>

        <form [formGroup]="editForm" (ngSubmit)="onSubmit()" class="space-y-4">
          <!-- Member Info (read-only) -->
          <div class="bg-base-200 p-4 rounded-lg mb-4">
            <p class="font-medium">{{ member.firstName }} {{ member.lastName }}</p>
            <p class="text-sm text-gray-500">{{ member.email }}</p>
          </div>

          <!-- Role -->
          <div class="form-control">
            <label class="label" for="roleId">
              <span class="label-text">Role</span>
            </label>
            <select id="roleId" formControlName="roleId" class="select select-bordered w-full">
              @for (role of roles; track role.id) {
                <option [ngValue]="role.id">{{ role.name }} - {{ role.description }}</option>
              }
            </select>
          </div>

          <!-- Status -->
          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-4">
              <input type="checkbox" formControlName="enabled" class="toggle toggle-primary" />
              <span class="label-text">Account Enabled</span>
            </label>
            <p class="text-sm text-gray-500 mt-1">Disabled users cannot log in to the account</p>
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
              [disabled]="editForm.invalid || isSubmitting"
              class="btn btn-primary"
            >
              @if (isSubmitting) {
                <span class="loading loading-spinner loading-sm"></span>
              }
              Save Changes
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
export class TeamMemberEditComponent implements OnChanges {
  @Input() member!: TeamMember;
  @Output() closeModal = new EventEmitter<void>();
  @Output() updated = new EventEmitter<void>();

  private fb = inject(FormBuilder);
  private teamService = inject(TeamService);

  editForm: FormGroup;
  isSubmitting = false;
  errorMessage = '';
  roles = this.teamService.getRoles();

  constructor() {
    this.editForm = this.fb.group({
      roleId: [null, Validators.required],
      enabled: [true],
    });
  }

  ngOnChanges(): void {
    if (this.member) {
      this.editForm.patchValue({
        roleId: this.member.roleId,
        enabled: this.member.enabled,
      });
    }
  }

  onSubmit(): void {
    if (this.editForm.invalid || this.isSubmitting) {
      return;
    }

    this.isSubmitting = true;
    this.errorMessage = '';

    this.teamService.updateMember(this.member.id, this.editForm.value).subscribe({
      next: () => {
        this.isSubmitting = false;
        this.updated.emit();
      },
      error: (error) => {
        this.isSubmitting = false;
        this.errorMessage = error.error?.message || 'Failed to update member';
      },
    });
  }
}
