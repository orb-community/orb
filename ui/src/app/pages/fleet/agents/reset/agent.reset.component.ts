import { Component, Input, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-agent-reset-component',
  templateUrl: './agent.reset.component.html',
  styleUrls: ['./agent.reset.component.scss'],
})

export class AgentResetComponent implements OnInit {
  strings = STRINGS.agents;
  @Input() selected: any[] = [];
  @Input() agent: Agent;

  validationInput: any;


  constructor(
    protected dialogRef: NbDialogRef<AgentResetComponent>,
  ) {
  }

  ngOnInit(): void {
    if (this.agent) {
      this.selected = [this.agent];
    }
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    if (this.agent) {
      return true;
    } else {
      return this.validationInput === this.selected.length;
    }
  }
}
