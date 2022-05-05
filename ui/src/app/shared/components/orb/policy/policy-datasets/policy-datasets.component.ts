import {
  AfterViewChecked, AfterViewInit, ChangeDetectorRef, Component, EventEmitter,
  Input, OnChanges, OnDestroy, OnInit, Output, SimpleChanges, TemplateRef,
  ViewChild,
} from '@angular/core';
import { NbDialogService } from '@nebular/theme';
import {
  ColumnMode, DatatableComponent, TableColumn,
} from '@swimlane/ngx-datatable';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import {
  AgentGroupsService,
} from 'app/common/services/agents/agent.groups.service';
import {
  DatasetPoliciesService,
} from 'app/common/services/dataset/dataset.policies.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import {
  DatasetFromComponent,
} from 'app/pages/datasets/dataset-from/dataset-from.component';
import { Subscription } from 'rxjs';
import { concatMap } from 'rxjs/operators';

interface FlexDataset extends Dataset {
  sinks?: Sink[];
  agent_group?: AgentGroup;
}

@Component({
  selector: 'ngx-policy-datasets',
  templateUrl: './policy-datasets.component.html',
  styleUrls: ['./policy-datasets.component.scss'],
})
export class PolicyDatasetsComponent implements OnInit, OnDestroy,
  AfterViewInit, AfterViewChecked, OnChanges {
  @Input()
  policy: AgentPolicy;

  @Output()
  refreshPolicy: EventEmitter<string>;

  datasets: FlexDataset[];

  isLoading: boolean;

  subscription: Subscription;

  errors;

  columnMode = ColumnMode;

  columns: TableColumn[];

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  // templates
  @ViewChild('nameTemplateCell') nameTemplateCell: TemplateRef<any>;

  @ViewChild('groupTemplateCell') groupTemplateCell: TemplateRef<any>;

  @ViewChild('validTemplateCell') validTemplateCell: TemplateRef<any>;

  @ViewChild('sinksTemplateCell') sinksTemplateCell: TemplateRef<any>;

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  constructor(
    private datasetService: DatasetPoliciesService,
    private groupsService: AgentGroupsService,
    private sinksService: SinksService,
    private dialogService: NbDialogService,
    private cdr: ChangeDetectorRef,
  ) {
    this.refreshPolicy = new EventEmitter<string>();
    this.policy = {};
    this.datasets = [];
    this.errors = {};
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes.policy) {
      this.retrieveInfo();
    }
  }

  retrieveInfo() {
    if (this.isLoading) {
      return;
    }
    this.isLoading = true;
    this.subscription = this.retrievePolicyDatasets()
      .pipe(
        concatMap(datasets => this.retrieveAgentGroups()),
        concatMap(sinks => this.retrieveSinks()))
      .subscribe(resp => {
        this.isLoading = false;
        if (this.table) {
          this.table.rows = this.datasets;
        }
        this.cdr.markForCheck();
      });
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        flexGrow: 3,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'agent_group',
        name: 'Agent Group',
        resizeable: false,
        flexGrow: 3,
        cellTemplate: this.groupTemplateCell,
      },
      {
        prop: 'valid',
        name: 'Valid',
        resizeable: false,
        flexGrow: 2,
        cellTemplate: this.validTemplateCell,
      },
      {
        prop: 'sinks',
        name: 'Sinks',
        resizeable: false,
        flexGrow: 9,
        cellTemplate: this.sinksTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  ngAfterViewChecked() {
    if (this.table && this.table.recalculate && (
      this.tableWrapper.nativeElement.clientWidth !== this.currentComponentWidth
    )) {
      this.currentComponentWidth = this.tableWrapper.nativeElement.clientWidth;
      this.table.recalculate();
      this.cdr.detectChanges();
      window.dispatchEvent(new Event('resize'));
    }
  }

  retrievePolicyDatasets() {
    return this.datasetService.getAllDatasets()
      .map(resp => {
        this.datasets = resp.data.filter(
          dataset => dataset.agent_policy_id === this.policy.id);
        if (this.table) {
          this.table.rows = this.datasets;
        }
        return this.datasets;
      });
  }

  // TODO this should be avoided
  retrieveAgentGroups() {
    return this.groupsService.getAllAgentGroups()
      .map(resp => {
        const groups = resp.data;
        this.datasets = this.datasets.map(dataset => {
          dataset.agent_group = groups.find(
            group => group.id === dataset.agent_group_id);
          return dataset;
        });
        if (this.table) {
          this.table.rows = this.datasets;
        }
        return resp;
      });
  }

  retrieveSinks() {
    return this.sinksService.getAllSinks()
      .map(resp => {
        const sinks = resp.data;
        this.datasets = this.datasets.map(dataset => {
          dataset.sinks = dataset.sink_ids.map(
            id => sinks.find(sink => sink.id === id));
          return dataset;
        });
        if (this.table) {
          this.table.rows = this.datasets;
        }
      });
  }

  onCreateDataset() {
    this.dialogService.open(DatasetFromComponent,
      {
        autoFocus: true,
        closeOnEsc: true,
        context: {
          policy: this.policy,
        },
        hasScroll: false,
        hasBackdrop: false,
      }).onClose.subscribe(resp => {
      if (resp === 'created') {
        // this.retrieveInfo();
        this.refreshPolicy.emit('refresh-from-dataset');
      }
    });
  }

  onOpenEdit(dataset) {
    this.dialogService.open(DatasetFromComponent,
      {
        autoFocus: true,
        closeOnEsc: true,
        context: {
          dataset,
        },
        hasScroll: false,
        hasBackdrop: false,
      }).onClose.subscribe(resp => {
      if (resp === 'changed' || 'deleted') {
        // this.retrieveInfo();
        this.refreshPolicy.emit('refresh-from-dataset');
      }
    });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
