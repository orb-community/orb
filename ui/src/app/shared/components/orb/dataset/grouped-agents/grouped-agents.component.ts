import { Component, Input, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { Subscription } from 'rxjs';
import { NbDialogService } from '@nebular/theme';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'ngx-grouped-agents',
  templateUrl: './grouped-agents.component.html',
  styleUrls: ['./grouped-agents.component.scss'],
})
export class GroupedAgentsComponent implements OnInit, OnChanges {
  @Input()
  agentGroup: AgentGroup;

  agents: Agent[];

  isLoading = true;

  subscription: Subscription;

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
    const { tags } = this.agentGroup;

    this.getMatchingAgents(tags);
  }

  ngOnChanges(changes: SimpleChanges) {
    const { tags } = changes.agentGroup.currentValue;

    this.getMatchingAgents(tags);
  }

  isTagsValid(tags) {
    return !!tags && Object.keys(tags).length !== 0;
  }

  getMatchingAgents(tags) {
    if (this.isTagsValid(tags)) {
      this.isLoading = true;
      this.subscription?.unsubscribe();

      this.agentsService.getAllAgents(tags)
        .subscribe(agents => {
          this.agents = agents;
          this.isLoading = false;
        });
    }
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
