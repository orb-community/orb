import { Component, Input } from '@angular/core';
import { NbDialogRef, NbDialogService } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentMatchComponent } from '../../agents/match/agent.match.component';

@Component({
  selector: 'ngx-agent-group-details-component',
  templateUrl: './agent.group.details.component.html',
  styleUrls: ['./agent.group.details.component.scss'],
})
export class AgentGroupDetailsComponent {
  strings = STRINGS.agentGroups;

  @Input() agentGroup: AgentGroup = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentGroupDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
    private dialogService: NbDialogService,
  ) {
  }

  onMatchingAgentsModal() {
    this.dialogService.open(AgentMatchComponent, {
      context: { agentGroup: this.agentGroup },
      autoFocus: true,
      closeOnEsc: true,
    });
  }

  onOpenEdit(agentGroup: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
