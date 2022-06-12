import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MyAppComponent } from './myapp.component';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'src/app/shared/shared.module';
import { HomeComponent } from './home/home.component';

export const routes = [
  {
    path: '',
    component: MyAppComponent,
    children: [
      {
        path: '',
        redirectTo: 'home',
        pathMatch: 'full',
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
  imports: [CommonModule, RouterModule.forChild(routes), SharedModule],
})
export class MyAppModule {}
