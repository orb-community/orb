import { Component } from '@angular/core';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss'],
})
export class FooterComponent {
    disclaimer: string = STRINGS.footer.disclaimer;
}
