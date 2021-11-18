import { Component } from '@angular/core';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})

export class DashboardComponent {
  title: string = STRINGS.home.title;
  message: string = STRINGS.home.step.message;
  agent: string = STRINGS.home.step.agent;
  agent_groups: string = STRINGS.home.step.agent_groups;
  policy: string = STRINGS.home.step.policy;
  sink: string = STRINGS.home.step.sink;
  dataset: string = STRINGS.home.step.dataset;

  constructor() {
  }
}
