import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

@Component({
  selector: 'ngx-policy-duplicate-component',
  templateUrl: './agent.policy.duplicate.confirmation.html',
  styleUrls: ['./agent.policy.duplicate.confirmation.scss'],
})

export class PolicyDuplicateComponent {
  @Input() policy: string
  constructor(
    protected dialogRef: NbDialogRef<PolicyDuplicateComponent>,
  ) {
  }

  onDuplicate() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

}