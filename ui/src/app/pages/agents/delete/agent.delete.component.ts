import {Component, Input} from '@angular/core';
import {NbDialogRef} from '@nebular/theme';
import { AgentsMockService } from 'app/common/services/agents/agents.mock.service';

@Component({
  selector: 'ngx-agent-delete-component',
  templateUrl: './agent.delete.component.html',
  styleUrls: ['./agent.delete.component.scss'],
})

export class AgentDeleteComponent {
  @Input() agent = {
    name: '',
    id: '',
  };

  agentName: string = '';

  constructor(
      protected dialogRef: NbDialogRef<AgentDeleteComponent>,
      protected agentService: AgentsMockService,
  ) {
  }

  onDelete() {
    this.agentService.deleteAgent(this.agent.id);
  }

  onClose() {
    this.dialogRef.close(true);
  }

  isEnabled(): boolean {
    return this.agentName.toLowerCase() === this.agent.name.toLowerCase();
  }
}
