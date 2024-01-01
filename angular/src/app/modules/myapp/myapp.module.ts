import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from 'src/app/shared/shared.module';
import { TopBarModule } from 'src/app/shared/modules/top-bar/top-bar.module';

import { MyAppComponent } from './myapp.component';
import { HomeComponent } from './home/home.component';

export const routes: Routes = [
  {
    path: '',
    component: MyAppComponent,
    children: [
      {
        path: '',
        redirectTo: 'home',
        pathMatch: 'full',
        title: 'Root redirect to /app/home',
      },
      {
        path: 'home',
        component: HomeComponent,
      },
    ],
  },
];

@NgModule({
  declarations: [MyAppComponent, HomeComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    SharedModule,
    TopBarModule,
  ],
})
export class MyAppModule { }
