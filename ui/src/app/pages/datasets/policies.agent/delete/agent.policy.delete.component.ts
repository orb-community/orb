import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

@Component({
  selector: 'ngx-agent-policy-delete-component',
  templateUrl: './agent.policy.delete.component.html',
  styleUrls: ['./agent.policy.delete.component.scss'],
})

export class AgentPolicyDeleteComponent {
  @Input() name: string;

  validationInput: string = '';

  constructor(
    protected dialogRef: NbDialogRef<AgentPolicyDeleteComponent>,
  ) {
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    return this.validationInput.toLowerCase() === this.name.toLowerCase();
  }
}
