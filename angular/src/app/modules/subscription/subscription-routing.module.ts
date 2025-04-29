import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SubscriptionPageComponent } from './subscription-page/subscription-page.component';
import { canActivateAnySubscription } from 'src/app/services/auth/subscription.guard';
import { SubscriptionStatusComponent } from './subscription-status/subscription-status.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'status',
    pathMatch: 'full', // Added pathMatch: 'full' for clarity and best practice
  },
  {
    path: 'status',
    component: SubscriptionStatusComponent,
    // Add guards if necessary, e.g., canActivate: [AuthGuard]
  },
  {
    path: 'select',
    component: SubscriptionPageComponent,
    canActivate: [canActivateAnySubscription],
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SubscriptionRoutingModule { }
