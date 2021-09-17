import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-agent-policy-delete-component',
  templateUrl: './agent.policy.delete.component.html',
  styleUrls: ['./agent.policy.delete.component.scss'],
})

export class AgentPolicyDeleteComponent {
  @Input() sink;

  strings = STRINGS.sink;

  userInput: string = '';

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
    return this.userInput === this.sink.name;
  }
}
