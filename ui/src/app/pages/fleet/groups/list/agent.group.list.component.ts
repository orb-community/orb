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

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import {
  FilterOption,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';

import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { FilterService } from 'app/common/services/filter.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { AgentGroupDeleteComponent } from 'app/pages/fleet/groups/delete/agent.group.delete.component';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { STRINGS } from 'assets/text/strings';
import { OrbService } from 'app/common/services/orb.service';
import { combineLatest, Observable } from 'rxjs';
import { filter, map, startWith } from 'rxjs/operators';

@Component({
  selector: 'ngx-agent-group-list-component',
  templateUrl: './agent.group.list.component.html',
  styleUrls: ['./agent.group.list.component.scss'],
})
export class AgentGroupListComponent
  implements AfterViewInit, AfterViewChecked {
  strings = STRINGS.agentGroups;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  searchPlaceholder = 'Search by name';

  // templates
  @ViewChild('agentGroupNameTemplateCell')
  agentGroupNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentGroupTemplateCell')
  agentGroupsTemplateCell: TemplateRef<any>;

  @ViewChild('agentGroupTagsTemplateCell')
  agentGroupTagsTemplateCell: TemplateRef<any>;

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
      filter: (agent, tag) =>
        Object.entries(agent?.tags).filter(([key, value]) =>
          `${key}:${value}`.includes(tag.replace(' ', '')),
        ).length > 0,
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

  groups$: Observable<AgentGroup[]>;
  filterOptions: FilterOption[];
  filters$!: Observable<FilterOption[]>;
  filteredGroups$: Observable<AgentGroup[]>;

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private agentGroupsService: AgentGroupsService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
    private orb: OrbService,
    private filters: FilterService,
  ) {
    this.groups$ = this.orb.getGroupListView();

    this.filters$ = this.filters.getFilters().pipe(startWith([]));

    this.filteredGroups$ = combineLatest([this.groups$, this.filters$]).pipe(
      map(([groups, _filters]) => {
        let filtered = groups;
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
        name: 'Tags',
        prop: 'tags',
        filter: (group: AgentGroup, tag: string) => {
          const values = Object.entries(group.tags)
            .map((entry) => `${entry[0]}: ${entry[1]}`);
          return values.reduce((acc, val) => {
            acc = acc || val.includes(tag.trim());
            return acc;
          }, false);
        },
        autoSuggestion: orb.getGroupsTags(),
        type: FilterTypes.AutoComplete,
      },
      // {
      //   name: 'Status',
      //   prop: 'state',
      //   filter: (agent: Agent, states: string[]) => {
      //     return states.reduce((prev, cur) => {
      //       return agent.state === cur || prev;
      //     }, false);
      //   },
      //   type: FilterTypes.MultiSelect,
      //   options: Object.values(AgentStates).map((value) => value as string),
      // },
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
        name: 'Name',
        flexGrow: 1,
        canAutoResize: true,
        resizeable: false,
        minWidth: 150,
        cellTemplate: this.agentGroupNameTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        flexGrow: 2,
        canAutoResize: true,
        resizeable: false,
        minWidth: 180,
        cellTemplate: this.agentGroupNameTemplateCell,
      },
      {
        prop: 'matching_agents',
        name: 'Agents',
        flexGrow: 1,
        canAutoResize: true,
        resizeable: false,
        minWidth: 80,
        comparator: (a, b) => a.total - b.total,
        cellTemplate: this.agentGroupsTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        flexGrow: 3,
        canAutoResize: true,
        resizeable: false,
        cellTemplate: this.agentGroupTagsTemplateCell,
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
        name: '',
        prop: 'actions',
        flexGrow: 2,
        canAutoResize: true,
        minWidth: 150,
        sortable: false,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(agentGroup: any) {
    this.router.navigate([`edit/${agentGroup.id}`], {
      state: { agentGroup: agentGroup, edit: true },
      relativeTo: this.route,
    });
  }

  onFilterSelected(selFilter) {
    this.searchPlaceholder = `Search by ${selFilter.label}`;
    this.filterValue = null;
  }

  applyFilter() {
    if (this.table.count === 0) return;

    if (!this.filterValue || this.filterValue === '') {
      this.table.rows = this.groups$;
    } else {
      this.table.rows = this.groups$.pipe(
        filter((sink) =>
          this.filterValue.split(/[,;]+/gm).reduce((prev, curr) => {
            return this.selectedFilter.filter(sink, curr) && prev;
          }, true),
        ),
      );
    }
  }

  openDeleteModal(row: any) {
    const { name, id } = row;
    this.dialogService
      .open(AgentGroupDeleteComponent, {
        context: { name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.agentGroupsService.deleteAgentGroup(id).subscribe(() => {
            this.notificationsService.success(
              'Agent Group successfully deleted',
              '',
            );
          });
          this.orb.refreshNow();
        }
      });
  }

  openDetailsModal(row: any) {
    this.dialogService
      .open(AgentGroupDetailsComponent, {
        context: { agentGroup: row },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEdit(row);
        }
      });
  }

  onMatchingAgentsModal(row: any) {
    this.dialogService.open(AgentMatchComponent, {
      context: { agentGroup: row },
      autoFocus: true,
      closeOnEsc: true,
    });
  }
}
