import { Component } from '@angular/core';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})

export class DashboardComponent {
  title: string = STRINGS.home.title;
  description: string = STRINGS.home.description;

  constructor(
  ) { }
}
