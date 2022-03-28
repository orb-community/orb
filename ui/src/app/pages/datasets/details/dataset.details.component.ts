import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { forkJoin, Subscription } from 'rxjs';

@Component({
  selector: 'ngx-dataset-details-component',
  templateUrl: './dataset.details.component.html',
  styleUrls: ['./dataset.details.component.scss'],
})
export class DatasetDetailsComponent implements OnInit, OnDestroy {
  @Input() dataset: Dataset = {};

  agentGroup: AgentGroup;

  agentPolicy: AgentPolicy;

  subscriptions: Subscription;

  errors: { [propName: string]: string };

  isLoading: boolean;

  constructor(
    protected dialogRef: NbDialogRef<DatasetDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
    protected agentGroupsService: AgentGroupsService,
    protected agentPoliciesService: AgentPoliciesService,
  ) {
    this.errors = {};
  }

  ngOnInit() {
    const { agent_group_id, agent_policy_id } = this.dataset;

    this.isLoading = true;
    this.subscriptions = forkJoin({
      agentGroup: this.agentGroupsService.getAgentGroupById(agent_group_id),
      agentPolicy: this.agentPoliciesService.getAgentPolicyById(agent_policy_id),
    })
      .subscribe(result => {

        // agent group
        if (!!result.agentGroup['error']) {
          this.errors['agentGroup'] = 'Failed to fetch Agent Group';
        } else {
          this.agentGroup = result.agentGroup;
        }

        // agent policy
        if (!!result.agentPolicy['error']) {
          this.errors['agentPolicy'] = 'Failed to fetch Agent Policy';
        } else {
          this.agentPolicy = result.agentPolicy;
        }

        this.isLoading = false;
      }, err => {
        this.errors['error'] = err.message;
      });
  }

  ngOnDestroy() {
    this.subscriptions?.unsubscribe();
  }

  onOpenEdit(dataset: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
