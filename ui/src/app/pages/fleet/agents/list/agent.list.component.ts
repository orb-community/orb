import { AfterViewInit, ChangeDetectorRef, Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination';
import { Debounce } from 'app/shared/decorators/utils';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { AgentDeleteComponent } from 'app/pages/fleet/agents/delete/agent.delete.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';


@Component({
  selector: 'ngx-agent-list-component',
  templateUrl: './agent.list.component.html',
  styleUrls: ['./agent.list.component.scss'],
})
export class AgentListComponent implements OnInit, AfterViewInit {
  strings = STRINGS.agents;

  columnMode = ColumnMode;
  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<Agent>;

  searchPlaceholder = 'Search by name';
  filterSelectedIndex = '0';

  // templates
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;
  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;
  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

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
      prop: 'tags',
      selected: false,
    },
  ];

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private agentService: AgentsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.agentService.clean();
    this.paginationControls = AgentsService.getDefaultPagination();
  }

  ngOnInit() {
    this.agentService.clean();
    this.getAgents();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        flexGrow: 1,
        minWidth: 90,
      },
      {
        prop: 'agent_tags',
        name: 'Agent Tags',
        resizeable: false,
        minWidth: 90,
        flexGrow: 2,
        cellTemplate: this.agentTagsTemplateCell,
      },
      {
        prop: 'state',
        name: 'Status',
        resizeable: false,
        minWidth: 100,
        flexGrow: 1,
        cellTemplate: this.agentStateTemplateRef,
      },
      {
        prop: 'orb_tags',
        name: 'Orb Tags',
        minWidth: 90,
        flexGrow: 3,
        cellTemplate: this.agentTagsTemplateCell,
      },
      {
        prop: 'ts_lst_hb',
        name: 'Last Activity',
        minWidth: 90,
        flexGrow: 3,
        resizeable: false,
        sortable: false,
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 130,
        resizeable: false,
        sortable: false,
        flexGrow: 1,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  @Debounce(400)
  getAgents(pageInfo: NgxDatabalePageInfo = null): void {
    const isFilter = pageInfo === null;
    if (isFilter) {
      pageInfo = {
        offset: this.paginationControls.offset,
        limit: this.paginationControls.limit,
      };
      if (this.paginationControls.name?.length > 0) pageInfo.name = this.paginationControls.name;
      if (this.paginationControls.tags?.length > 0) pageInfo.tags = this.paginationControls.tags;
    }

    this.loading = true;
    this.agentService.getAgents(pageInfo, isFilter).subscribe(
      (resp: OrbPagination<Agent>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo.offset;
        this.loading = false;
      },
    );
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(agent: any) {
    this.router.navigate([`edit/${agent.id}`], {
      state: {agent: agent, edit: true},
      relativeTo: this.route,
    });
  }

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${this.tableFilters[selectedIndex].label}`;
  }

  openDeleteModal(row: any) {
    const {name, id} = row;
    this.dialogService.open(AgentDeleteComponent, {
      context: {name},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.agentService.deleteAgent(id).subscribe(() => this.getAgents());
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(AgentDetailsComponent, {
      context: {agent: row},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEdit(row);
      } else {
        this.getAgents();
      }
    });
  }

  searchAgentByName(input) {
    this.getAgents({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }

  filterByActive = (agent) => agent.state === 'active';
}
