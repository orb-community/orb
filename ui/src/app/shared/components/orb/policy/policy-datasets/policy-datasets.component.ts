import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  Input,
  OnDestroy,
  OnInit, TemplateRef,
  ViewChild,
} from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Subscription } from 'rxjs';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { NbDialogService } from '@nebular/theme';
import { DatasetFromComponent } from 'app/pages/datasets/dataset-from/dataset-from.component';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { concatMap } from 'rxjs/operators';

@Component({
  selector: 'ngx-policy-datasets',
  templateUrl: './policy-datasets.component.html',
  styleUrls: ['./policy-datasets.component.scss'],
})
export class PolicyDatasetsComponent implements OnInit, OnDestroy, AfterViewInit, AfterViewChecked {
  @Input()
  policy: AgentPolicy;

  datasets: Dataset[];

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
    protected datasetService: DatasetPoliciesService,
    protected groupsService: AgentGroupsService,
    protected dialogService: NbDialogService,
    protected cdr: ChangeDetectorRef,
    ) {
    this.policy = {};
    this.datasets = [];
    this.errors = {};
  }

  ngOnInit(): void {
    this.subscription = this.retrievePolicyDatasets()
      .pipe(concatMap(datasets => this.retrieveAgentGroups()))
      .subscribe(resp => {
        this.isLoading = false;
      });
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        flexGrow: 5,
        minWidth: 90,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'agent_group_name',
        name: 'Agent Group Name',
        resizeable: false,
        flexGrow: 5,
        minWidth: 90,
        cellTemplate: this.groupTemplateCell,
      },
      {
        prop: 'valid',
        name: 'Valid',
        resizeable: false,
        flexGrow: 1,
        minWidth: 25,
        cellTemplate: this.validTemplateCell,
      },
      {
        prop: 'sink_ids',
        name: 'Sinks',
        resizeable: false,
        flexGrow: 1,
        minWidth: 25,
        cellTemplate: this.sinksTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  ngAfterViewChecked() {
    if (this.table && this.table.recalculate && (this.tableWrapper.nativeElement.clientWidth !== this.currentComponentWidth)) {
      this.currentComponentWidth = this.tableWrapper.nativeElement.clientWidth;
      this.table.recalculate();
      this.cdr.detectChanges();
      window.dispatchEvent(new Event('resize'));
    }
  }

  retrievePolicyDatasets() {
    return this.datasetService.getAllDatasets()
      .map(resp => {
        this.datasets = resp.data.filter(dataset => dataset.agent_policy_id === this.policy.id);
        return this.datasets;
      });
  }

  // TODO this should be avoided
  retrieveAgentGroups() {
    return this.groupsService.getAllAgentGroups()
      .map(resp => {
        const groups = resp.data;
        this.datasets.forEach(dataset => {
          dataset['agent_group_name'] = groups.find(group => group.id === dataset.agent_group_id).name;
        });
        return resp.data;
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
      }).onClose.subscribe(() => {});
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
