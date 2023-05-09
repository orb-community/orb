import { AfterViewInit, Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { Router } from '@angular/router';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-agent-match-component',
  templateUrl: './agent.match.component.html',
  styleUrls: ['./agent.match.component.scss'],
})

export class AgentMatchComponent implements OnInit, AfterViewInit {
  strings = STRINGS.agents;

  @Input()
  agentGroup: AgentGroup;

  @Input()
  policy!: AgentPolicy;

  agents: Agent[];

  isLoading = false;

  columnMode = ColumnMode;

  columns: TableColumn[];

  // templates
  @ViewChild('agentNameTemplateCell') agentNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

  @ViewChild('agentEspecificPolicyStateTemplateCell') agentEspecificPolicyStateTemplateRef: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
    },
    {
      id: '1',
      label: 'Tags',
      prop: 'combined_tags',
      selected: false,
    },
  ];

  constructor(
    protected dialogRef: NbDialogRef<AgentMatchComponent>,
    protected agentsService: AgentsService,
    protected router: Router,
  ) {
  }

  ngOnInit() {
    this.agents = [];
    this.updateMatchingAgents();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Agent Name',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 3,
        minWidth: 300,
        width: 450,
        maxWidth: 600,
        cellTemplate: this.agentNameTemplateCell,
      },
      {
        prop: 'combined_tags',
        name: 'Tags',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 3,
        minWidth: 300,
        width: 450,
        maxWidth: 1000,
        cellTemplate: this.agentTagsTemplateCell,
        comparator: (a, b) => Object.entries(a)
          .map(([key, value]) => `${key}:${value}`)
          .join(',')
          .localeCompare(Object.entries(b)
            .map(([key, value]) => `${key}:${value}`)
            .join(',')),
      },
      {
        prop: 'state',
        name: 'Agent Status',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 3,
        minWidth: 90,
        width: 344,
        maxWidth: 344,
        cellTemplate: this.agentStateTemplateRef,
      },
      {
        prop: 'policy_agg_info',
        name: 'Policy Status',
        resizeable: false,
        flexGrow: 3,
        canAutoResize: true,
        minWidth: 150,
        cellTemplate: this.agentEspecificPolicyStateTemplateRef,
      },
    ];
  }

  onOpenView(agent: any) {
    this.router.navigateByUrl(`pages/fleet/agents/view/${ agent.id }`);
    this.dialogRef.close();
  }

  updateMatchingAgents() {
    const { tags } = this.agentGroup;
    const tagsList = Object.keys(tags).map(key => ({ [key]: tags[key] }));
    this.agentsService.getAllAgents(tagsList).subscribe(
      resp => {
        if(!!this.policy){
          this.agents = resp.map((agent)=>{
            const {policy_state} = agent;
            const policy_agg_info = !!policy_state && policy_state[this.policy.id].state || "Not Applied";
            
            return {...agent, policy_agg_info };
          })
        } else {
          this.agents = resp;
        }
      },
    );
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
