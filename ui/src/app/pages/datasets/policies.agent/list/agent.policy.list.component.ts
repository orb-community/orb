import { DatePipe } from '@angular/common';
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
import { AgentPolicy, AgentPolicyUsage } from 'app/common/interfaces/orb/agent.policy.interface';
import {
  filterNumber,
  FilterOption, filterString, filterTags,
  FilterTypes,
  filterMultiSelect,
} from 'app/common/interfaces/orb/filter-option';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { FilterService } from 'app/common/services/filter.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { AgentPolicyDeleteComponent } from 'app/pages/datasets/policies.agent/delete/agent.policy.delete.component';
import { DeleteSelectedComponent } from 'app/shared/components/delete/delete.selected.component';
import { combineLatest, Observable, Subscription } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { STRINGS } from '../../../../../assets/text/strings';
import { PolicyDuplicateComponent } from '../duplicate/agent.policy.duplicate.confirmation';


@Component({
  selector: 'ngx-agent-policy-list-component',
  templateUrl: './agent.policy.list.component.html',
  styleUrls: ['./agent.policy.list.component.scss'],
})
export class AgentPolicyListComponent
  implements AfterViewInit, AfterViewChecked, OnDestroy {
  strings = STRINGS.agents;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  selected: any[] = [];

  private policiesSubscription: Subscription;

  @ViewChild('nameTemplateCell') nameTemplateCell: TemplateRef<any>;

  @ViewChild('versionTemplateCell') versionTemplateCell: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('usageStateTemplateCell') usageStateTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateCell') checkboxTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateHeader') checkboxTemplateHeader: TemplateRef<any>;

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  @ViewChild('tagsTemplateCell') tagsTemplateCell: TemplateRef<any>;

  private currentComponentWidth;

  policies$: Observable<AgentPolicy[]>;
  filterOptions: FilterOption[];
  filters$!: Observable<FilterOption[]>;
  filteredPolicies$: Observable<AgentPolicy[]>;

  contextMenuRow: any;

  showContextMenu = false;
  menuPositionLeft: number;
  menuPositionTop: number;

  policyContextMenu = [
    {icon: 'search-outline', action: 'openview'},
    {icon:'edit-outline', action: 'openview'},
    {icon: 'copy-outline', action: 'openduplicate'},
    {icon: 'trash-outline', action: 'opendelete'},
  ];

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private datePipe: DatePipe,
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
    private orb: OrbService,
    private filters: FilterService,
  ) {
    this.filters$ = this.filters.getFilters();
    this.selected = [];
    this.policies$ = combineLatest([
      this.orb.getPolicyListView(),
      this.orb.getDatasetListView(),
    ]).pipe(
      filter(([policies, datasets]) => policies !== undefined && policies !== null && datasets !== undefined && datasets !== null),
      map(([policies, datasets]) => {
        return policies.map((policy) => {
          const dataset = datasets.filter((d) => d.valid && d.agent_policy_id === policy.id);
          return { ...policy, policy_usage: dataset.length > 0 ? AgentPolicyUsage.inUse : AgentPolicyUsage.notInUse };
        });
      }),
    );

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
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Version',
        prop: 'version',
        filter: filterNumber,
        type: FilterTypes.Number,
      },
      {
        name: 'Description',
        prop: 'description',
        filter: filterString,
        type: FilterTypes.Input,
      },
      {
        name: 'Usage',
        prop: 'policy_usage',
        filter: filterMultiSelect,
        type: FilterTypes.MultiSelect,
        options: Object.values(AgentPolicyUsage).map((value) => value as string),
      },
    ];

    this.filteredPolicies$ = this.filters.createFilteredList()(
      this.policies$,
      this.filters$,
      this.filterOptions,
    );
  }

  onTableContextMenu(event) {
    event.event.preventDefault();
    event.event.stopPropagation();
    if (event.type === 'body') {
      this.contextMenuRow = {
        objectType: 'policy',
        ...event.content
      }
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
  
  onOpenDuplicatePolicy(agentPolicy: any) {
    const policy = agentPolicy.name;
    this.dialogService
      .open(PolicyDuplicateComponent, {
        context: { policy },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.duplicatePolicy(agentPolicy);
        }
      });
  }
  duplicatePolicy(agentPolicy: any) {
    this.agentPoliciesService
    .duplicateAgentPolicy(agentPolicy.id)
    .subscribe((newAgentPolicy) => {
      if (newAgentPolicy?.id) {
        this.notificationsService.success(
          'Agent Policy Duplicated',
          `New Agent Policy Name: ${newAgentPolicy?.name}`,
        );
        this.router.navigate([`view/${newAgentPolicy.id}`], {
          relativeTo: this.route,
        });
      }
    });
  }

  ngOnDestroy(): void {
    if (this.policiesSubscription) {
      this.policiesSubscription.unsubscribe();
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
        width: 1,
        minWidth: 62,
        canAutoResize: true,
        sortable: false,
        cellTemplate: this.checkboxTemplateCell,
        headerTemplate: this.checkboxTemplateHeader,
      },
      {
        prop: 'name',
        name: 'Policy Name',
        resizeable: true,
        canAutoResize: true,
        width: 230,
        minWidth: 100,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'policy_usage',
        name: 'Usage',
        resizeable: true,
        canAutoResize: true,
        width: 130,
        minWidth: 100,
        cellTemplate: this.usageStateTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: true,
        width: 280,
        minWidth: 100,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'tags',
        width: 170,
        canAutoResize: true,
        name: 'Tags',
        minWidth: 150,
        cellTemplate: this.tagsTemplateCell,
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
        name: 'Version',
        resizeable: true,
        width: 100,
        minWidth: 50,
        cellTemplate: this.versionTemplateCell,
      },
      {
        prop: 'ts_last_modified',
        pipe: {
          transform: (value) =>
            this.datePipe.transform(value, 'M/d/yy, HH:mm z'),
        },
        name: 'Last Modified',
        minWidth: 110,
        width: 150,
        resizeable: true,
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 200,
        resizeable: true,
        sortable: false,
        width: 150,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(agentPolicy: any) {
    this.router.navigate([`edit/${agentPolicy.id}`], {
      state: { agentPolicy: agentPolicy, edit: true },
      relativeTo: this.route,
    });
  }

  onOpenView(agentPolicy: any) {
    this.router.navigate([`view/${agentPolicy.id}`], {
      relativeTo: this.route,
    });
  }

  openDeleteModal(row: any) {
    const { name: name, id } = row as AgentPolicy;
    this.dialogService
      .open(AgentPolicyDeleteComponent, {
        context: { name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.agentPoliciesService.deleteAgentPolicy(id).subscribe(() => {
            this.notificationsService.success(
              'Agent Policy successfully deleted',
              '',
            );
            this.orb.refreshNow();
          });
        }
      });
  }
  onOpenDeleteSelected() {
    const elementName = 'Policies';
    const selected = this.selected;
    this.dialogService
      .open(DeleteSelectedComponent, {
        context: { selected, elementName },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.deleteSelectedAgentsPolicy();
          this.selected = [];
          this.orb.refreshNow();
        }
      });
  }

  deleteSelectedAgentsPolicy() {
    this.selected.forEach((policy) => {
      this.agentPoliciesService.deleteAgentPolicy(policy.id).subscribe();
    });
    this.notificationsService.success('All selected Policies delete requests succeeded', '');
  }
  public onCheckboxChange(event: any, row: any): void {
    const policySelected = {
      id: row.id,
      name: row.name,
      usage: row.policy_usage,
    };
    if (this.getChecked(row) === false) {
      this.selected.push(policySelected);
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
    if (event.target.checked && this.filteredPolicies$) {
      this.policiesSubscription = this.filteredPolicies$.subscribe(rows => {
        this.selected = [];
        rows.forEach(row => {
          const policySelected = {
            id: row.id,
            name: row.name,
            usage: row.policy_usage,
          };
          this.selected.push(policySelected);
        });
      });
    } else {
      if (this.policiesSubscription) {
        this.policiesSubscription.unsubscribe();
      }
      this.selected = [];
    }
  }
}
