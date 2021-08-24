import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { AgentsService } from 'app/common/services/agents/agents.service';

@Component({
  selector: 'ngx-agent-delete-component',
  templateUrl: './agent.delete.component.html',
  styleUrls: ['./agent.delete.component.scss'],
})

export class AgentDeleteComponent {
  @Input() agentGroup = {
    name: '',
    id: '',
  };

  agentName: string = '';

  constructor(
    protected dialogRef: NbDialogRef<AgentDeleteComponent>,
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
