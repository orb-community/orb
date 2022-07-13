import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { forkJoin } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { ActivatedRoute, Router } from '@angular/router';
import { DatasetFromComponent } from 'app/pages/datasets/dataset-from/dataset-from.component';
import { NbDialogService } from '@nebular/theme';
import {
  AgentPolicyState,
  AgentPolicyStates,
} from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-agent-policies-datasets',
  templateUrl: './agent-policies-datasets.component.html',
  styleUrls: ['./agent-policies-datasets.component.scss'],
})
export class AgentPoliciesDatasetsComponent implements OnInit {
  @Input() agent: Agent;

  @Output()
  refreshAgent: EventEmitter<string>;

  policyStates = AgentPolicyStates;

  datasets: { [id: string]: Dataset };

  policies: AgentPolicyState[];

  isLoading: boolean;

  errors;

  constructor(
    private datasetService: DatasetPoliciesService,
    private router: Router,
    private route: ActivatedRoute,
    private dialogService: NbDialogService,
  ) {
    this.refreshAgent = new EventEmitter<string>();
    this.datasets = {};
    this.errors = {};
  }

  ngOnInit(): void {
    this.policies = this.getPoliciesStates(
      this?.agent?.last_hb_data?.policy_state,
    );
    const datasetIds = this.getDatasetIds(this.policies);
    this.retrieveDatasets(datasetIds);
  }

  getPoliciesStates(policyStates: { [id: string]: AgentPolicyState }) {
    if (!policyStates || policyStates === {}) {
      this.errors['nodatasets'] = 'Agent has no defined policies.';
      return [];
    }
    return Object.entries(policyStates).map(([id, policy]) => {
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
      .map((state) => state?.datasets)
      .reduce((acc, curr) => curr.concat(acc), []);

    return datasetIds;
  }

  retrieveDatasets(datasetIds: string[]) {
    if (!datasetIds) {
      return;
    }
    this.isLoading = true;
    forkJoin(
      datasetIds.map((id) => this.datasetService.getDatasetById(id)),
    ).subscribe((resp) => {
      this.datasets = resp.reduce(
        (acc: { [id: string]: Dataset }, curr: Dataset) => {
          acc[curr.id] = curr;
          return acc;
        },
        {},
      );
      this.isLoading = false;
    });
  }

  onOpenViewPolicy(policy: any) {
    this.router.navigate([`/pages/datasets/policies/view/${policy.id}`], {
      state: { policy: policy },
      relativeTo: this.route,
    });
  }

  onOpenViewDataset(dataset: any) {
    this.dialogService
      .open(DatasetFromComponent, {
        autoFocus: true,
        closeOnEsc: true,
        context: {
          dataset,
        },
        hasScroll: false,
        hasBackdrop: false,
      })
      .onClose.subscribe((resp) => {
        if (resp === 'changed' || 'deleted') {
          this.refreshAgent.emit('refresh-from-dataset');
        }
      });
  }
}
