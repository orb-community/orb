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

import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import {
  FilterOption,
  filterSubstr,
  filterTags,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';

import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { FilterService } from 'app/common/services/filter.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { AgentGroupDeleteComponent } from 'app/pages/fleet/groups/delete/agent.group.delete.component';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { STRINGS } from 'assets/text/strings';
import { Observable } from 'rxjs';

@Component({
  selector: 'ngx-agent-group-list-component',
  templateUrl: './agent.group.list.component.html',
  styleUrls: ['./agent.group.list.component.scss'],
})
export class AgentGroupListComponent
  implements AfterViewInit, AfterViewChecked, OnDestroy {
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

    this.filters$ = this.filters.getFilters();

    this.filterOptions = [
      {
        name: 'Name',
        prop: 'name',
        filter: filterSubstr,
        type: FilterTypes.Input,
      },
      {
        name: 'Tags',
        prop: 'tags',
        filter: filterTags,
        autoSuggestion: orb.getGroupsTags(),
        type: FilterTypes.AutoComplete,
      },
      // {
      //   name: 'Status',
      //   prop: 'state',
      //   filter: filterMultiSelect,
      //   type: FilterTypes.MultiSelect,
      //   options: Object.values(AgentStates).map((value) => value as string),
      // },
    ];

    this.filteredGroups$ = this.filters.createFilteredList()(
      this.groups$,
      this.filters$,
      this.filterOptions,
    );
  }

  ngOnDestroy(): void {
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
