import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';

@Component({
  selector: 'ngx-agent-group-details-component',
  templateUrl: './agent.group.details.component.html',
  styleUrls: ['./agent.group.details.component.scss'],
})
export class AgentGroupDetailsComponent {
  strings = STRINGS.agents;

  @Input() agentGroup: AgentGroup = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentGroupDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }


  onOpenEdit(agentGroup: any) {
    this.dialogRef.close();
    this.router.navigate([`pages/agent-group/edit/${agentGroup.id}`], {
      state: {agentGroup: agentGroup, edit: true},
    });
  }

  onClose() {
    this.dialogRef.close();
  }
}
