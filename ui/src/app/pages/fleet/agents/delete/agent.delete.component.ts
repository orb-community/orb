import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-agent-delete-component',
  templateUrl: './agent.delete.component.html',
  styleUrls: ['./agent.delete.component.scss'],
})

export class AgentDeleteComponent {
  strings = STRINGS.agents;
  @Input() name: string;

  validationInput: string = '';

  constructor(
    protected dialogRef: NbDialogRef<AgentDeleteComponent>,
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
