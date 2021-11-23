import { RouterModule, Routes } from '@angular/router';

import { LoginComponent } from 'app/pages/login/login.component';
import { LogoutComponent } from 'app/pages/logout/logout.component';
import { RegisterComponent } from 'app/pages/register/register.component';
import { NbResetPasswordComponent } from '@nebular/auth';
import { NgModule } from '@angular/core';
import { RequestPasswordComponent } from 'app/pages/login/request-password.component';

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
    component: NbResetPasswordComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class AuthRoutingModule {
}
