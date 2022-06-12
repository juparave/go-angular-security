import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AppComponent } from './app.component';
import { RouterModule } from '@angular/router';

export const routes = [
  {
    path: '',
    component: AppComponent,
  },
];

@NgModule({
  declarations: [AppComponent],
  imports: [CommonModule, RouterModule.forChild(routes)],
})
export class AppModule {}
