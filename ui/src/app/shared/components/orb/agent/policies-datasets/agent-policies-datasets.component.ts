import { Component, Input, OnInit } from '@angular/core';
import { Agent, AgentPolicyState } from 'app/common/interfaces/orb/agent.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
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

  datasets: Dataset[];

  isLoading: boolean;

  errors;

  constructor(
    protected datasetService: DatasetPoliciesService,
  ) {
    this.datasets = [];
    this.errors = {};
  }

  ngOnInit(): void {
    const datasetIds = this.getDatasetIds(this?.agent?.last_hb_data?.policy_state);
    this.retrieveDatasets(datasetIds);
  }

  getDatasetIds(policyState: AgentPolicyState[]) {
    if (!policyState) {
      this.errors['nodatasets'] = 'Agent has no defined datasets.';
      return [];
    }

    const datasetIds = Object.values(policyState)
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
        this.datasets = resp;
        this.isLoading = false;
      });
  }

}
