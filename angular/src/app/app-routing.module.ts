import { NgModule } from '@angular/core';
import { RouterModule, Routes, provideRouter, withComponentInputBinding } from '@angular/router';
import { UnAuthGuard } from 'src/app/services/auth/unauth.guard';
import { AuthGuard } from './services/auth/auth.guard';
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
    canActivate: [AuthGuard],
    title: 'Dashboard :: Austral QC',
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
