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

  redirectDelay: number = 0;
  showMessages: any = {};
  strategy: string = '';

  submitted = false;
  errors: string[] = [];
  messages: string[] = [];
  user: any = {};

  showPassword = false;
  groupOptions = {
    container_selector: 'pactsafe-container',
    display_all: true,
    signer_id_selector: 'input-email',
    test_mode: true,
  };

  constructor(
    @Inject(NB_AUTH_OPTIONS) protected options: {},
    protected authService: NbAuthService,
    protected cd: ChangeDetectorRef,
    protected router: Router,
  ) {
    super(authService, options, cd, router);
    _ps('create', _sid);
    _ps.debug = true;
    _ps('load', _groupKey, this.groupOptions);
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

  // Return whether to block the submission or not.
  blockSubmission() {
    // Check to ensure we're able to get the Group successfully.
    if (_ps.getByKey(_groupKey)) {

      // Return if we should block the submission using the .block() method
      // provided by the Group object.
      return _ps.getByKey(_groupKey).block();
    } else {
      // We weren't able to get the group,
      // so blocking form submission may be needed.
      return true;
    }
  }

  register(event?: any) {
    // Prevent the form from automatically submitting without
    // checking PactSafe acceptance first.

    event?.preventDefault();
    if (!this.blockSubmission()) {
      // We don't need to block the form submission,
      // so submit the form.
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
    } else {
      // We can get the alert message if set on the group
      // or define our own if it's not.
      const acceptanceAlertLanguage =
        (_ps.getByKey(_groupKey) && _ps.getByKey(_groupKey).get('alert_message')) ?
          _ps.getByKey(_groupKey).get('alert_message') :
          'Please accept our Terms and Conditions.';

      // Alert the user that the Terms need to be accepted before continuing.
      alert(acceptanceAlertLanguage);
    }

  }
}
