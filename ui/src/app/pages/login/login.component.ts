import { ChangeDetectionStrategy, Component } from '@angular/core';
import { NbLoginComponent } from '@nebular/auth';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-orb-login',
  templateUrl: './login.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent extends NbLoginComponent {
  strings = STRINGS.login;
}
