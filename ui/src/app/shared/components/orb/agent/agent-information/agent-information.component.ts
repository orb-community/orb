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

  getAgentBackend() {
    return Object.keys(this.agent.agent_metadata.backends)[0] || '-';
  }

  getAgentBackendVersion() {
    const backend = Object.keys(this.agent.agent_metadata.backends)[0];
    return backend ? this.agent.agent_metadata.backends[backend].version : '-';
  }

  notifyResetSuccess() {
    this.notificationService.success('Agent Reset Requested', '');
  }
}
