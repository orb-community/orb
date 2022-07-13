import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import {
  ColumnMode,
  DatatableComponent,
  TableColumn,
} from '@swimlane/ngx-datatable';
import { STRINGS } from '../../../../../assets/text/strings';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { NbDialogService } from '@nebular/theme';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentPolicyDeleteComponent } from 'app/pages/datasets/policies.agent/delete/agent.policy.delete.component';
import { DatePipe } from '@angular/common';
import { combineLatest, Observable } from 'rxjs';
import { OrbService } from 'app/common/services/orb.service';
import { map, startWith } from 'rxjs/operators';
import {
  FilterOption,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';
import { FilterService } from 'app/common/services/filter.service';

@Component({
  selector: 'ngx-agent-policy-list-component',
  templateUrl: './agent.policy.list.component.html',
  styleUrls: ['./agent.policy.list.component.scss'],
})
export class AgentPolicyListComponent
  implements AfterViewInit, AfterViewChecked {
  strings = STRINGS.agents;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  @ViewChild('nameTemplateCell') nameTemplateCell: TemplateRef<any>;

  @ViewChild('versionTemplateCell') versionTemplateCell: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  policies$: Observable<AgentPolicy[]>;
  filterOptions: FilterOption[];
  filters$!: Observable<FilterOption[]>;
  filteredPolicies$: Observable<AgentPolicy[]>;

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
    this.policies$ = this.orb.getPolicyListView();
    this.filters$ = this.filters.getFilters().pipe(startWith([]));

    this.filteredPolicies$ = combineLatest([
      this.policies$,
      this.filters$,
    ]).pipe(
      map(([policies, _filters]) => {
        let filtered = policies;
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
        filter: (policy: AgentPolicy, name: string) => {
          return policy.name?.includes(name);
        },
        type: FilterTypes.Input,
      },
      {
        name: 'Description',
        prop: 'description',
        filter: (policy: AgentPolicy, description: string) => {
          return policy.description?.includes(description);
        },
        type: FilterTypes.Input,
      },
      {
        name: 'Version',
        prop: 'version',
        filter: (policy: AgentPolicy, version: number) => {
          return policy.version === version;
        },
        type: FilterTypes.Input,
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
        name: 'Policy Name',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 2,
        minWidth: 150,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: false,
        flexGrow: 4,
        minWidth: 150,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'version',
        name: 'Version',
        resizeable: false,
        flexGrow: 1,
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
        minWidth: 180,
        flexGrow: 3,
        resizeable: false,
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 150,
        resizeable: false,
        sortable: false,
        flexGrow: 2,
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
}
