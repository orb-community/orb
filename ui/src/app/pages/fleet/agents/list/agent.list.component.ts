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
import { DeleteSelectedComponent } from 'app/shared/components/delete/delete.selected.component';
import { STRINGS } from 'assets/text/strings';
import { Observable, Subscription } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { AgentResetComponent } from '../reset/agent.reset.component';

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

  selected: any[] = [];

  canResetAgents: boolean;

  isResetting: boolean;
  
  private agentsSubscription: Subscription;


  // templates
  @ViewChild('agentNameTemplateCell') agentNameTemplateCell: TemplateRef<any>;

  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

  @ViewChild('agentPolicyStateTemplateCell') agentPolicyStateTemplateRef: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateCell') checkboxTemplateCell: TemplateRef<any>;

  @ViewChild('agentLastActivityTemplateCell')
  agentLastActivityTemplateCell: TemplateRef<any>;

  @ViewChild('agentVersionTemplateCell') agentVersionTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateHeader') checkboxTemplateHeader: TemplateRef<any>;

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
    protected agentsService: AgentsService,
    protected notificationService: NotificationsService,
  ) {
    this.isResetting = false;
    this.selected = [];
    this.agents$ = this.orb.getAgentListView().pipe(
      map(agents => {
        return agents.map(agent => {
          let version: string;
          if (agent.state !== 'new') {
            version = agent.agent_metadata.orb_agent.version;
          } else {
            version = '-';
          }
          return {
            ...agent,
            version,
          };
        });
      })
    );

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
      {
        name: 'Version',
        prop: 'version',
        filter: filterString,
        type: FilterTypes.Input,
      },
    ];

    this.filteredAgents$ = this.filters.createFilteredList()(
      this.agents$,
      this.filters$,
      this.filterOptions,
    );
  }

  ngOnDestroy() {
    if (this.agentsSubscription) {
      this.agentsSubscription.unsubscribe();
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
    this.orb.refreshNow();
    this.columns = [
      {
        name: '',
        prop: 'checkbox',
        flexGrow: 1,
        minWidth: 62,
        canAutoResize: true,
        sortable: false,
        cellTemplate: this.checkboxTemplateCell,
        headerTemplate: this.checkboxTemplateHeader,
      },
      {
        prop: 'name',
        flexGrow: 5,
        canAutoResize: true,
        minWidth: 150,
        name: 'Name',
        cellTemplate: this.agentNameTemplateCell,
      },
      {
        prop: 'state',
        flexGrow: 3,
        canAutoResize: true,
        name: 'Status',
        cellTemplate: this.agentStateTemplateRef,
      },
      {
        prop: 'policy_agg_info',
        flexGrow: 4,
        canAutoResize: true,
        minWidth: 150,
        name: 'Policies',
        cellTemplate: this.agentPolicyStateTemplateRef,
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
        prop: 'version',
        flexGrow: 5,
        minWidth: 150,
        canAutoResize: true,
        name: 'Version',
        sortable: false,
        cellTemplate: this.agentVersionTemplateCell,
      },
      {
        prop: 'ts_last_hb',
        flexGrow: 4,
        minWidth: 150,
        canAutoResize: true,
        name: 'Last Activity',
        sortable: false,
        cellTemplate: this.agentLastActivityTemplateCell,
      },
      {
        name: '',
        prop: 'actions',
        flexGrow: 2.5,
        minWidth: 150,
        canAutoResize: true,
        sortable: false,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
  }


  public onCheckboxChange(event: any, row: any): void { 
    let selectedAgent = {
      id: row.id,
      resetable: true,
      name: row.name,
      state: row.state,
    }
    if (this.getChecked(row) === false) {
      let resetable = true;
      if (row.state === 'new' || row.state === 'offline') {
        resetable = false;
      }
      selectedAgent.resetable = resetable;
      this.selected.push(selectedAgent);
    } else {
      for (let i = 0; i < this.selected.length; i++) {
        if (this.selected[i].id === row.id) {
          this.selected.splice(i, 1);
          break;
        }
      }
    }
    const reset = this.selected.filter((e) => e.resetable === false);
    this.canResetAgents = reset.length > 0 ? true : false;
  }



  public getChecked(row: any): boolean {
    const item = this.selected.filter((e) => e.id === row.id);
    return item.length > 0 ? true : false;
  }

  notifyResetSuccess() {
    this.notificationService.success('All Agents Resets Requested', '');
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
  onOpenDeleteSelected() {
    const selected = this.selected;
    const elementName = "Agents"
    this.dialogService
      .open(DeleteSelectedComponent, {
        context: { selected, elementName },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.deleteSelectedAgents();
          this.selected = [];
          this.orb.refreshNow();
        }
      });
  }

  deleteSelectedAgents() {
    this.selected.forEach((agent) => {
      this.agentService.deleteAgent(agent.id).subscribe();
    })
    this.notificationsService.success('All selected Agents delete requests succeeded', '');
  }

  onOpenResetAgents() {
    const size = this.selected.length;
    this.dialogService
      .open(AgentResetComponent, {
        context: { size },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.resetAgents();
          this.orb.refreshNow();
        }
      })
  }
  resetAgents() {
    if (!this.isResetting) {
      this.isResetting = true;
      this.selected.forEach((agent) => {
        this.agentService.resetAgent(agent.id).subscribe();
      })
      this.notifyResetSuccess();
      this.selected = [];
      this.isResetting = false;
    }
  }

  onHeaderCheckboxChange(event: any) {
    if (event.target.checked && this.filteredAgents$) {
      this.agentsSubscription = this.filteredAgents$.subscribe(rows => {
        this.selected = [];
        rows.forEach(row => {
          const policySelected = {
            id: row.id,
            name: row.name,
            state: row.state,
            resetable: row.state === 'new' || row.state === 'offline' ? false : true,
          }
          this.selected.push(policySelected);
        });
      });
    } else {
      if (this.agentsSubscription) {
        this.agentsSubscription.unsubscribe();
      }
      this.selected = [];
    }
    const reset = this.selected.filter((e) => e.resetable === false);
    this.canResetAgents = reset.length > 0 ? true : false;
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
