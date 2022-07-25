import { Component, Input, OnInit } from '@angular/core';
import {
  AgentGroup,
  AgentGroupState,
} from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { NbDialogService } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { forkJoin } from 'rxjs';

@Component({
  selector: 'ngx-agent-groups',
  templateUrl: './agent-groups.component.html',
  styleUrls: ['./agent-groups.component.scss'],
})
export class AgentGroupsComponent implements OnInit {
  @Input() agent: Agent;

  groups: AgentGroup[];

  isLoading: boolean;

  errors;

  constructor(
    protected groupsService: AgentGroupsService,
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute,
  ) {
    this.groups = [];
    this.errors = {};
  }

  ngOnInit(): void {
    this.retrieveGroups(this.agent?.last_hb_data?.group_state);
  }

  retrieveGroups(groupState: AgentGroupState) {
    if (!groupState || groupState === {}) {
      this.errors['nogroup'] = 'This agent does not belong to any group.';
      return;
    }

    this.isLoading = true;

    const groupIds = Object.keys(groupState);

    forkJoin(
      groupIds.map((id) => this.groupsService.getAgentGroupById(id)),
    ).subscribe((resp) => {
      this.groups = resp.filter((group) => !group.error);
      this.errors.notfound = resp
        .filter((group) => !!group.error)
        .map((value) => `${value.id}: ${value.status} ${value.statusText}`)
        .join(',\n');
      this.isLoading = false;
    });
  }

  showAgentGroupDetail(agentGroup) {
    this.dialogService
      .open(AgentGroupDetailsComponent, {
        context: { agentGroup },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEditAgentGroup(agentGroup);
        }
      });
  }

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`/pages/fleet/groups/edit/${agentGroup.id}`], {
      state: { agentGroup: agentGroup, edit: true },
      relativeTo: this.route,
    });
  }
}
