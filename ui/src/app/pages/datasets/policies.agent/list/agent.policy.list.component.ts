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
import { Debounce } from 'app/shared/decorators/utils';
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

  filterSelectedIndex = '0';

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
    },
    {
      id: '1',
      label: 'Version',
      prop: 'version',
      selected: false,
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
        flexGrow: 3,
        minWidth: 90,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: false,
        flexGrow: 4,
        minWidth: 180,
      },
      {
        prop: 'version',
        name: 'Version',
        resizeable: false,
        flexGrow: 2,
        minWidth: 60,
      },
      {
        prop: 'ts_last_modified',
        pipe: {transform: (value) => this.datePipe.transform(value, 'MMM d, y, HH:mm:ss z')},
        name: 'Last Modified',
        minWidth: 90,
        flexGrow: 2,
        resizeable: false,
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

  @Debounce(500)
  getAgentsPolicies(pageInfo: NgxDatabalePageInfo = null): void {
    const isFilter = this.paginationControls.name?.length > 0 || this.paginationControls.tags?.length > 0;

    if (isFilter) {
      pageInfo = {
        offset: this.paginationControls.offset,
        limit: this.paginationControls.limit,
      };
      if (this.paginationControls.name?.length > 0) pageInfo.name = this.paginationControls.name;
      if (this.paginationControls.tags?.length > 0) pageInfo.tags = this.paginationControls.tags;
    }

    this.loading = true;
    this.agentPoliciesService.getAgentsPolicies(pageInfo, isFilter).subscribe(
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

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${ this.tableFilters[selectedIndex].label }`;
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

  searchAgentByName(input) {
    this.getAgentsPolicies({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }
}
