import { Component, OnDestroy, OnInit } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { defer, forkJoin, Observable, of, Subscription } from 'rxjs';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { concatMap } from 'rxjs/operators';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.view.component.html',
  styleUrls: ['./agent.view.component.scss'],
})
export class AgentViewComponent implements OnInit, OnDestroy {
  strings = STRINGS.agents;

  agentStates = AgentStates;

  isLoading: boolean = true;

  agent: Agent;

  datasets: { [id: string]: Dataset };

  policies: AgentPolicy[];

  groups: AgentGroup[];

  agentID;

  agentSubscription: Subscription;

  routerSubscription: Subscription;

  constructor(
    private agentsService: AgentsService,
    private datasetService: DatasetPoliciesService,
    private groupsService: AgentGroupsService,
    private policiesService: AgentPoliciesService,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    this.agent = this.router.getCurrentNavigation()?.extras?.state?.agent as Agent || null;
    this.agentID = this.route.snapshot.paramMap.get('id');

    this.datasets = {};
    this.groups = [];
    this.policies = [];

    this.routerSubscription = this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        this.loadData();
      }
    });

  }

  ngOnInit() {
    this.isLoading = true;
    this.loadData();
  }

  loadData() {
    this.agentSubscription = this.retrieveAgentDetails()
      .subscribe({
        next: resp => {
          this.agent = resp.agent;
          this.datasets = resp?.datasets.reduce((acc, dataset) => {
            acc[dataset.id] = dataset;
            return acc;
          }, {});
          this.policies = resp?.policies;
          this.groups = resp?.groups;
        },
        complete: () => {
          this.isLoading = false;
        },
      });
  }

  retrieveAgentDetails() {
    return !!this.agentID
      && this.agentsService
        // for each AGENT
        .getAgentById(this.agentID)
        .pipe(
          // retrieve policies
          concatMap(agent => forkJoin({
              // defer execution until subscription
              // either has policies to query or not
              policies: defer(() => !!agent?.last_hb_data?.policy_state
                // fork all requests and await complete all
                && forkJoin(Object.keys(agent?.last_hb_data?.policy_state)
                  // map policy IDs to request
                  .map(policyId => this.policiesService
                    .getAgentPolicyById(policyId)))
                // or no requests at all
                || of(null)),
              // defer execution until subscription
              // and datasets for each policy too
              datasets: defer(() => !!agent?.last_hb_data?.policy_state
                // fork all requests and await complete all
                && forkJoin(Object.values(agent?.last_hb_data?.policy_state)
                  // summarize all datasets to request
                  .reduce((acc: Observable<Dataset>[], { datasets }) => {
                    return acc.concat(datasets
                      // map each datasetID to request
                      .map(dataset => this.datasetService
                        .getDatasetById(dataset)));
                  }, []) as Observable<Dataset>[])
                // or no requests at all
                || of(null)),
              groups: defer(() => !!agent?.last_hb_data?.group_state
                // fork all requests and await complete all
                && forkJoin(Object.keys(agent?.last_hb_data?.group_state)
                  // summarize all groups to request
                  .map(groupId => this.groupsService
                    .getAgentGroupById(groupId)))
                // or no requests at all
                // or no requests at all
                || of(null)),
            }),
            // emit once when all emitters emit(), completes,
            // and take(1) unsubscribes all inner observables at
            // first emission.
            (outer, inner) => ({ agent: outer, ...inner })),
        );
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
    this.routerSubscription?.unsubscribe();
  }


}
