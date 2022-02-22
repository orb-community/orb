declare var _ps;

import { ChangeDetectorRef, Component, Inject, OnInit } from '@angular/core';

import { NB_AUTH_OPTIONS, NbAuthService, NbRegisterComponent } from '@nebular/auth';
import { Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'ngx-register-component',
  templateUrl: 'register.component.html',
  styleUrls: ['register.component.scss'],
})
export class RegisterComponent extends NbRegisterComponent implements OnInit {
  strings = STRINGS.login;

  _isProduction = environment.production;

  /**
   * Pactsafe
   */
  _sid = environment.PS.SID;

  _groupKey = environment.PS.GROUP_KEY;

  _psEnabled = !!this._sid && !!this._groupKey;

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
    _ps = window['_ps'];
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
    if (!this._psEnabled) return false;
    // Check to ensure we're able to get the Group successfully.
    if (_ps.getByKey(this._groupKey)) {

      // Return if we should block the submission using the .block() method
      // provided by the Group object.
      return _ps.getByKey(this._groupKey).block();
    } else {
      // We weren't able to get the group,
      // so blocking form submission may be needed.
      return true;
    }
  }

  register(event?: any) {
    // Prevent the form from automatically submitting without
    // checking PactSafe acceptance first.
    this.errors = this.messages = [];
    this.submitted = true;

    event?.preventDefault();
    if (!this.blockSubmission()) {
      // We don't need to block the form submission,
      // so submit the form.
      const { email, password, company } = this.user;
      this.authService.register(this.strategy, {
        email,
        password,
        company,
      }).subscribe(
        respReg => {
          const first_name = this.user.fullname.split(' ')[0];
          const last_name = this.user.fullName.replace(`${first_name} `, '');
          _ps.getByKey(this._groupKey).send('updated', { custom_data: { first_name, last_name } });

          this.submitted = false;

          if (respReg.isSuccess()) {
            this.messages = respReg.getMessages();
          } else {
            this.errors = respReg.getErrors();
          }

          this.cd.detectChanges();

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
        (_ps.getByKey(this._groupKey) && _ps.getByKey(this._groupKey).get('alert_message')) ?
          _ps.getByKey(this._groupKey).get('alert_message') :
          'Please accept our Terms and Conditions.';

      // Alert the user that the Terms need to be accepted before continuing.
      alert(acceptanceAlertLanguage);
    }
  }
}
