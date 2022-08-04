import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { forkJoin, Subscription } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { ActivatedRoute, Router } from '@angular/router';
import { DatasetFromComponent } from 'app/pages/datasets/dataset-from/dataset-from.component';
import { NbDialogService } from '@nebular/theme';
import {
  AgentPolicyState,
  AgentPolicyStates,
} from 'app/common/interfaces/orb/agent.policy.interface';
import { OrbService } from 'app/common/services/orb.service';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-agent-policies-datasets',
  templateUrl: './agent-policies-datasets.component.html',
  styleUrls: ['./agent-policies-datasets.component.scss'],
})
export class AgentPoliciesDatasetsComponent implements OnInit {
  @Input() agent: Agent;

  @Input()  
  datasets: { [id: string]: Dataset };
  
  @Output()
  refreshAgent: EventEmitter<string>;

  policyStates = AgentPolicyStates;

  policies: AgentPolicyState[];

  errors;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private dialogService: NbDialogService,
  ) {
    this.refreshAgent = new EventEmitter<string>();
    this.datasets = {};
    this.policies = [];
    this.errors = {};
  }

  ngOnInit(): void {
    const policiesStates = this.agent?.last_hb_data?.policy_state;
    if (!policiesStates || policiesStates === []) {
      this.errors['nodatasets'] = 'Agent has no defined datasets.';
    } else {
      this.policies = this.getPoliciesStates(policiesStates);
    }
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
        hasBackdrop: true,
        closeOnBackdropClick: true,
      })
      .onClose.subscribe((resp) => {
        if (resp === 'changed' || 'deleted') {
          this.refreshAgent.emit('refresh-from-dataset');
        }
      });
  }
}
