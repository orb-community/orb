import {
  NB_AUTH_OPTIONS,
  NbAuthService,
  NbRequestPasswordComponent,
} from '@nebular/auth';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject } from '@angular/core';
import { Router } from '@angular/router';
import { STRINGS } from '../../../assets/text/strings';

@Component({
  selector: 'ngx-orb-request-password',
  styleUrls: ['./request-password.component.scss'],
  templateUrl: './request-password.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RequestPasswordComponent extends NbRequestPasswordComponent {
  strings = STRINGS.login;
  showPassword = false;

  constructor(protected service: NbAuthService,
    @Inject(NB_AUTH_OPTIONS) protected options = {},
    protected cd: ChangeDetectorRef,
    protected router: Router) {
    super(service, options, cd, router);
  }
}
