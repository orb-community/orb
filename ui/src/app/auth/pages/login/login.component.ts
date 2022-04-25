import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, OnInit } from '@angular/core';
import { NB_AUTH_OPTIONS, NbAuthService, NbLoginComponent } from '@nebular/auth';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { environment } from '../../../../environments/environment';

@Component({
  selector: 'ngx-orb-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent extends NbLoginComponent implements OnInit {
  strings = STRINGS.login;

  showPassword = false;

  isOverride = false;

  isMaintenance = false;

  dueTime;

  constructor(
    @Inject(NB_AUTH_OPTIONS) protected options: {},
    protected authService: NbAuthService,
    protected cd: ChangeDetectorRef,
    protected router: Router,
    protected route: ActivatedRoute,
    ) {
    super(authService, options, cd, router);

    const tsDue = environment.MAINTENANCE;

    this.dueTime = parseInt(tsDue, 10) * 1000;

    this.isMaintenance = Date.now() < this.dueTime;

    this.isOverride = this.route.snapshot.queryParams.override === '1';
  }

  ngOnInit() { // In the ngOnInit() or in the constructor
    const el = document.getElementById('nb-global-spinner');
    if (el) {
      el.style['display'] = 'none';
    }

    if (this.isMaintenance) {
      this.showMessages.maintenance = true;
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

  loginDisabled() {
    return (this.submitted || this.isMaintenance) && !this.isOverride;
  }
}
