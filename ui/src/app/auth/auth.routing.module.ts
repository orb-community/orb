import { RouterModule, Routes } from '@angular/router';

import { LoginComponent } from 'app/pages/login/login.component';
import { LogoutComponent } from 'app/pages/logout/logout.component';
import { RegisterComponent } from 'app/pages/register/register.component';
import { NgModule } from '@angular/core';
import { RequestPasswordComponent } from 'app/pages/request-password/request-password.component';
import { ResetPasswordComponent } from 'app/pages/reset-password/reset-password.component';

export const routes: Routes = [
  {
    path: '',
    component: LoginComponent,
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: 'register',
    component: RegisterComponent,
  },
  {
    path: 'logout',
    component: LogoutComponent,
  },
  {
    path: 'request-password',
    component: RequestPasswordComponent,
  },
  {
    path: 'reset-password',
    component: ResetPasswordComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class AuthRoutingModule {
}
