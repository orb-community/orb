import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { AgentsService } from 'app/common/services/agents/agents.service';

@Component({
  selector: 'ngx-agent-group-delete-component',
  templateUrl: './agentGroup.delete.component.html',
  styleUrls: ['./agentGroup.delete.component.scss'],
})

export class AgentGroupDeleteComponent {
  @Input() agentGroup = {
    name: '',
    id: '',
  };

  agentName: string = '';

  constructor(
    protected dialogRef: NbDialogRef<AgentGroupDeleteComponent>,
    protected agentService: AgentsService,
  ) {
  }

  onDelete() {
    this.agentService.deleteAgentGroup(this.agentGroup.id);
  }

  onClose() {
    this.dialogRef.close(true);
  }

  isEnabled(): boolean {
    return this.agentName.toLowerCase() === this.agentGroup.name.toLowerCase();
  }
}
