import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { mergeMap } from 'rxjs/operators';
import { forkJoin, Subscription } from 'rxjs';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { NbDialogService } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'ngx-policy-groups',
  templateUrl: './policy-groups.component.html',
  styleUrls: ['./policy-groups.component.scss'],
})
export class PolicyGroupsComponent implements OnInit, OnDestroy {
  @Input() policy: AgentPolicy;

  datasets: Dataset[];

  groups: AgentGroup[];

  isLoading: boolean;

  subscription: Subscription;

  errors;

  constructor(protected datasetService: DatasetPoliciesService,
    protected groupService: AgentGroupsService,
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute) {
    this.policy = {};
    this.datasets = [];
    this.groups = [];
    this.errors = {};
  }

  ngOnInit(): void {
    this.subscription = this.retrievePolicyDatasets()
      .pipe(mergeMap(datasets => this.retrieveAgentGroups(datasets))).subscribe(resp => {
        this.datasets = resp;

        if (!this.datasets || this.datasets.length === 0) {
          this.errors['nogroup'] = 'This policy is not applied to any group.';
        }

        this.isLoading = false;
      });
  }

  retrievePolicyDatasets() {
    return this.datasetService.getAllDatasets()
      .map(resp => {
        return resp.data.filter(dataset => dataset.agent_policy_id === this.policy.id);
      });
  }

  retrieveAgentGroups(datasets: Dataset[]) {
    const groupsIds = datasets.map(dataset => dataset.agent_group_id);

    return forkJoin(groupsIds.map(id => this.groupService.getAgentGroupById(id))).map(groups => {
      this.groups = groups;
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

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`/pages/fleet/groups/edit/${ agentGroup.id }`], {
      state: { agentGroup: agentGroup, edit: true }, relativeTo: this.route,
    });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
