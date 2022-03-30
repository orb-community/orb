import { Component } from '@angular/core';
import { map } from 'rxjs/operators';
import { Breakpoints, BreakpointObserver } from '@angular/cdk/layout';
import { STRINGS } from '../../../assets/text/strings';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent {
  strings = { dashboard: STRINGS.dashboard };
  /** Based on the screen size, switch from standard to one column per row */
  cards = this.breakpointObserver.observe(Breakpoints.Handset).pipe(
    map(({ matches }) => {
      const { step } = this.strings.dashboard;
      if (matches) {
        return [
          { title: '', content: step.agent, cols: 2, rows: 1 },
          { title: '', content: step.agent_groups, cols: 2, rows: 1 },
          { title: '', content: step.sink, cols: 2, rows: 1 },
          { title: '', content: step.policy, cols: 2, rows: 1 },
          { title: '', content: step.dataset, cols: 2, rows: 1 },
        ];
      }

      return [
        { title: '', content: step.agent, cols: 1, rows: 1 },
        { title: '', content: step.agent_groups, cols: 1, rows: 1 },
        { title: '', content: step.sink, cols: 1, rows: 1 },
        { title: '', content: step.policy, cols: 1, rows: 1 },
        { title: '', content: step.dataset, cols: 1, rows: 1 },
      ];
    })
  );

  constructor(private breakpointObserver: BreakpointObserver) {}
}
