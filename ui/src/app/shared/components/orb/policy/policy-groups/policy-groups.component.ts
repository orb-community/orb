import {
  Component, EventEmitter, Input, OnChanges, OnDestroy, OnInit, Output,
  SimpleChanges,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import {
  AgentGroupsService,
} from 'app/common/services/agents/agent.groups.service';
import {
  DatasetPoliciesService,
} from 'app/common/services/dataset/dataset.policies.service';
import {
  AgentMatchComponent,
} from 'app/pages/fleet/agents/match/agent.match.component';
import {
  AgentGroupDetailsComponent,
} from 'app/pages/fleet/groups/details/agent.group.details.component';
import { forkJoin, Subscription } from 'rxjs';
import { mergeMap } from 'rxjs/operators';

@Component({
  selector: 'ngx-policy-groups',
  templateUrl: './policy-groups.component.html',
  styleUrls: ['./policy-groups.component.scss'],
})
export class PolicyGroupsComponent implements OnInit, OnChanges, OnDestroy {
  @Input() policy: AgentPolicy;

  @Output()
  refreshPolicy: EventEmitter<string>;

  datasets: Dataset[];

  groups: AgentGroup[];

  isLoading: boolean;

  subscription: Subscription;

  errors;

  constructor(
    protected datasetService: DatasetPoliciesService,
    protected groupService: AgentGroupsService,
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute,
  ) {
    this.refreshPolicy = new EventEmitter<string>();
    this.policy = {};
    this.datasets = [];
    this.groups = [];
    this.errors = {};
  }

  ngOnInit(): void {

  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes.policy) {
      if (this.subscription) {
        this.subscription?.unsubscribe();
      }
      this.subscription = this.retrievePolicyDatasets()
        .pipe(mergeMap(datasets => this.retrieveAgentGroups(datasets)))
        .subscribe(resp => {
          this.datasets = resp;

          if (!this.datasets || this.datasets.length === 0) {
            this.errors['nogroup'] = 'This policy is not applied to any group.';
          }

          this.isLoading = false;
        });
    }
  }

  retrievePolicyDatasets() {
    return this.datasetService.getAllDatasets()
      .map(resp => {
        return resp.data.filter(
          dataset => dataset.agent_policy_id === this.policy.id);
      });
  }

  retrieveAgentGroups(datasets: Dataset[]) {
    const groupsIds = datasets.map(dataset => dataset.agent_group_id)
      .filter(this.unique);

    if (!groupsIds || groupsIds.length === 0) {
      this.errors['nogroup'] = 'This policy is not in use by any agent group.';
    }

    return forkJoin(
      groupsIds.map(id => this.groupService.getAgentGroupById(id)))
      .map(groups => {
        this.groups = groups.filter(group => !group.error);
        this.errors.notfound = groups
          .filter(group => !!group.error)
          .map(value => `${value.id}: ${value.status} ${value.statusText}`)
          .join(',\n');
        return groups;
      });
  }

  showAgentGroupDetail(agentGroup) {
    this.dialogService.open(AgentGroupDetailsComponent, {
      context: { agentGroup }, autoFocus: true, closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEditAgentGroup(agentGroup);
      }
    });
  }

  showAgentGroupMatches(agentGroup) {
    this.dialogService.open(AgentMatchComponent, {
      context: { agentGroup },
      autoFocus: true,
      closeOnEsc: true,
    });
  }

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`/pages/fleet/groups/edit/${agentGroup.id}`], {
      state: { agentGroup: agentGroup, edit: true }, relativeTo: this.route,
    });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }

  unique(value, index, self) {
    return self.indexOf(value) === index;
  }
}
