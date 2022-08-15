import {
  Component,
  Input,
  OnChanges,
  OnInit,
  SimpleChanges,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';

@Component({
  selector: 'ngx-policy-groups',
  templateUrl: './policy-groups.component.html',
  styleUrls: ['./policy-groups.component.scss'],
})
export class PolicyGroupsComponent implements OnInit, OnChanges {
  @Input() groups: AgentGroup[];

  isLoading: boolean;

  errors;

  constructor(
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute,
  ) {
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

  unique(value, index, self) {
    return self.indexOf(value) === index;
  }
}
