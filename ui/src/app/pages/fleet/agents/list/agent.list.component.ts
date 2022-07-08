import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
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

import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import {
  FilterOption,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { FilterService } from 'app/common/services/filter.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { AgentDeleteComponent } from 'app/pages/fleet/agents/delete/agent.delete.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { STRINGS } from 'assets/text/strings';
import { combineLatest, Observable } from 'rxjs';
import { map, startWith } from 'rxjs/operators';

@Component({
  selector: 'ngx-agent-list-component',
  templateUrl: './agent.list.component.html',
  styleUrls: ['./agent.list.component.scss'],
})
export class AgentListComponent implements AfterViewInit, AfterViewChecked {
  strings = STRINGS.agents;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  // templates
  @ViewChild('agentNameTemplateCell') agentNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

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

    this.filters$ = this.filters.getFilters().pipe(startWith([]));

    this.filteredAgents$ = combineLatest([this.agents$, this.filters$]).pipe(
      map(([agents, _filters]) => {
        let filtered = agents;
        _filters.forEach((_filter) => {
          filtered = filtered.filter((value) => {
            const paramValue = _filter.param;
            const result = _filter.filter(value, paramValue);
            return result;
          });
        });

        return filtered;
      }),
    );

    this.filterOptions = [
      {
        name: 'Name',
        prop: 'name',
        filter: (agent: Agent, name: string) => {
          return agent.name?.includes(name);
        },
        type: FilterTypes.Input,
      },
      {
        name: 'Agent Tags',
        prop: 'agent_tags',
        filter: (agent: Agent, tag: string) => {
          const values = Object.entries(agent.agent_tags)
            .map((entry) => `${entry[0]}: ${entry[1]}`)
            .reduce((acc, val) => acc.concat(val), []);
          return values.reduce((acc, val) => {
            acc |= val.includes(tag.trim());
            return acc;
          }, false);
        },
        autoSuggestion: orb.getAgentsTags(),
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Orb Tags',
        prop: 'orb_tags',
        filter: (agent: Agent, tag: string) => {
          const values = Object.entries(agent.orb_tags)
            .map((entry) => `${entry[0]}: ${entry[1]}`)
            .reduce((acc, val) => acc.concat(val), []);
          return values.reduce((acc, val) => {
            acc |= val.includes(tag.trim());
            return acc;
          }, false);
        },
        autoSuggestion: orb.getAgentsTags(),
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Status',
        prop: 'state',
        filter: (agent: Agent, states: string[]) => {
          return states.reduce((prev, cur) => {
            return agent.state === cur || prev;
          }, false);
        },
        type: FilterTypes.MultiSelect,
        options: Object.values(AgentStates).map((value) => value as string),
      },
    ];
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
        prop: 'combined_tags',
        flexGrow: 10,
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
        flexGrow: 4,
        canAutoResize: true,
        name: 'Last Activity',
        sortable: false,
        cellTemplate: this.agentLastActivityTemplateCell,
      },
      {
        name: '',
        prop: 'actions',
        flexGrow: 3,
        minWidth: 150,
        canAutoResize: true,
        sortable: false,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
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
