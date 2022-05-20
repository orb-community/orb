import { Injector, ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  NbAlertModule,
  NbButtonModule,
  NbCheckboxModule,
  NbIconModule,
  NbInputModule,
  NbLayoutModule,
} from '@nebular/theme';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AuthComponent } from 'app/auth/components/auth.component';
import { LoginComponent } from 'app/auth/pages/login/login.component';
import {
  NB_AUTH_FALLBACK_TOKEN,
  NB_AUTH_INTERCEPTOR_HEADER,
  NB_AUTH_OPTIONS,
  NB_AUTH_STRATEGIES,
  NB_AUTH_TOKEN_INTERCEPTOR_FILTER,
  NB_AUTH_TOKENS,
  NB_AUTH_USER_OPTIONS,
  NbAuthModule,
  NbAuthOptions,
  NbAuthService,
  NbAuthSimpleToken,
  NbAuthTokenParceler,
  NbDummyAuthStrategy,
  nbNoOpInterceptorFilter,
  NbOAuth2AuthStrategy,
  nbOptionsFactory,
  NbPasswordAuthStrategy,
  nbStrategiesFactory,
  NbTokenLocalStorage,
  NbTokenService,
  nbTokensFactory,
  NbTokenStorage,
} from '@nebular/auth';
import { AuthRoutingModule } from 'app/auth/auth.routing.module';
import { RegisterComponent } from 'app/auth/pages/register/register.component';
import { LogoutComponent } from 'app/auth/pages/logout/logout.component';
import { RequestPasswordComponent } from 'app/auth/pages/request-password/request-password.component';
import { ResetPasswordComponent } from 'app/auth/pages/reset-password/reset-password.component';
import { PSModule } from '@pactsafe/pactsafe-angular-sdk';

@NgModule({
  imports: [
    PSModule.forRoot(),
    CommonModule,
    NbLayoutModule,
    NbCheckboxModule,
    NbAlertModule,
    NbInputModule,
    NbButtonModule,
    RouterModule,
    FormsModule,
    NbIconModule,
    NbAuthModule,
  ],
  declarations: [
    AuthComponent,
    LoginComponent,
    RequestPasswordComponent,
    ResetPasswordComponent,
    RegisterComponent,
    LogoutComponent,
  ],
  exports: [
    AuthRoutingModule,
    AuthComponent,
    LoginComponent,
    RequestPasswordComponent,
    ResetPasswordComponent,
    RegisterComponent,
    LogoutComponent,
  ],
})
export class AuthModule extends NbAuthModule {
  static forRoot(nbAuthOptions?: NbAuthOptions): ModuleWithProviders<AuthModule> {
    return {
      ngModule: AuthModule,
      providers: [
        {provide: NB_AUTH_USER_OPTIONS, useValue: nbAuthOptions},
        {provide: NB_AUTH_OPTIONS, useFactory: nbOptionsFactory, deps: [NB_AUTH_USER_OPTIONS]},
        {provide: NB_AUTH_STRATEGIES, useFactory: nbStrategiesFactory, deps: [NB_AUTH_OPTIONS, Injector]},
        {provide: NB_AUTH_TOKENS, useFactory: nbTokensFactory, deps: [NB_AUTH_STRATEGIES]},
        {provide: NB_AUTH_FALLBACK_TOKEN, useValue: NbAuthSimpleToken},
        {provide: NB_AUTH_INTERCEPTOR_HEADER, useValue: 'Authorization'},
        {provide: NB_AUTH_TOKEN_INTERCEPTOR_FILTER, useValue: nbNoOpInterceptorFilter},
        {provide: NbTokenStorage, useClass: NbTokenLocalStorage},
        NbAuthTokenParceler,
        NbAuthService,
        NbTokenService,
        NbDummyAuthStrategy,
        NbPasswordAuthStrategy,
        NbOAuth2AuthStrategy,
      ],
    };
  }
}
