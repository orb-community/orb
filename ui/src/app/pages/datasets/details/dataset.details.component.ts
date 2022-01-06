import { Component, Input, OnChanges, OnDestroy, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { combineLatest, forkJoin, Observable, of, Subscription, zip } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';
import { concatMapTo, filter, mergeMap, switchMap, tap } from 'rxjs/operators';

@Component({
  selector: 'ngx-dataset-details-component',
  templateUrl: './dataset.details.component.html',
  styleUrls: ['./dataset.details.component.scss'],
})
export class DatasetDetailsComponent implements OnInit {
  @Input() dataset: Dataset = {};

  fullDataset$: Observable<any>;

  agentGroup: AgentGroup;

  agentPolicy: AgentPolicy;

  sinks: Sink[];

  isLoading: boolean;

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
    this.datasetsService
      .getDatasetById(id)
      .pipe(
        tap(dataset => this.dataset = dataset),
        switchMap(dataset => forkJoin({
          agentGroup: this.agentGroupsService.getAgentGroupById(dataset.agent_group_id),
          agentPolicy: this.agentPoliciesService.getAgentPolicyById(dataset.agent_policy_id),
          sinks: this.sinksService.getSinks(null),
        })),
      )
      .subscribe(result => {
        this.agentGroup = result.agentGroup;
        this.agentPolicy = result.agentPolicy;
        this.sinks = result.sinks.data
          .filter(sink => this.dataset?.sink_ids.indexOf(sink.id) >= 0);
      });
  }

  onOpenEdit(dataset: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
