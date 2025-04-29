import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SubscriptionPageComponent } from './subscription-page/subscription-page.component';
import { SubscriptionStatusComponent } from './subscription-status/subscription-status.component';
import { SubscriptionSuccessComponent } from './subscription-success/subscription-success.component'; // Added import
import { canActivateNoActiveSubscription } from 'src/app/services/auth/no-active-subscription.guard'; // Keep this import
import { AuthGuard } from 'src/app/services/auth/auth.guard';
import { BaseComponent } from '../app/base/base.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'status',
    pathMatch: 'full',
  },
  {
    path: 'status',
    component: BaseComponent,
    children: [
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
    ],
  },
  {
    path: 'select',
    component: SubscriptionPageComponent,
    canActivate: [canActivateNoActiveSubscription],
  },
  { // Added route for success page
    path: 'success',
    component: SubscriptionSuccessComponent,
    canActivate: [AuthGuard]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SubscriptionRoutingModule { }
