import { Component, ChangeDetectorRef, Inject, OnInit } from '@angular/core';

import { NbAuthService, NB_AUTH_OPTIONS, NbRegisterComponent } from '@nebular/auth';
import { Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-register-component',
  templateUrl: 'register.component.html',
  styleUrls: ['register.component.scss'],
})
export class RegisterComponent extends NbRegisterComponent implements OnInit {
  strings = STRINGS.login;
  showPassword = false;

  constructor(
    @Inject(NB_AUTH_OPTIONS) protected options: {},
    protected authService: NbAuthService,
    protected cd: ChangeDetectorRef,
    protected router: Router,
  ) {
    super(authService, options, cd, router);
  }

  ngOnInit() { // In the ngOnInit() or in the constructor
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

  register() {
    const {email, password, company} = this.user;
    this.authService.register(this.strategy, {
      email,
      password,
      company,
    }).subscribe(
      respReg => {
        this.authService.authenticate(this.strategy, {
          email,
          password,
        }).subscribe(
          respAuth => {
            this.router.navigateByUrl('/pages/dashboard');
          },
        );
      },
    );
  }
}
