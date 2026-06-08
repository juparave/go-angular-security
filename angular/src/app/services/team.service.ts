import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

export interface TeamMember {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  role: string;
  roleId: number;
  enabled: boolean;
  passwordChangeRequired: boolean;
  lastLoginAt?: string;
}

export interface InviteMemberRequest {
  email: string;
  firstName?: string;
  lastName?: string;
  roleId: number;
}

export interface UpdateMemberRequest {
  roleId?: number;
  enabled?: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class TeamService {
  private baseUrl = `${environment.apiUrl}/team`;
  private http = inject(HttpClient);

  constructor() {}

  /**
   * Get all team members for the current account
   */
  getMembers(): Observable<TeamMember[]> {
    return this.http.get<TeamMember[]>(this.baseUrl);
  }

  /**
   * Invite a new team member
   */
  inviteMember(data: InviteMemberRequest): Observable<TeamMember> {
    return this.http.post<TeamMember>(this.baseUrl, data);
  }

  /**
   * Update a team member's role or status
   */
  updateMember(memberId: string, data: UpdateMemberRequest): Observable<{ message: string }> {
    return this.http.put<{ message: string }>(`${this.baseUrl}/${memberId}`, data);
  }

  /**
   * Remove a team member
   */
  deleteMember(memberId: string): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.baseUrl}/${memberId}`);
  }

  /**
   * Resend invitation to a team member
   */
  resendInvitation(memberId: string): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.baseUrl}/${memberId}/resend`, {});
  }

  /**
   * Get available roles
   */
  getRoles(): { id: number; name: string; description: string }[] {
    return [
      { id: 1, name: 'Admin', description: 'Full access to all features' },
      { id: 2, name: 'Editor', description: 'Can create, edit, and view content' },
      { id: 3, name: 'Viewer', description: 'Can only view content' },
    ];
  }
}
