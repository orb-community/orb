import { AfterViewInit, Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';

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
        flexGrow: 1,
        minWidth: 90,
      },
      {
        prop: 'combined_tags',
        name: 'Tags',
        resizeable: false,
        minWidth: 100,
        flexGrow: 2,
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
        minWidth: 90,
        flexGrow: 1,
        cellTemplate: this.agentStateTemplateRef,
      },
    ];
  }

  updateMatchingAgents() {
    const { tags } = this.agentGroup;
    const tagsList = Object.keys(tags).map(key => ({ [key]: tags[key] }));
    this.agentsService.getMatchingAgents(tagsList).subscribe(
      resp => {
        this.agents = resp;
      },
    );
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
