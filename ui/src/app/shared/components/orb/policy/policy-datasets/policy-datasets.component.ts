import { Component, Input, OnInit } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-policy-datasets',
  templateUrl: './policy-datasets.component.html',
  styleUrls: ['./policy-datasets.component.scss'],
})
export class PolicyDatasetsComponent implements OnInit {
  @Input()
  policy: AgentPolicy;

  constructor() { }

  ngOnInit(): void {
  }

}
