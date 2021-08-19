import { Component } from '@angular/core';
import { NbAuthComponent } from '@nebular/auth';

@Component({
  selector: 'ngx-orb-auth',
  styleUrls: ['./auth.component.scss'],
  template: `
    <nb-layout>
      <nb-layout-column>

        <nav class="navigation">
          <a href="#" (click)="back()" class="link back-link" aria-label="Back">
            <nb-icon icon="arrow-back"></nb-icon>
          </a>
        </nav>

        <nb-auth-block>
          <router-outlet></router-outlet>
        </nb-auth-block>

      </nb-layout-column>
    </nb-layout>
  `,
})
export class AuthComponent extends NbAuthComponent {
}
