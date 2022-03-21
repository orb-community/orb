import { Component, Input, OnInit } from '@angular/core';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { NbDialogService } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'ngx-agent-groups',
  templateUrl: './agent-groups.component.html',
  styleUrls: ['./agent-groups.component.scss'],
})
export class AgentGroupsComponent implements OnInit {
  @Input() groups: AgentGroup[];

  constructor(
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute,
    ) { }

  ngOnInit(): void {
  }

  showAgentGroupDetail(agentGroup) {
    this.dialogService.open(AgentGroupDetailsComponent, {
      context: { agentGroup },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEditAgentGroup(agentGroup);
      }
    });
  }

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`../../../groups/edit/${ agentGroup.id }`], {
      state: { agentGroup: agentGroup, edit: true },
      relativeTo: this.route,
    });
  }
}
