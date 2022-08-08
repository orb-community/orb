import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { OrbService } from 'app/common/services/orb.service';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';

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

  datasets: { [id: string]: Dataset };

  groups: AgentGroup[];

  agentID;

  agentSubscription: Subscription;

  constructor(
    protected agentsService: AgentsService,
    protected route: ActivatedRoute,
    protected router: Router,
    protected orb: OrbService,
  ) {
    this.agent = {};
    this.datasets = {};
    this.isLoading = true;
  }

  ngOnInit() {
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.agentSubscription = this.orb.getAgentFullView(this.agentID).subscribe({
      next: ({ agent, datasets, groups }) => {
        this.agent = agent;
        this.datasets = datasets as {[id: string]: Dataset};
        this.groups = groups;
        this.isLoading = false;
      },
      error: (err) => {
        this.isLoading = false;
      },
    });
    this.isLoading = true;
    this.retrieveAgent();
  }

  retrieveAgent() {
    return this.agentsService.getAgentById(this.agentID).subscribe((agent) => {
      this.agent = agent;
      this.isLoading = false;
    });
  }

  isToday() {
    const today = new Date(Date.now());
    const date = new Date(this?.agent?.ts_last_hb);

    return (
      today.getDay() === date.getDay() &&
      today.getMonth() === date.getMonth() &&
      today.getFullYear() === date.getFullYear()
    );
  }

  ngOnDestroy() {
    this.agentSubscription?.unsubscribe();
  }
}
