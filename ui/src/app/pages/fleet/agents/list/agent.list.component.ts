import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';
import {
  ColumnMode,
  DatatableComponent,
  TableColumn,
} from '@swimlane/ngx-datatable';

import {Agent, AgentPolicyAggStates, AgentStates} from 'app/common/interfaces/orb/agent.interface';
import {
  filterMultiSelect,
  FilterOption, filterString,
  filterTags,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { FilterService } from 'app/common/services/filter.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { AgentDeleteComponent } from 'app/pages/fleet/agents/delete/agent.delete.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { STRINGS } from 'assets/text/strings';
import { Observable } from 'rxjs';

@Component({
  selector: 'ngx-agent-list-component',
  templateUrl: './agent.list.component.html',
  styleUrls: ['./agent.list.component.scss'],
})
export class AgentListComponent implements AfterViewInit, AfterViewChecked, OnDestroy {
  strings = STRINGS.agents;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  // templates
  @ViewChild('agentNameTemplateCell') agentNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

  @ViewChild('agentPolicyStateTemplateCell') agentPolicyStateTemplateRef: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('agentLastActivityTemplateCell')
  agentLastActivityTemplateCell: TemplateRef<any>;

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;
  agents$: Observable<Agent[]>;
  filterOptions: FilterOption[];
  filters$!: Observable<FilterOption[]>;
  filteredAgents$: Observable<Agent[]>;
  private currentComponentWidth;

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private agentService: AgentsService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
    private orb: OrbService,
    private filters: FilterService,
  ) {
    this.agents$ = this.orb.getAgentListView();
    this.columns = [];

    this.filters$ = this.filters.getFilters();

    this.filterOptions = [
      {
        name: 'Name',
        prop: 'name',
        filter: filterString,
        type: FilterTypes.Input,
      },
      {
        name: 'Tags',
        prop: 'combined_tags',
        filter: filterTags,
        autoSuggestion: orb.getAgentsTags(),
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Status',
        prop: 'state',
        filter: filterMultiSelect,
        type: FilterTypes.MultiSelect,
        options: Object.values(AgentStates).map((value) => value as string),
      },
      {
        name: 'Policies',
        prop: 'policy_agg_state',
        filter: filterMultiSelect,
        type: FilterTypes.MultiSelect,
        options: Object.values(AgentPolicyAggStates).map((value) => value as string),
      },
    ];

    this.filteredAgents$ = this.filters.createFilteredList()(
      this.agents$,
      this.filters$,
      this.filterOptions,
    );
  }

  ngOnDestroy() {
    this.orb.killPolling.next();
  }

  ngAfterViewChecked() {
    if (
      this.table &&
      this.table.recalculate &&
      this.tableWrapper.nativeElement.clientWidth !== this.currentComponentWidth
    ) {
      this.currentComponentWidth = this.tableWrapper.nativeElement.clientWidth;
      this.table.recalculate();
      this.cdr.detectChanges();
      window.dispatchEvent(new Event('resize'));
    }
  }

  ngAfterViewInit() {
    this.orb.refreshNow();
    this.columns = [
      {
        prop: 'name',
        flexGrow: 4,
        canAutoResize: true,
        minWidth: 150,
        name: 'Name',
        cellTemplate: this.agentNameTemplateCell,
      },
      {
        prop: 'state',
        flexGrow: 2,
        canAutoResize: true,
        name: 'Status',
        cellTemplate: this.agentStateTemplateRef,
      },
      {
        prop: 'policy_agg_info',
        flexGrow: 2,
        canAutoResize: true,
        minWidth: 100,
        name: 'Policies',
        cellTemplate: this.agentPolicyStateTemplateRef,
      },
      {
        prop: 'combined_tags',
        flexGrow: 7,
        canAutoResize: true,
        name: 'Tags',
        cellTemplate: this.agentTagsTemplateCell,
        comparator: (a, b) =>
          Object.entries(a)
            .map(([key, value]) => `${key}:${value}`)
            .join(',')
            .localeCompare(
              Object.entries(b)
                .map(([key, value]) => `${key}:${value}`)
                .join(','),
            ),
      },
      {
        prop: 'ts_last_hb',
        flexGrow: 2,
        minWidth: 150,
        canAutoResize: true,
        name: 'Last Activity',
        sortable: false,
        cellTemplate: this.agentLastActivityTemplateCell,
      },
      {
        name: '',
        prop: 'actions',
        flexGrow: 4,
        minWidth: 200,
        canAutoResize: true,
        sortable: false,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
  }

  onOpenMetrics(agent: any) {
    this.router.navigate([`metrics/${agent.id}`], {
      relativeTo: this.route,
    });
  }

  onOpenView(agent: any) {
    this.router.navigate([`view/${agent.id}`], {
      relativeTo: this.route,
    });
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(agent: any) {
    this.router.navigate([`edit/${agent.id}`], {
      state: { agent: agent, edit: true },
      relativeTo: this.route,
    });
  }

  openDeleteModal(row: any) {
    const { name, id } = row;
    this.dialogService
      .open(AgentDeleteComponent, {
        context: { name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.agentService.deleteAgent(id).subscribe(() => {
            this.notificationsService.success('Agent successfully deleted', '');
            this.orb.refreshNow();
          });
        }
      });
  }

  openDetailsModal(row: any) {
    this.dialogService
      .open(AgentDetailsComponent, {
        context: { agent: row },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEdit(row);
        }
      });
  }

  filterByError = (agent) => !!agent && agent?.error_state && agent.error_state;

  mapRegion = (agent) =>
    !!agent &&
    agent?.orb_tags &&
    !!agent.orb_tags['region'] &&
    agent.orb_tags['region']

  filterValid = (value) => !!value && typeof value === 'string';

  countUnique = (value, index, self) => self.indexOf(value) === index;
}
