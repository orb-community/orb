import { Component, Input } from '@angular/core';
import { NbDialogRef, NbDialogService } from '@nebular/theme';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-delete-selected-component',
  templateUrl: './delete.selected.component.html',
  styleUrls: ['./delete.selected.component.scss'],
})

export class DeleteSelectedComponent {
  strings = STRINGS.agents;
  @Input() selected: any[] = [];
  @Input() elementName: String;

  validationInput: Number;

  constructor(
    private dialogService: NbDialogService,
    protected dialogRef: NbDialogRef<DeleteSelectedComponent>,
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
  onMatchingAgentsModal(row: any) {
    this.dialogService.open(AgentMatchComponent, {
      context: { agentGroup: row },
      autoFocus: true,
      closeOnEsc: true,
    });
  }
}
