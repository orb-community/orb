import { Component, Input, OnInit } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-agent-policies-datasets',
  templateUrl: './agent-policies-datasets.component.html',
  styleUrls: ['./agent-policies-datasets.component.scss'],
})
export class AgentPoliciesDatasetsComponent implements OnInit {
  @Input() agent: Agent;

  @Input() policies: AgentPolicy[];

  constructor() { }

  ngOnInit(): void {
  }

}
