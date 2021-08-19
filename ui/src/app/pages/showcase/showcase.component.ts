import { Component } from '@angular/core';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-showcase-component',
  templateUrl: './showcase.component.html',
  styleUrls: ['./showcase.component.scss'],
})
export class ShowcaseComponent {
  title: string = STRINGS.home.title;

  constructor() {
  }
}
