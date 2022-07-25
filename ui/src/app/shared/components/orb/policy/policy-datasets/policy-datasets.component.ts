import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  OnInit,
  Output,
  SimpleChanges,
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
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetFromComponent } from 'app/pages/datasets/dataset-from/dataset-from.component';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';
import { Subscription } from 'rxjs';

@Component({
  selector: 'ngx-policy-datasets',
  templateUrl: './policy-datasets.component.html',
  styleUrls: ['./policy-datasets.component.scss'],
})
export class PolicyDatasetsComponent
  implements OnInit, OnDestroy, AfterViewInit, AfterViewChecked, OnChanges {
  @Input()
  datasets: Dataset[];

  @Input()
  policy: AgentPolicy;

  @Output()
  refreshPolicy: EventEmitter<string>;

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
    private dialogService: NbDialogService,
    private cdr: ChangeDetectorRef,
    protected router: Router,
    protected route: ActivatedRoute,
  ) {
    this.refreshPolicy = new EventEmitter<string>();
    this.datasets = [];
    this.errors = {};
  }

  ngOnInit(): void {}

  ngOnChanges(changes: SimpleChanges) {}

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        canAutoResize: true,
        minWidth: 90,
        width: 120,
        maxWidth: 150,
        flexGrow: 3,
        cellTemplate: this.nameTemplateCell,
      },
      {
        prop: 'agent_group',
        name: 'Agent Group',
        resizeable: false,
        canAutoResize: true,
        minWidth: 90,
        width: 120,
        maxWidth: 150,
        flexGrow: 3,
        cellTemplate: this.groupTemplateCell,
      },
      {
        prop: 'valid',
        name: 'Valid',
        resizeable: false,
        canAutoResize: true,
        minWidth: 65,
        width: 80,
        maxWidth: 100,
        flexGrow: 2,
        cellTemplate: this.validTemplateCell,
      },
      {
        prop: 'sinks',
        name: 'Sinks',
        resizeable: false,
        canAutoResize: true,
        minWidth: 200,
        width: 300,
        maxWidth: 500,
        flexGrow: 6,
        cellTemplate: this.sinksTemplateCell,
      },
    ];

    this.cdr.detectChanges();
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

  onCreateDataset() {
    this.dialogService
      .open(DatasetFromComponent, {
        autoFocus: true,
        closeOnEsc: true,
        context: {
          policy: this.policy,
        },
        hasScroll: false,
        hasBackdrop: true,
        closeOnBackdropClick: true,
      })
      .onClose.subscribe((resp) => {
        if (resp === 'created') {
          this.refreshPolicy.emit('refresh-from-dataset');
        }
      });
  }

  onOpenEdit(dataset) {
    this.dialogService
      .open(DatasetFromComponent, {
        autoFocus: true,
        closeOnEsc: false,
        context: {
          dataset,
        },
        hasScroll: false,
        closeOnBackdropClick: true,
        hasBackdrop: true,
      })
      .onClose.subscribe((resp) => {
        if (resp === 'changed' || resp === 'deleted') {
          this.refreshPolicy.emit('refresh-from-dataset');
        }
      });
  }

  onOpenGroupDetails(agentGroup) {
    this.dialogService
      .open(AgentGroupDetailsComponent, {
        autoFocus: true,
        closeOnEsc: true,
        context: { agentGroup },
        hasScroll: false,
        hasBackdrop: false,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEditAgentGroup(agentGroup);
        }
      });
  }

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`/pages/fleet/groups/edit/${agentGroup.id}`], {
      state: { agentGroup: agentGroup, edit: true },
      relativeTo: this.route,
    });
  }

  onOpenSinkDetails(sink) {
    this.dialogService
      .open(SinkDetailsComponent, {
        autoFocus: true,
        closeOnEsc: true,
        context: { sink },
        hasScroll: false,
        hasBackdrop: false,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEditSink(sink);
        }
      });
  }

  onOpenEditSink(sink: any) {
    this.router.navigate([`pages/sinks/edit/${sink.id}`], {
      relativeTo: this.route,
      state: { sink: sink, edit: true },
    });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
