import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';

@Component({
  selector: 'ngx-agent-details-component',
  templateUrl: './agent.details.component.html',
  styleUrls: ['./agent.details.component.scss'],
})
export class AgentDetailsComponent {
  strings = STRINGS.agents;

  @Input() agentGroup: AgentGroup = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }


  onOpenEdit(row: any) {
    this.router.navigate(['../agents/edit'], {
      relativeTo: this.route,
      queryParams: {id: row.id},
      state: {agentGroup: row},
    });
  }

  onClose() {
    this.dialogRef.close();
  }
}
