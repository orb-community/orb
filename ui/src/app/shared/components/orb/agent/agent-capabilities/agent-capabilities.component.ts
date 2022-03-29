import { Component, Input, OnInit } from '@angular/core';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';

@Component({
  selector: 'ngx-agent-capabilities',
  templateUrl: './agent-capabilities.component.html',
  styleUrls: ['./agent-capabilities.component.scss'],
})
export class AgentCapabilitiesComponent implements OnInit {
  @Input() agent: Agent;

  agentStates = AgentStates;

  constructor() {
  }

  ngOnInit(): void {
  }

}
