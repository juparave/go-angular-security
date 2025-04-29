import { NgModule } from '@angular/core';
import { RouterModule, Routes, provideRouter, withComponentInputBinding } from '@angular/router';
import { UnAuthGuard } from 'src/app/services/auth/unauth.guard';
import { AuthGuard, canActivateAuth } from './services/auth/auth.guard'; // Import canActivateAuth
import { canActivateAnySubscription } from './services/auth/subscription.guard'; // Import subscription guard
import { canActivateNoActiveSubscription } from './services/auth/no-active-subscription.guard'; // Import NoActiveSubscriptionGuard
import { PageNotFoundComponent } from './shared/components/page-not-found/page-not-found.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: '/login',
    pathMatch: 'full',
    title: 'Root redirect to /login',
  },
  {
    path: 'login',
    loadChildren: () =>
      import('./modules/login/login.module').then((m) => m.LoginModule),
    canActivate: [UnAuthGuard],
    title: 'Login',
  },
  {
    path: 'app',
    loadChildren: () =>
      import('./modules/app/app.module').then((m) => m.AppModule),
    canActivate: [canActivateAuth, canActivateAnySubscription],
    title: 'Dashboard :: Austral QC',
  },
  {
    path: 'subscription',
    loadChildren: () =>
      import('./modules/subscription/subscription.module').then((m) => m.SubscriptionModule),
    canActivate: [canActivateAuth],
    title: 'Subscription',
  },
  {
    path: '**',
    component: PageNotFoundComponent,
  }
];

@NgModule({
  exports: [RouterModule],
  providers: [
    provideRouter(routes, withComponentInputBinding()),
  ]
})
export class AppRoutingModule { }
