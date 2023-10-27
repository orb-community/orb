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
  filterMultiSelect,
  FilterOption, filterString,
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
import { DeleteSelectedComponent } from 'app/shared/components/delete/delete.selected.component';
import { STRINGS } from 'assets/text/strings';
import { Observable, Subscription } from 'rxjs';

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

  selected: any[] = [];

  private groupsSubscription: Subscription;

  // templates
  @ViewChild('agentGroupNameTemplateCell')
  agentGroupNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentGroupTemplateCell')
  agentGroupsTemplateCell: TemplateRef<any>;

  @ViewChild('agentGroupTagsTemplateCell')
  agentGroupTagsTemplateCell: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateCell') checkboxTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateHeader') checkboxTemplateHeader: TemplateRef<any>;

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

  contextMenuRow: any;

  showContextMenu = false;
  menuPositionLeft: number;
  menuPositionTop: number;

  agentGroupContextMenu = [
    {icon: 'search-outline', action: 'openview'},
    {icon: 'edit-outline', action: 'openedit'},
    {icon: 'trash-outline', action: 'opendelete'},
  ];

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
    this.selected = [];

    this.groups$ = this.orb.getGroupListView();

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
        prop: 'tags',
        filter: filterTags,
        autoSuggestion: orb.getGroupsTags(),
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Description',
        prop: 'description',
        filter: filterString,
        type: FilterTypes.Input,
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

  onTableContextMenu(event) {
    event.event.preventDefault();
    event.event.stopPropagation();
    if (event.type === 'body') {
      this.contextMenuRow = {
        objectType: 'group',
        ...event.content,
      };
      this.menuPositionLeft = event.event.clientX;
      this.menuPositionTop = event.event.clientY;
      this.showContextMenu = true;
    }
  }
  handleContextClick() {
    if (this.showContextMenu) {
      this.showContextMenu = false;
    }
  }

  ngOnDestroy(): void {
    if (this.groupsSubscription) {
      this.groupsSubscription.unsubscribe();
    }
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
    this.columns = [
      {
        name: '',
        prop: 'checkbox',
        width: 62,
        minWidth: 62,
        canAutoResize: false,
        resizeable: false,
        sortable: false,
        cellTemplate: this.checkboxTemplateCell,
        headerTemplate: this.checkboxTemplateHeader,
      },
      {
        prop: 'name',
        name: 'Name',
        width: 230,
        canAutoResize: true,
        resizeable: true,
        minWidth: 150,
        cellTemplate: this.agentGroupNameTemplateCell,
      },
      {
        prop: 'matching_agents',
        name: 'Agents',
        width: 150,
        canAutoResize: true,
        resizeable: true,
        minWidth: 80,
        comparator: (a, b) => a.total - b.total,
        cellTemplate: this.agentGroupsTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        width: 350,
        canAutoResize: true,
        resizeable: true,
        minWidth: 180,
        cellTemplate: this.agentGroupNameTemplateCell,
      },
      {
        prop: 'tags',
        width: 460,
        canAutoResize: true,
        resizeable: true,
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
        width: 150,
        resizeable: true,
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
  onOpenDeleteSelected() {
    const selected = this.selected;
    const elementName = 'Agent Groups';
    this.dialogService
      .open(DeleteSelectedComponent, {
        context: { selected, elementName },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.deleteSelectedAgentGroups();
          this.selected = [];
          this.orb.refreshNow();
        }
      });
  }

  deleteSelectedAgentGroups() {
    this.selected.forEach((group) => {
      this.agentGroupsService.deleteAgentGroup(group.id).subscribe();
    });
    this.notificationsService.success('All selected Groups delete requests succeeded', '');
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
  public onCheckboxChange(event: any, row: any): void {
    if (this.getChecked(row) === false) {
      this.selected.push(row);
    } else {
      for (let i = 0; i < this.selected.length; i++) {
        if (this.selected[i].id === row.id) {
          this.selected.splice(i, 1);
          break;
        }
      }
    }
  }

  public getChecked(row: any): boolean {
    const item = this.selected.filter((e) => e.id === row.id);
    return item.length > 0 ? true : false;
  }

  onHeaderCheckboxChange(event: any) {
    if (event.target.checked && this.filteredGroups$) {
      this.groupsSubscription = this.filteredGroups$.subscribe(rows => {
        this.selected = [];
        rows.forEach(row => {
          this.selected.push(row);
        });
      });
    } else {
      if (this.groupsSubscription) {
        this.groupsSubscription.unsubscribe();
      }
      this.selected = [];
    }
  }
}
