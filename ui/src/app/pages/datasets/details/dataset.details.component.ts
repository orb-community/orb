import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { forkJoin, Subscription } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { switchMap, tap } from 'rxjs/operators';

@Component({
  selector: 'ngx-dataset-details-component',
  templateUrl: './dataset.details.component.html',
  styleUrls: ['./dataset.details.component.scss'],
})
export class DatasetDetailsComponent implements OnInit, OnDestroy {
  @Input() dataset: Dataset = {};

  agentGroup: AgentGroup;

  agentPolicy: AgentPolicy;

  sinks: Sink[];

  subscriptions: Subscription;

  constructor(
    protected dialogRef: NbDialogRef<DatasetDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
    protected agentGroupsService: AgentGroupsService,
    protected agentPoliciesService: AgentPoliciesService,
    protected datasetsService: DatasetPoliciesService,
    protected sinksService: SinksService,
  ) {

  }

  ngOnInit() {
    const { id } = this.dataset;
    this.subscriptions = this.datasetsService
      .getDatasetById(id)
      .pipe(
        tap(dataset => this.dataset = dataset),
        switchMap(dataset => forkJoin({
          agentGroup: this.agentGroupsService.getAgentGroupById(dataset.agent_group_id),
          agentPolicy: this.agentPoliciesService.getAgentPolicyById(dataset.agent_policy_id),
          sinks: this.sinksService.getAllSinks(),
        })),
      )
      .subscribe(result => {
        this.agentGroup = result.agentGroup;
        this.agentPolicy = result.agentPolicy;
        this.sinks = result.sinks.data
          .filter(sink => this.dataset?.sink_ids.indexOf(sink.id) >= 0);
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
