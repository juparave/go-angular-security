import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TeamService, TeamMember } from '../../../services/team.service';

@Component({
  selector: 'app-team-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container mx-auto px-4 py-8">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold">Team Members</h1>
        <button class="btn btn-primary" (click)="showInviteModal = true">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5 mr-2"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
              clip-rule="evenodd"
            />
          </svg>
          Invite Member
        </button>
      </div>

      @if (loading) {
        <div class="flex justify-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      } @else if (members.length === 0) {
        <div class="text-center py-8">
          <p class="text-gray-500">No team members yet. Invite your first team member!</p>
        </div>
      } @else {
        <div class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Last Login</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (member of members; track member.id) {
                <tr>
                  <td>{{ member.firstName }} {{ member.lastName }}</td>
                  <td>{{ member.email }}</td>
                  <td>
                    <span [class]="'badge ' + getRoleBadgeClass(member.role)">
                      {{ member.role }}
                    </span>
                  </td>
                  <td>
                    @if (member.enabled) {
                      <span class="badge badge-success">Active</span>
                    } @else {
                      <span class="badge badge-error">Disabled</span>
                    }
                  </td>
                  <td>
                    {{ member.lastLoginAt ? (member.lastLoginAt | date: 'medium') : 'Never' }}
                  </td>
                  <td>
                    <div class="flex gap-2">
                      <button
                        class="btn btn-sm btn-ghost"
                        (click)="editMember(member)"
                        title="Edit"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          viewBox="0 0 20 20"
                          fill="currentColor"
                        >
                          <path
                            d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
                          />
                        </svg>
                      </button>
                      <button
                        class="btn btn-sm btn-ghost text-error"
                        (click)="confirmDelete(member)"
                        title="Remove"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          viewBox="0 0 20 20"
                          fill="currentColor"
                        >
                          <path
                            fill-rule="evenodd"
                            d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                            clip-rule="evenodd"
                          />
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      }

      <!-- Invite Modal -->
      @if (showInviteModal) {
        <app-team-invite (closeModal)="showInviteModal = false" (invited)="onMemberInvited()" />
      }

      <!-- Edit Modal -->
      @if (showEditModal && selectedMember) {
        <app-team-member-edit
          [member]="selectedMember"
          (closeModal)="closeEditModal()"
          (updated)="onMemberUpdated()"
        />
      }

      <!-- Delete Confirmation Modal -->
      @if (showDeleteModal && memberToDelete) {
        <div class="modal modal-open">
          <div class="modal-box">
            <h3 class="font-bold text-lg">Remove Team Member</h3>
            <p class="py-4">
              Are you sure you want to remove <strong>{{ memberToDelete.email }}</strong> from the
              team?
            </p>
            <div class="modal-action">
              <button class="btn btn-ghost" (click)="showDeleteModal = false">Cancel</button>
              <button class="btn btn-error" (click)="deleteMember()">Remove</button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class TeamListComponent implements OnInit {
  members: TeamMember[] = [];
  loading = true;
  showInviteModal = false;
  showEditModal = false;
  showDeleteModal = false;
  selectedMember: TeamMember | null = null;
  memberToDelete: TeamMember | null = null;

  private teamService = inject(TeamService);

  constructor() {}

  ngOnInit(): void {
    this.loadMembers();
  }

  loadMembers(): void {
    this.loading = true;
    this.teamService.getMembers().subscribe({
      next: (members) => {
        this.members = members;
        this.loading = false;
      },
      error: (error) => {
        console.error('Failed to load team members:', error);
        this.loading = false;
      },
    });
  }

  getRoleBadgeClass(role: string): string {
    switch (role.toLowerCase()) {
      case 'admin':
        return 'badge-primary';
      case 'editor':
        return 'badge-secondary';
      default:
        return 'badge-ghost';
    }
  }

  editMember(member: TeamMember): void {
    this.selectedMember = member;
    this.showEditModal = true;
  }

  closeEditModal(): void {
    this.showEditModal = false;
    this.selectedMember = null;
  }

  onMemberUpdated(): void {
    this.closeEditModal();
    this.loadMembers();
  }

  onMemberInvited(): void {
    this.showInviteModal = false;
    this.loadMembers();
  }

  confirmDelete(member: TeamMember): void {
    this.memberToDelete = member;
    this.showDeleteModal = true;
  }

  deleteMember(): void {
    if (!this.memberToDelete) return;

    this.teamService.deleteMember(this.memberToDelete.id).subscribe({
      next: () => {
        this.showDeleteModal = false;
        this.memberToDelete = null;
        this.loadMembers();
      },
      error: (error) => {
        console.error('Failed to delete member:', error);
      },
    });
  }
}
