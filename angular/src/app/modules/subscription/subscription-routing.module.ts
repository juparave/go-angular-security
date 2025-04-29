import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SubscriptionPageComponent } from './subscription-page/subscription-page.component';
import { SubscriptionStatusComponent } from './subscription-status/subscription-status.component';
import { canActivateNoActiveSubscription } from 'src/app/services/auth/no-active-subscription.guard'; // Keep this import
import { AuthGuard } from 'src/app/services/auth/auth.guard';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'status',
    pathMatch: 'full',
  },
  {
    path: 'status',
    component: SubscriptionStatusComponent,
    canActivate: [AuthGuard]
  },
  {
    path: 'select',
    component: SubscriptionPageComponent,
    canActivate: [canActivateNoActiveSubscription],
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SubscriptionRoutingModule { }
