import { AfterViewInit, Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { Router } from '@angular/router';

@Component({
  selector: 'ngx-agent-match-component',
  templateUrl: './agent.match.component.html',
  styleUrls: ['./agent.match.component.scss'],
})

export class AgentMatchComponent implements OnInit, AfterViewInit {
  strings = STRINGS.agents;

  @Input()
  agentGroup: AgentGroup;

  agents: Agent[];

  isLoading = false;

  columnMode = ColumnMode;

  columns: TableColumn[];

  // templates
  @ViewChild('agentNameTemplateCell') agentNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

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
        flexGrow: 2,
        minWidth: 90,
        width: 120,
        maxWidth: 200,
        cellTemplate: this.agentNameTemplateCell,
      },
      {
        prop: 'combined_tags',
        name: 'Tags',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 6,
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
        name: 'Status',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 2,
        minWidth: 90,
        width: 90,
        maxWidth: 90,
        cellTemplate: this.agentStateTemplateRef,
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
        this.agents = resp;
      },
    );
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
