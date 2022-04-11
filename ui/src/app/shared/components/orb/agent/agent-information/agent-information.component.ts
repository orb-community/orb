import { Component, Input, OnInit } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { AgentsService } from 'app/common/services/agents/agents.service';

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

  ngOnInit(): void {
  }

  resetAgent() {
    if (!this.isResetting) {
      this.isResetting = true;
      this.agentsService.resetAgent(this.agent.id).subscribe(() => {
        this.isResetting = false;
        this.notifyResetSuccess();
      });
    }
  }

  notifyResetSuccess() {
    this.notificationService.success('Agent Reset Requested', '');
  }
}
