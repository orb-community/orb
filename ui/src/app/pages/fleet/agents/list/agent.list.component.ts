import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { AgentDeleteComponent } from 'app/pages/fleet/agents/delete/agent.delete.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';


@Component({
  selector: 'ngx-agent-list-component',
  templateUrl: './agent.list.component.html',
  styleUrls: ['./agent.list.component.scss'],
})
export class AgentListComponent implements OnInit, AfterViewInit, AfterViewChecked {
  strings = STRINGS.agents;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<Agent>;

  searchPlaceholder = 'Search by name';

  // templates
  @ViewChild('agentNameTemplateCell') agentNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('agentLastActivityTemplateCell') agentLastActivityTemplateCell: TemplateRef<any>;

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
      filter: (agent, tag) => Object.entries(agent?.combined_tags)
        .filter(([key, value]) => `${key}:${value}`.includes(tag.replace(' ', ''))).length > 0,
    },
    {
      id: '2',
      label: 'Status',
      prop: 'state',
      selected: false,
      filter: (agent, state) => agent?.state.includes(state),
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
    private agentService: AgentsService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.agentService.clean();
    this.paginationControls = AgentsService.getDefaultPagination();
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
    this.agentService.clean();
    this.getAllAgents();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        flexGrow: 2,
        minWidth: 90,
        cellTemplate: this.agentNameTemplateCell,
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
        prop: 'combined_tags',
        name: 'Tags',
        minWidth: 150,
        flexGrow: 4,
        cellTemplate: this.agentTagsTemplateCell,
        comparator: (a, b) => Object.entries(a)
          .map(([key, value]) => `${key}:${value}`)
          .join(',')
          .localeCompare(Object.entries(b)
            .map(([key, value]) => `${key}:${value}`)
            .join(',')),
      },
      {
        prop: 'ts_last_hb',
        name: 'Last Activity',
        minWidth: 90,
        flexGrow: 2,
        resizeable: false,
        sortable: false,
        cellTemplate: this.agentLastActivityTemplateCell,
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 130,
        resizeable: false,
        sortable: false,
        flexGrow: 2,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  getAllAgents(): void {
    this.agentService.getAllAgents().subscribe(resp => {
      this.paginationControls.data = resp.data;
      this.paginationControls.total = resp.data.length;
      this.paginationControls.offset = resp.offset / resp.limit;
      this.loading = false;
      this.cdr.markForCheck();
    });
  }

  getAgents(pageInfo: NgxDatabalePageInfo = null): void {
    const finalPageInfo = { ...pageInfo };
    finalPageInfo.dir = 'desc';
    finalPageInfo.order = 'name';
    finalPageInfo.limit = this.paginationControls.limit;
    finalPageInfo.offset = pageInfo?.offset * pageInfo?.limit || 0;

    this.loading = true;
    this.agentService.getAgents(finalPageInfo).subscribe(
      (resp: OrbPagination<Agent>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo?.offset || 0;
        this.paginationControls.total = resp.total;
        this.loading = false;
      },
    );
  }

  onOpenView(agent: any) {
    this.router.navigate([`view/${ agent.id }`], {
      relativeTo: this.route,
    });
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(agent: any) {
    this.router.navigate([`edit/${ agent.id }`], {
      state: { agent: agent, edit: true },
      relativeTo: this.route,
    });
  }

  onFilterSelected(filter) {
    this.searchPlaceholder = `Search by ${ filter.label }`;
    this.filterValue = null;
  }

  applyFilter() {
    if (!this.paginationControls || !this.paginationControls?.data) return;

    if (!this.filterValue || this.filterValue === '') {
      this.table.rows = this.paginationControls.data;
    } else {
      this.table.rows = this.paginationControls.data
        .filter(agent => this.filterValue.split(/[,;]+/gm).reduce((prev, curr) => {
        return this.selectedFilter.filter(agent, curr) && prev;
      }, true));
    }
    this.paginationControls.offset = 0;
  }

  openDeleteModal(row: any) {
    const { name, id } = row;
    this.dialogService.open(AgentDeleteComponent, {
      context: { name },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.agentService.deleteAgent(id).subscribe(() => {
            this.notificationsService.success('Agent successfully deleted', '');
            this.getAllAgents();
          });
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(AgentDetailsComponent, {
      context: { agent: row },
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

  filterByError = (agent) => !!agent && agent?.error_state && agent.error_state;

  mapRegion = (agent) => !!agent && agent?.orb_tags && !!agent.orb_tags['region'] && agent.orb_tags['region'];

  filterValid = (value) => !!value && typeof value === 'string';

  countUnique = (value, index, self) => {
    return self.indexOf(value) === index;
  }
}
