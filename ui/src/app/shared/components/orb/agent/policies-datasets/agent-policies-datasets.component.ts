import { Component, Input, OnInit } from '@angular/core';
import { Agent, AgentPolicyState } from 'app/common/interfaces/orb/agent.interface';
import { forkJoin } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';

@Component({
  selector: 'ngx-agent-policies-datasets',
  templateUrl: './agent-policies-datasets.component.html',
  styleUrls: ['./agent-policies-datasets.component.scss'],
})
export class AgentPoliciesDatasetsComponent implements OnInit {
  @Input() agent: Agent;

  datasets: { [id: string]: Dataset };

  policies: AgentPolicyState[];

  isLoading: boolean;

  errors;

  constructor(
    protected datasetService: DatasetPoliciesService,
  ) {
    this.datasets = {};
    this.errors = {};
  }

  ngOnInit(): void {
    this.policies = this.getPoliciesStates(this?.agent?.last_hb_data?.policy_state);
    const datasetIds = this.getDatasetIds(this.policies);
    this.retrieveDatasets(datasetIds);
  }

  getPoliciesStates(policyStates: { [id: string]: AgentPolicyState }) {
    return Object.entries(policyStates)
      .map(([id, policy]) => {
        policy.id = id;
        return policy;
      });
  }

  getDatasetIds(policiesStates: AgentPolicyState[]) {
    if (!policiesStates || policiesStates === []) {
      this.errors['nodatasets'] = 'Agent has no defined datasets.';
      return [];
    }

    const datasetIds = policiesStates
      .map(state => state?.datasets)
      .reduce((acc, curr) => curr.concat(acc), []);

    return datasetIds;
  }

  retrieveDatasets(datasetIds: string[]) {
    if (!datasetIds) {
      return;
    }
    this.isLoading = true;
    forkJoin(datasetIds.map(id => this.datasetService.getDatasetById(id)))
      .subscribe(resp => {
        this.datasets = resp.reduce((acc: { [id: string]: Dataset }, curr: Dataset) => {
          acc[curr.id] = curr;
          return acc;
        }, {});
        this.isLoading = false;
      });
  }

}
