import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { TeamListComponent } from './list/team-list.component';
import { TeamInviteComponent } from './invite/team-invite.component';
import { TeamMemberEditComponent } from './edit/team-member-edit.component';

const routes: Routes = [
  {
    path: '',
    component: TeamListComponent,
  },
];

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    RouterModule.forChild(routes),
    TeamListComponent,
    TeamInviteComponent,
    TeamMemberEditComponent,
  ],
  exports: [TeamListComponent, TeamInviteComponent, TeamMemberEditComponent],
})
export class TeamModule {}
