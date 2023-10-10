import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-agent-reset-component',
  templateUrl: './agent.reset.component.html',
  styleUrls: ['./agent.reset.component.scss'],
})

export class AgentResetComponent {
  strings = STRINGS.agents;
  @Input() selected: any[] = [];

  validationInput: Number;

  constructor(
    protected dialogRef: NbDialogRef<AgentResetComponent>,
  ) {
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    return this.validationInput === this.selected.length;
  }
}
