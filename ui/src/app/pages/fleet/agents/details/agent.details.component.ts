import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent } from 'app/common/interfaces/orb/agent.interface';

@Component({
  selector: 'ngx-agent-details-component',
  templateUrl: './agent.details.component.html',
  styleUrls: ['./agent.details.component.scss'],
})
export class AgentDetailsComponent {
  strings = STRINGS.agents;

  @Input() agent: Agent = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }


  onOpenEdit(agent: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
