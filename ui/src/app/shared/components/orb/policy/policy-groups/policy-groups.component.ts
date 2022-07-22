import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  OnInit,
  Output,
  SimpleChanges,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { Subscription } from 'rxjs';

@Component({
  selector: 'ngx-policy-groups',
  templateUrl: './policy-groups.component.html',
  styleUrls: ['./policy-groups.component.scss'],
})
export class PolicyGroupsComponent implements OnInit, OnChanges, OnDestroy {
  @Input() policy: AgentPolicy;

  @Input() groups: AgentGroup[];

  @Output()
  refreshPolicy: EventEmitter<string>;

  datasets: Dataset[];

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

  ngOnInit(): void {}

  ngOnChanges(changes: SimpleChanges) {}

  showAgentGroupMatches(agentGroup) {
    this.dialogService.open(AgentMatchComponent, {
      context: { agentGroup },
      autoFocus: true,
      closeOnEsc: true,
    });
  }

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`/pages/fleet/groups/edit/${agentGroup.id}`], {
      state: { agentGroup: agentGroup, edit: true },
      relativeTo: this.route,
    });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }

  unique(value, index, self) {
    return self.indexOf(value) === index;
  }
}
