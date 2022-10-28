import { Component, Input, OnInit } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-agent-information',
  templateUrl: './agent-information.component.html',
  styleUrls: ['./agent-information.component.scss'],
})
export class AgentInformationComponent implements OnInit {
  @Input() agent: Agent;

  isResetting: boolean;

  constructor(
    protected agentsService: AgentsService,
    protected notificationService: NotificationsService,
  ) {
    this.isResetting = false;
  }

  ngOnInit(): void {}

  resetAgent() {
    if (!this.isResetting) {
      this.isResetting = true;
      this.agentsService.resetAgent(this.agent.id).subscribe(() => {
        this.isResetting = false;
        this.notifyResetSuccess();
      });
    }
  }

  getAgentVersion() {
    const agentVersion = this.agent?.agent_metadata?.orb_agent?.version;

    return agentVersion ? agentVersion : '-';
  }

  getAgentBackend() {
    const backends = this.agent?.agent_metadata?.backends;
    const backend = !!backends && Object.keys(backends).length > 0 ? Object.keys(backends)[0] : '-';
    return backend;
  }

  getAgentBackendState() {
    const backends_states = this.agent?.last_hb_data?.backend_state;
    let formatted = Object.entries(backends_states).map(([key, value]) => {
      return { backend: key, state: value };
    });
    return backends_states;
  }

  getAgentBackendVersion() {
    const backends = this.agent?.agent_metadata?.backends;
    const version = !!backends && Object.keys(backends).length > 0 ? Object.values(backends)[0]['version'] : '-';
    return version;
  }

  notifyResetSuccess() {
    this.notificationService.success('Agent Reset Requested', '');
  }
}
