import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-agent-policy-details-component',
  templateUrl: './agent.policy.details.component.html',
  styleUrls: ['./agent.policy.details.component.scss'],
})
export class AgentPolicyDetailsComponent {
  strings = STRINGS.sink;

  @Input() sink: Sink = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentPolicyDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    !this.sink.tags ? this.sink.tags = {} : null;
  }

  onOpenEdit(sink: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
