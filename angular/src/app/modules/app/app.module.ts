import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BaseComponent } from './base/base.component';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { DashboardComponent } from './dashboard/dashboard.component';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatCardModule } from '@angular/material/card';
import { MatMenuModule } from '@angular/material/menu';
import { RouterModule, Routes } from '@angular/router';
import { CustomSidenavComponent } from './custom-sidenav/custom-sidenav.component';

import { RolePipe } from 'src/app/shared/pipes/role.pipe';

export const routes: Routes = [
  {
    path: '',
    component: BaseComponent,
    children: [
      {
        path: '',
        redirectTo: 'home',
        pathMatch: 'full',
        title: 'Root redirect to /app/home',
      },
      {
        path: 'home',
        component: DashboardComponent,
      },
    ],
  },
];


@NgModule({
  declarations: [
    BaseComponent,
    CustomSidenavComponent,
    DashboardComponent
  ],
  imports: [
    CommonModule,

    RouterModule.forChild(routes),

    MatToolbarModule,
    MatButtonModule,
    MatSidenavModule,
    MatIconModule,
    MatListModule,
    MatGridListModule,
    MatCardModule,
    MatMenuModule,

    RolePipe,
  ]
})
export class AppModule { }
