import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-agent-policy-details-component',
  templateUrl: './agent.policy.details.component.html',
  styleUrls: ['./agent.policy.details.component.scss'],
})
export class AgentPolicyDetailsComponent {
  @Input() agentPolicy: AgentPolicy = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentPolicyDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    !this.agentPolicy.tags ? this.agentPolicy.tags = {} : null;
  }

  onOpenEdit(agentPolicy: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
