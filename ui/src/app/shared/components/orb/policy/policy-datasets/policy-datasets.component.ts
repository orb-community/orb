import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Subscription } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';

@Component({
  selector: 'ngx-policy-datasets',
  templateUrl: './policy-datasets.component.html',
  styleUrls: ['./policy-datasets.component.scss'],
})
export class PolicyDatasetsComponent implements OnInit, OnDestroy {
  @Input()
  policy: AgentPolicy;

  datasets: Dataset[];

  isLoading: boolean;

  subscription: Subscription;

  errors;

  constructor(protected datasetService: DatasetPoliciesService) {
    this.policy = {};
    this.datasets = [];
    this.errors = {};
  }

  ngOnInit(): void {
    this.subscription = this.retrievePolicyDatasets()
      .subscribe(resp => {
        this.isLoading = false;
      });
  }

  retrievePolicyDatasets() {
    return this.datasetService.getAllDatasets()
      .map(resp => {
        this.datasets = resp.data.filter(dataset => dataset.agent_policy_id === this.policy.id);
        return this.datasets;
      });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
