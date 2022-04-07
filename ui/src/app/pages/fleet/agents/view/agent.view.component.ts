import { Component, OnDestroy, OnInit } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { Subscription } from 'rxjs';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.view.component.html',
  styleUrls: ['./agent.view.component.scss'],
})
export class AgentViewComponent implements OnInit, OnDestroy {
  strings = STRINGS.agents;

  agentStates = AgentStates;

  isLoading: boolean;

  agent: Agent;

  agentID;

  agentSubscription: Subscription;

  isResetting: boolean;

  constructor(
    protected agentsService: AgentsService,
    protected route: ActivatedRoute,
    protected router: Router,
    protected notificationService: NotificationsService,
  ) {
    this.agent = {};
    this.isLoading = false;
    this.isResetting = false;
  }

  ngOnInit() {
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.retrieveAgent();
  }

  retrieveAgent() {
    this.isLoading = true;
    return this.agentsService
      .getAgentById(this.agentID)
      .subscribe((agent) => {
        this.agent = agent;
        this.isLoading = false;
      });
  }

  resetAgent() {
    if (!this.isResetting) {
      this.isResetting = true;
      this.agentsService.resetAgent(this.agentID).subscribe(() => {
        this.isResetting = false;
        this.notifyResetSuccess();
      });
    }
  }

  notifyResetSuccess() {
    this.notificationService.success('Agent Reset Requested', '');
  }

  isToday() {
    const today = new Date(Date.now());
    const date = new Date(this?.agent?.ts_last_hb);

    return today.getDay() === date.getDay()
      && today.getMonth() === date.getMonth()
      && today.getFullYear() === date.getFullYear();

  }

  ngOnDestroy() {
    this.agentSubscription?.unsubscribe();
  }


}
