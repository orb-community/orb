import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import {NbDialogService} from '@nebular/theme';

import {DropdownFilterItem} from 'app/common/interfaces/mainflux.interface';
import {ActivatedRoute, Router} from '@angular/router';
import {STRINGS} from 'assets/text/strings';
import {AgentGroupDeleteComponent} from 'app/pages/fleet/groups/delete/agent.group.delete.component';
import {ColumnMode, DatatableComponent, TableColumn} from '@swimlane/ngx-datatable';
import {AgentGroupsService} from 'app/common/services/agents/agent.groups.service';
import {NgxDatabalePageInfo, OrbPagination} from 'app/common/interfaces/orb/pagination.interface';
import {AgentGroup} from 'app/common/interfaces/orb/agent.group.interface';
import {AgentMatchComponent} from 'app/pages/fleet/agents/match/agent.match.component';
import {NotificationsService} from 'app/common/services/notifications/notifications.service';


@Component({
  selector: 'ngx-agent-group-list-component',
  templateUrl: './agent.group.list.component.html',
  styleUrls: ['./agent.group.list.component.scss'],
})
export class AgentGroupListComponent implements OnInit, AfterViewInit, AfterViewChecked {
  strings = STRINGS.agentGroups;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<AgentGroup>;

  searchPlaceholder = 'Search by name';

  // templates
  @ViewChild('agentGroupNameTemplateCell') agentGroupNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentGroupTemplateCell') agentGroupsTemplateCell: TemplateRef<any>;

  @ViewChild('agentGroupTagsTemplateCell') agentGroupTagsTemplateCell: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
      filter: (agent, name) => agent?.name.includes(name),
    },
    {
      id: '1',
      label: 'Tags',
      prop: 'tags',
      selected: false,
      filter: (agent, tag) => Object.entries(agent?.tags)
          .filter(([key, value]) => `${key}:${value}`.includes(tag.replace(' ', ''))).length > 0,
    },
    {
      id: '2',
      label: 'Description',
      prop: 'description',
      selected: false,
      filter: (agent, description) => agent?.description.includes(description),
    },
  ];

  selectedFilter = this.tableFilters[0];

  filterValue = null;

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  constructor(
      private cdr: ChangeDetectorRef,
      private dialogService: NbDialogService,
      private agentGroupsService: AgentGroupsService,
      private notificationsService: NotificationsService,
      private route: ActivatedRoute,
      private router: Router,
  ) {
    this.paginationControls = AgentGroupsService.getDefaultPagination();
  }

  ngAfterViewChecked() {
    if (this.table && this.table.recalculate && (this.tableWrapper.nativeElement.clientWidth !== this.currentComponentWidth)) {
      this.currentComponentWidth = this.tableWrapper.nativeElement.clientWidth;
      this.table.recalculate();
      this.cdr.detectChanges();
      window.dispatchEvent(new Event('resize'));
    }
  }

  ngOnInit() {
    this.getAllAgentGroups();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        canAutoResize: true,
        resizeable: false,
        flexGrow: 2,
        minWidth: 90,
        cellTemplate: this.agentGroupNameTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: false,
        minWidth: 100,
        flexGrow: 2,
      },
      {
        prop: 'matching_agents',
        name: 'Agents',
        resizeable: false,
        minWidth: 25,
        flexGrow: 1,
        comparator: (a, b) => a.total - b.total,
        cellTemplate: this.agentGroupsTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        minWidth: 300,
        flexGrow: 3,
        resizeable: false,
        cellTemplate: this.agentGroupTagsTemplateCell,
        comparator: (a, b) => Object.entries(a)
            .map(([key, value]) => `${key}:${value}`)
            .join(',')
            .localeCompare(Object.entries(b)
                .map(([key, value]) => `${key}:${value}`)
                .join(',')),
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 150,
        resizeable: false,
        sortable: false,
        flexGrow: 1,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  getAllAgentGroups(): void {
    this.agentGroupsService.clean();
    this.agentGroupsService.getAllAgentGroups().subscribe(resp => {
      this.paginationControls.data = resp.data;
      this.paginationControls.total = resp.data.length;
      this.paginationControls.offset = resp.offset / resp.limit;
      this.loading = false;
      this.cdr.markForCheck();
    });
  }

  getAgentGroups(pageInfo: NgxDatabalePageInfo = null): void {
    const finalPageInfo = {...pageInfo};
    finalPageInfo.dir = 'desc';
    finalPageInfo.order = 'name';
    finalPageInfo.limit = this.paginationControls.limit;
    finalPageInfo.offset = pageInfo?.offset * pageInfo?.limit || 0;

    this.loading = true;
    this.agentGroupsService.getAgentGroups(pageInfo).subscribe(
        (resp: OrbPagination<AgentGroup>) => {
          this.paginationControls = resp;
          this.paginationControls.offset = pageInfo?.offset || 0;
          this.paginationControls.total = resp.total;
          this.loading = false;
        },
    );
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(agentGroup: any) {
    this.router.navigate([`edit/${agentGroup.id}`], {
      state: {agentGroup: agentGroup, edit: true},
      relativeTo: this.route,
    });
  }

  onFilterSelected(filter) {
    this.searchPlaceholder = `Search by ${filter.label}`;
    this.filterValue = null;
  }

  applyFilter() {
    if (!this.paginationControls || !this.paginationControls?.data) return;

    if (!this.filterValue || this.filterValue === '') {
      this.table.rows = this.paginationControls.data;
    } else {
      this.table.rows = this.paginationControls.data.filter(sink => this.filterValue.split(/[,;]+/gm).reduce((prev, curr) => {
        return this.selectedFilter.filter(sink, curr) && prev;
      }, true));
    }
    this.paginationControls.offset = 0;
  }

  openDeleteModal(row: any) {
    const {name, id} = row;
    this.dialogService.open(AgentGroupDeleteComponent, {
      context: {name},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
        confirm => {
          if (confirm) {
            this.agentGroupsService.deleteAgentGroup(id).subscribe(() => {
              this.notificationsService.success('Agent Group successfully deleted', '');
              this.getAllAgentGroups();
            });
          }
        },
    );
  }

  onMatchingAgentsModal(row: any) {
    this.dialogService.open(AgentMatchComponent, {
      context: {agentGroup: row},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(_ => {
      this.getAgentGroups();
    });
  }
}
