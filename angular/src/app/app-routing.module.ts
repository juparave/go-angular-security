import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UnAuthGuard } from 'src/app/services/auth/unauth.guard';
import { AuthGuard } from './services/auth/auth.guard';

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
      import('./modules/myapp/myapp.module').then((m) => m.MyAppModule),
    canActivate: [AuthGuard],
    title: 'My App',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule { }
