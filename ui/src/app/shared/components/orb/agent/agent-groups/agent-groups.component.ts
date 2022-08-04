import { Component, Input, OnInit } from '@angular/core';
import {
  AgentGroup,
  AgentGroupState,
} from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { NbDialogService } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent } from 'app/common/interfaces/orb/agent.interface';

@Component({
  selector: 'ngx-agent-groups',
  templateUrl: './agent-groups.component.html',
  styleUrls: ['./agent-groups.component.scss'],
})
export class AgentGroupsComponent implements OnInit {
  @Input() agent: Agent;

  @Input()
  groups: AgentGroup[];

  errors;

  constructor(
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute,
  ) {
    this.groups = [];
    this.errors = {};
  }

  ngOnInit(): void {
    if (!this.groups || this.groups.length === 0) {
      this.errors['nogroup'] = 'This agent does not belong to any group.';
      return;
    }
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
