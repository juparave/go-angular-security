import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from 'src/app/modules/login/login.component';
import { ResetPasswordComponent } from 'src/app/modules/login/reset-password/reset-password.component';
import { RequestResetPasswordComponent } from 'src/app/modules/login/request-reset-password/request-reset-password.component';
import { LoginFormComponent } from 'src/app/modules/login/form/login-form.component';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from 'src/app/shared/shared.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SignUpComponent } from './sign-up/sign-up.component';
import { BackendErrorMessagesModule } from 'src/app/shared/modules/backend-error-messages/backend-error-messages.module';

export const routes: Routes = [
  {
    path: '',
    component: LoginComponent,
    children: [
      { path: '', redirectTo: 'form', pathMatch: 'full' },
      { path: 'form', component: LoginFormComponent },
      {
        path: 'request-reset-password',
        component: RequestResetPasswordComponent,
      },
      {
        path: 'reset-password/:token/:email',
        component: ResetPasswordComponent,
      },
      {
        path: 'sign-up',
        component: SignUpComponent,
      },
    ],
  },
];

@NgModule({
  declarations: [
    LoginComponent,
    ResetPasswordComponent,
    RequestResetPasswordComponent,
    LoginFormComponent,
    SignUpComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,

    BackendErrorMessagesModule,

    SharedModule,
    RouterModule.forChild(routes),
  ],
})
export class LoginModule {}
