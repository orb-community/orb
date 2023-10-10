import { ChangeDetectorRef, Component, Inject, OnInit } from '@angular/core';

import { Router } from '@angular/router';
import {
  NbAuthService,
  NbRegisterComponent,
  NB_AUTH_OPTIONS,
} from '@nebular/auth';
import { STRINGS } from 'assets/text/strings';
import { environment } from '../../../../environments/environment';

@Component({
  selector: 'ngx-register-component',
  templateUrl: 'register.component.html',
  styleUrls: ['register.component.scss'],
})
export class RegisterComponent extends NbRegisterComponent implements OnInit {
  strings = STRINGS.login;

  _isProduction = environment.production;

  // TODO
  orbErrors = {};

  showPassword = false;

  repeatedEmail = null;

  constructor(
    @Inject(NB_AUTH_OPTIONS) protected options: {},
    protected authService: NbAuthService,
    protected cd: ChangeDetectorRef,
    protected router: Router,
  ) {
    super(authService, options, cd, router);
  }

  ngOnInit() {
    // In the ngOnInit() or in the constructor
    const el = document.getElementById('nb-global-spinner');
    if (el) {
      el.style['display'] = 'none';
    }
  }

  getInputType() {
    if (this.showPassword) {
      return 'text';
    }
    return 'password';
  }

  toggleShowPassword() {
    this.showPassword = !this.showPassword;
  }

  register(event?: any) {
    this.orbErrors = {};
    this.errors = this.messages = [];
    this.submitted = true;
    this.repeatedEmail = null;
    
    const { email, password, company } = this.user;
    this.authService
      .register(this.strategy, {
        email,
        password,
        metadata: {
          company: company,
          fullName: this.user.fullName,
        },
      })
      .subscribe((respReg) => {
        const first_name = this.user.fullName.split(' ')[0];
        const last_name = this.user.fullName.replace(`${first_name} `, '');

        this.submitted = false;

        if (respReg.isSuccess()) {
          this.messages = respReg.getMessages();

          this.authenticateAndRedirect(email, password);
        } else {
          if (respReg.getResponse().status === 409) {
            this.repeatedEmail = email;
            this.errors = [respReg.getResponse().error.error];
          }
        }
      });
    
  }

  authenticateAndRedirect(email, password) {
    this.authService
      .authenticate(this.strategy, {
        email,
        password,
      })
      .subscribe((respAuth) => {
        this.router.navigateByUrl('/pages/dashboard');
      });
  }
}
