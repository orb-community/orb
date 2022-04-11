import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-agent-group-delete-component',
  templateUrl: './agent.group.delete.component.html',
  styleUrls: ['./agent.group.delete.component.scss'],
})

export class AgentGroupDeleteComponent {
  strings = STRINGS.agentGroups;
  @Input() name: string;

  validationInput: string = '';

  constructor(
    protected dialogRef: NbDialogRef<AgentGroupDeleteComponent>,
  ) {
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    return this.validationInput === this.name;
  }
}
