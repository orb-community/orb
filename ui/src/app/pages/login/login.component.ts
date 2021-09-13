import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { NbLoginComponent } from '@nebular/auth';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-orb-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent extends NbLoginComponent implements OnInit {
  strings = STRINGS.login;
  showPassword = false;

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
}
