import { Component } from '@angular/core';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-components-showcase',
  templateUrl: './components-showcase.component.html',
  styleUrls: ['./components-showcase.component.scss'],
})

export class ComponentsShowcaseComponent {
  title: string = STRINGS.home.title;

  constructor() {
  }
}
