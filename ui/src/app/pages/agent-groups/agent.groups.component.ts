import { AfterViewInit, ChangeDetectorRef, Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { AgentGroupDeleteComponent } from 'app/pages/agent-groups/delete/agent.group.delete.component';
import { AgentGroupDetailsComponent } from 'app/pages/agent-groups/details/agent.group.details.component';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Debounce } from 'app/shared/decorators/utils';


@Component({
  selector: 'ngx-agent-groups-component',
  templateUrl: './agent.groups.component.html',
  styleUrls: ['./agent.groups.component.scss'],
})
export class AgentGroupsComponent implements OnInit, AfterViewInit {
  strings = STRINGS.agents;

  columnMode = ColumnMode;
  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<AgentGroup>;

  searchPlaceholder = 'Search by name';
  filterSelectedIndex = '0';

  // templates

  @ViewChild('agentGroupTemplateCell') agentGroupsTemplateCell: TemplateRef<any>;
  @ViewChild('agentGroupTagsTemplateCell') agentGroupTagsTemplateCell: TemplateRef<any>;
  @ViewChild('addAgentGroupTemplateRef') addAgentGroupTemplateRef: TemplateRef<any>;
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
    private agentGroupsService: AgentGroupsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.agentGroupsService.clean();
    this.paginationControls = AgentGroupsService.getDefaultPagination();
  }

  ngOnInit() {
    this.agentGroupsService.clean();
    this.getAgentGroups();
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
        prop: 'description',
        name: 'Description',
        resizeable: false,
        minWidth: 100,
        flexGrow: 2,
      },
      {
        prop: 'matching_agents.total',
        name: 'Agents',
        resizeable: false,
        minWidth: 100,
        flexGrow: 1,
        cellTemplate: this.agentGroupsTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        minWidth: 90,
        flexGrow: 3,
        cellTemplate: this.agentGroupTagsTemplateCell,
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
  getAgentGroups(pageInfo: NgxDatabalePageInfo = null): void {
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
    this.agentGroupsService.getAgentGroups(pageInfo, isFilter).subscribe(
      (resp: OrbPagination<AgentGroup>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo.offset;
        this.loading = false;
      },
    );
  }

  onOpenAdd() {
    this.router.navigate(['../agent-group/add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(row: any) {
    this.router.navigate(['../agent-group/edit'], {
      relativeTo: this.route,
      queryParams: {id: row.id},
      state: {agentGroup: row, edit: true},
    });
  }

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${this.tableFilters[selectedIndex].label}`;
  }

  openDeleteModal(row: any) {
    const {name} = row;
    this.dialogService.open(AgentGroupDeleteComponent, {
      context: {name},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.agentGroupsService.deleteAgentGroup(row.id).subscribe(() => this.getAgentGroups());
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(AgentGroupDetailsComponent, {
      context: {agentGroup: row},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(() => this.getAgentGroups());
  }

  searchAgentByName(input) {
    this.getAgentGroups({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }

  filterByActive = (agent) => agent.status === 'active';
}
