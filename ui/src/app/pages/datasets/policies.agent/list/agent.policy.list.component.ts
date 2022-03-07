import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { STRINGS } from '../../../../../assets/text/strings';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { NbDialogService } from '@nebular/theme';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentPolicyDeleteComponent } from 'app/pages/datasets/policies.agent/delete/agent.policy.delete.component';
import { AgentPolicyDetailsComponent } from 'app/pages/datasets/policies.agent/details/agent.policy.details.component';
import { DatePipe } from '@angular/common';

@Component({
  selector: 'ngx-agent-policy-list-component',
  templateUrl: './agent.policy.list.component.html',
  styleUrls: ['./agent.policy.list.component.scss'],
})
export class AgentPolicyListComponent implements OnInit, AfterViewInit, AfterViewChecked {
  strings = STRINGS.agents;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<AgentPolicy>;

  searchPlaceholder = 'Search by name';

  @ViewChild('nameTemplateCell') nameTemplateCell: TemplateRef<any>;

  @ViewChild('versionTemplateCell') versionTemplateCell: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
      filter: (policy, name) => policy?.name.includes(name),
    },
    {
      id: '1',
      label: 'Description',
      prop: 'description',
      selected: false,
      filter: (policy, description) => policy?.description.includes(description),
    },
    {
      id: '2',
      label: 'Version',
      prop: 'version',
      selected: false,
      filter: (policy, version) => policy?.version.includes(version),
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
    private datePipe: DatePipe,
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.agentPoliciesService.clean();
    this.paginationControls = AgentPoliciesService.getDefaultPagination();
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
    this.agentPoliciesService.clean();
    this.getAgentsPolicies();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Policy Name',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 2,
        minWidth: 120,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: false,
        flexGrow: 4,
        minWidth: 120,
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
        pipe: { transform: (value) => this.datePipe.transform(value, 'M/d/yy, HH:mm z') },
        name: 'Last Modified',
        minWidth: 140,
        flexGrow: 2,
        resizeable: false,
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 100,
        resizeable: false,
        sortable: false,
        flexGrow: 2,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  getAllPolicies(): void {
    this.agentPoliciesService.getAllAgentPolicies().subscribe(resp => {
      this.paginationControls.data = resp.data;
      this.paginationControls.total = resp.data.length;
      this.paginationControls.offset = resp.offset / resp.limit;
      this.loading = false;
      this.cdr.markForCheck();
    });
  }

  getAgentsPolicies(pageInfo: NgxDatabalePageInfo = null): void {
    const finalPageInfo = { ...pageInfo };
    finalPageInfo.dir = 'desc';
    finalPageInfo.order = 'name';
    finalPageInfo.limit = this.paginationControls.limit;
    finalPageInfo.offset = pageInfo?.offset * pageInfo?.limit || 0;

    this.loading = true;
    this.agentPoliciesService.getAgentsPolicies(finalPageInfo).subscribe(
      (resp: OrbPagination<AgentPolicy>) => {
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

  onOpenEdit(agentPolicy: any) {
    this.router.navigate([`edit/${ agentPolicy.id }`], {
      state: { agentPolicy: agentPolicy, edit: true },
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
      this.table.rows = this.paginationControls.data.filter(sink => this.filterValue.split(/[,;]+/gm).reduce((prev, curr) => {
        return this.selectedFilter.filter(sink, curr) && prev;
      }, true));
    }
    this.paginationControls.offset = 0;
  }

  openDeleteModal(row: any) {
    const { name, id } = row as AgentPolicy;
    this.dialogService.open(AgentPolicyDeleteComponent, {
      context: { name },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.agentPoliciesService.deleteAgentPolicy(id).subscribe(() => {
            this.getAgentsPolicies();
            this.notificationsService.success('Agent Policy successfully deleted', '');
          });
        }
      },
    );
  }

  openDetailsModal(agentPolicy: any) {
    this.dialogService.open(AgentPolicyDetailsComponent, {
      context: { agentPolicy: agentPolicy },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEdit(agentPolicy);
      } else {
        this.getAgentsPolicies();
      }
    });
  }
}
