import {
  getDeepFromObject,
  NB_AUTH_OPTIONS,
  NbAuthResult,
  NbAuthService,
  NbRequestPasswordComponent,
} from '@nebular/auth';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject } from '@angular/core';
import { Router } from '@angular/router';
import { STRINGS } from '../../../../assets/text/strings';

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

  requestPass(): void {
    this.errors = this.messages = [];
    this.submitted = true;

    this.service.requestPassword(this.strategy, this.user).subscribe((result: NbAuthResult) => {
      this.submitted = false;
      if (result.isSuccess()) {
        this.messages = result.getMessages();
      } else {
        this.errors = result.getErrors();
      }

      const redirect = result.getRedirect();
      if (redirect) {
        setTimeout(() => {
          return this.router.navigateByUrl(redirect);
        }, this.redirectDelay);
      }
      this.cd.detectChanges();
    });
  }

  getConfigValue(key: string): any {
    return getDeepFromObject(this.options, key, null);
  }
}
