import { Component, OnDestroy, OnInit } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { Subscription } from 'rxjs';
import { OrbService } from 'app/common/services/orb.service';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';

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
    console.log(this.agentID);
    this.agentSubscription = this.orb.getAgentFullView(this.agentID).subscribe({
      next: ({agent, datasets, groups}) => {
        this.agent = agent;
        this.datasets = datasets;
        this.groups = groups;
        console.table(this.agent);
        console.table(this.datasets);
        console.table(this.groups);
        this.isLoading = false;
      }, error: (err) => {
        this.isLoading = false;
      }
    })
    this.isLoading = true;
    this.retrieveAgent();
  }

  retrieveAgent() {
    return this.agentsService
      .getAgentById(this.agentID)
      .subscribe((agent) => {
        this.agent = agent;
        this.isLoading = false;
      });
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
