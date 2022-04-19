import { Component, Input, OnInit } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-policy-details',
  templateUrl: './policy-details.component.html',
  styleUrls: ['./policy-details.component.scss'],
})
export class PolicyDetailsComponent implements OnInit {
  @Input()
  policy: AgentPolicy;

  constructor() {
  }

  ngOnInit(): void {
  }

}
