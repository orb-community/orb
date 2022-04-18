import { Component, Input, OnInit } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

const CONFIG = {
  SINKS: 'SINKS',
  GROUPS: 'GROUPS',
  POLICIES: 'POLICIES',
  DATASET: 'DATASET',
};

@Component({
  selector: 'ngx-dataset-from',
  templateUrl: './dataset-from.component.html',
  styleUrls: ['./dataset-from.component.scss'],
})
export class DatasetFromComponent implements OnInit {
  @Input()
  dataset: Dataset;

  @Input()
  policy: AgentPolicy;

  @Input()
  group: AgentGroup;

  isEdit: boolean;

  selectedGroup: string;

  selectedPolicy: string;

  selectedSinks: Sink[];

  availableAgentGroups: AgentGroup[];

  availableAgentPolicies: AgentPolicy[];

  availableSinks: Sink[];

  form: FormGroup;

  // #load controls
  loading = Object.entries(CONFIG)
    .reduce((acc, [value]) => {
      acc[value] = false;
      return acc;
    }, {});

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentPoliciesService: AgentPoliciesService,
    private sinksService: SinksService,
    private _formBuilder: FormBuilder,
  ) {
    this.isEdit = false;
    this.availableAgentGroups = [];
    this.availableAgentPolicies = [];
    this.availableSinks = [];
    this.selectedSinks = [];
  }

  readyForms() {
    const {
      name,
      agent_policy_id,
      agent_group_id,
      sink_ids,
    } = this?.dataset || {
      name: '',
      agent_group_id: '',
      agent_policy_id: '',
      sink_ids: [],
    } as Dataset;

    this.form = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      agent_policy_id: [agent_policy_id, [Validators.required]],
      agent_group_id: [agent_group_id, [Validators.required]],
      sink_ids: [sink_ids, [Validators.minLength(1)]],
    });
  }

  ngOnInit(): void {
    if (!!this.dataset) {
      this.isEdit = true;
    }

    if (!!this.group) {
      this.selectedGroup = this.group.id;
    }

    if (!!this.policy) {
      this.selectedPolicy = this.policy.id;
    }

    this.getDatasetAvailableConfigList();
  }

  getDatasetAvailableConfigList() {
    Promise.all([this.getAvailableAgentGroups(), this.getAvailableAgentPolicies(), this.getAvailableSinks()])
      .then(value => {

      }, reason => console.warn(`Cannot retrieve available configurations - reason: ${ JSON.parse(reason) }`))
      .catch(reason => {
        console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`);
      });
  }

  getAvailableAgentGroups() {
    return new Promise((resolve) => {
      this.loading[CONFIG.GROUPS] = true;
      const pageInfo = { ...AgentGroupsService.getDefaultPagination(), limit: 100 };
      this.agentGroupsService
        .getAgentGroups(pageInfo, false)
        .subscribe((resp: OrbPagination<AgentGroup>) => {
          this.availableAgentGroups = resp.data;
          this.loading[CONFIG.GROUPS] = false;

          resolve(this.availableAgentGroups);
        });
    });
  }

  getAvailableAgentPolicies() {
    return new Promise((resolve) => {
      this.loading[CONFIG.POLICIES] = true;
      const pageInfo = { ...AgentPoliciesService.getDefaultPagination(), limit: 100 };
      this.agentPoliciesService
        .getAgentsPolicies(pageInfo, false)
        .subscribe((resp: OrbPagination<AgentPolicy>) => {
          this.availableAgentPolicies = resp.data;
          this.loading[CONFIG.POLICIES] = false;

          resolve(this.availableAgentPolicies);
        });
    });
  }

  getAvailableSinks() {
    return new Promise((resolve) => {
      this.loading[CONFIG.SINKS] = true;
      const pageInfo = { ...SinksService.getDefaultPagination(), limit: 100 };
      this.sinksService
        .getSinks(pageInfo, false)
        .subscribe((resp: OrbPagination<Sink>) => {
          this.selectedSinks.forEach((sink) => {
            sink.name = resp.data.find(anotherSink => anotherSink.id === sink.id).name;
          });
          const sinkIDMap = this.selectedSinks.map(sink => sink.id);
          this.availableSinks = resp.data.filter(sink => !sinkIDMap.includes(sink.id));
          this.loading[CONFIG.SINKS] = false;

          resolve(this.availableSinks);
        });
    });
  }

  isLoading() {
    return Object.values<boolean>(this.loading).reduce((prev, curr) => prev && curr);
  }

  onClose() {

  }

}
