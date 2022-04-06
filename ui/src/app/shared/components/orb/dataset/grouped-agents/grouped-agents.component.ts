import { Component, Input, OnInit } from '@angular/core';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { tap } from 'rxjs/operators';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { from } from 'rxjs';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { NbDialogService } from '@nebular/theme';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'ngx-grouped-agents',
  templateUrl: './grouped-agents.component.html',
  styleUrls: ['./grouped-agents.component.scss'],
})
export class GroupedAgentsComponent implements OnInit {
  @Input()
  agentGroup: AgentGroup;

  agents: Agent[];

  isLoading = true;

  errors;

  constructor(
    protected agentsService: AgentsService,
    protected dialogService: NbDialogService,
    protected router: Router,
    protected route: ActivatedRoute,
  ) {
    this.agentGroup = {};
    this.errors = {};
  }

  ngOnInit(): void {
    this.getMatchingAgents().subscribe(() => {
      this.isLoading = false;
    });
  }

  getMatchingAgents() {
    this.isLoading = true;
    const { tags } = this.agentGroup;

    return this.agentsService.getMatchingAgents(tags)
      .pipe(
        tap(agents => {
          this.agents = agents;
        }),
      );
  }

  showAgentDetails(agent) {
    this.dialogService.open(AgentDetailsComponent, {
      context: { agent },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEditAgent(agent);
      }
    });
  }

  onOpenEditAgent(agent: any) {
    this.router.navigate([`/pages/fleet/agents/edit/${ agent.id }`], {
      state: { agent: agent, edit: true },
      relativeTo: this.route,
    });
  }

}
