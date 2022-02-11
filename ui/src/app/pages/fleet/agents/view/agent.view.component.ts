import { Component, OnDestroy } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService, AvailableOS } from 'app/common/services/agents/agents.service';
import { forkJoin, Observable, Subscription } from 'rxjs';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { concatMap, take, tap } from 'rxjs/operators';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.view.component.html',
  styleUrls: ['./agent.view.component.scss'],
})
export class AgentViewComponent implements OnDestroy {
  strings = STRINGS.agents;

  agentStates = AgentStates;

  isLoading: boolean = true;

  agent: Agent;

  groups: AgentGroup[];

  datasets: Dataset[];

  policies: AgentPolicy[];

  agentID;

  command2copy: string;

  copyCommandIcon: string;

  availableOS = [AvailableOS.DOCKER];

  selectedOS = AvailableOS.DOCKER;

  command2show: string;

  hideCommand: boolean;

  subscription: Subscription;

  constructor(
    private agentsService: AgentsService,
    private agentGroupsService: AgentGroupsService,
    private policiesService: AgentPoliciesService,
    private datasetService: DatasetPoliciesService,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    this.agent = this.router.getCurrentNavigation().extras.state?.agent as Agent || null;
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.command2copy = '';
    this.command2show = '';
    this.copyCommandIcon = 'clipboard-outline';
    this.hideCommand = this.agent?.state !== this.agentStates.new;
    this.subscription = this.loadData()
      .subscribe(resp => {
        this.agent = resp.agent;
        this.datasets = resp.datasets;
        this.policies = resp.policies;
        this.makeCommand2Copy();
        this.isLoading = false;
      });
  }

  loadData() {
    return !!this.agentID
      && this.agentsService
        // for each AGENT
        .getAgentById(this.agentID)
        .pipe(
          // retrieve policies
          concatMap(agent => forkJoin({
              policies: forkJoin(Object.keys(agent?.last_hb_data?.policy_state)
                .map(policyId => this.policiesService.getAgentPolicyById(policyId).pipe(take(1)))).pipe(take(1)),
              // and datasets for each policy too
              datasets: forkJoin(Object.values(agent?.last_hb_data?.policy_state)
                .reduce((acc: Observable<Dataset>[], { datasets }) => {
                  return acc.concat(datasets.map(dataset => this.datasetService.getDatasetById(dataset).pipe(take(1))));
                }, []) as Observable<Dataset>[]).pipe(take(1)),
            }),
            // emit once
            (outer, inner) => ({ agent: outer, ...inner })),
        );
  }

  toggleIcon(target) {
    if (target === 'command') {
      this.copyCommandIcon = 'checkmark-outline';
    }
  }

  isToday() {
    const today = new Date(Date.now());
    const date = new Date(this?.agent?.ts_last_hb);

    return today.getDay() === date.getDay()
      && today.getMonth() === date.getMonth()
      && today.getFullYear() === date.getFullYear();

  }

  makeCommand2Copy() {
    // TODO: future - store this elsewhere
    if (this.selectedOS === AvailableOS.DOCKER) {
      this.command2copy = `docker run -d --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY="PASTE_AGENT_KEY" \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=mock \\
ns1labs/orb-agent:develop`;

      this.command2show = `docker run -d --net=host \n
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \n
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \n
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \n
-e ORB_CLOUD_MQTT_KEY=<mark>{{ AGENT KEY }}</mark> \n
-e PKTVISOR_PCAP_IFACE_DEFAULT=<mark>mock</mark> \n

ns1labs/orb-agent:develop`;
    }
  }

  toggleProvisioningCommand() {
    this.hideCommand = !this.hideCommand;
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
