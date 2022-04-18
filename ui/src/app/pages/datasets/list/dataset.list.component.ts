import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { ActivatedRoute, Router } from '@angular/router';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';
import { NbDialogService } from '@nebular/theme';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { DatasetDetailsComponent } from 'app/pages/datasets/details/dataset.details.component';

@Component({
  selector: 'ngx-dataset-list-component',
  templateUrl: './dataset.list.component.html',
  styleUrls: ['./dataset.list.component.scss'],
})
export class DatasetListComponent implements OnInit, AfterViewInit, AfterViewChecked {
  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<Dataset>;

  searchPlaceholder = 'Search by name';

  // templates
  @ViewChild('nameTemplateCell') nameTemplateCell: TemplateRef<any>;

  @ViewChild('validTemplateCell') validTemplateCell: TemplateRef<any>;

  @ViewChild('sinksTemplateCell') sinksTemplateCell: TemplateRef<any>;

  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
      filter: (dataset, name) => dataset?.name.includes(name),
    },
    {
      id: '1',
      label: 'Valid',
      prop: 'valid',
      selected: false,
      filter: (dataset, valid) => `${ dataset?.valid }`.includes(name),
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
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
    private datasetPoliciesService: DatasetPoliciesService,
  ) {
    this.paginationControls = DatasetPoliciesService.getDefaultPagination();
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
    this.getAllDatasets();
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
      {
        name: '',
        prop: 'actions',
        minWidth: 150,
        resizeable: false,
        sortable: false,
        flexGrow: 1,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  getAllDatasets(): void {
    this.datasetPoliciesService.clean();
    this.datasetPoliciesService.getAllDatasets().subscribe(resp => {
      this.paginationControls.data = resp.data;
      this.paginationControls.total = resp.data.length;
      this.paginationControls.offset = resp.offset / resp.limit;
      this.loading = false;
      this.cdr.markForCheck();
    });
  }

  getDatasets(pageInfo: NgxDatabalePageInfo = null): void {
    const finalPageInfo = { ...pageInfo };
    finalPageInfo.dir = 'desc';
    finalPageInfo.order = 'name';
    finalPageInfo.limit = this.paginationControls.limit;
    finalPageInfo.offset = pageInfo?.offset * pageInfo?.limit || 0;

    this.loading = true;
    this.datasetPoliciesService.getDatasetPolicies(pageInfo).subscribe(
      (resp: OrbPagination<Dataset>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo?.offset || 0;
        this.paginationControls.total = resp.total;
        this.loading = false;
      },
    );
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route.parent,
    });
  }

  onOpenEdit(dataset: any) {
    this.router.navigate(
      [`edit/${ dataset.id }`],
      {
        relativeTo: this.route.parent,
        state: { dataset: dataset, edit: true },
      },
    );
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
    const { id } = row;
    this.dialogService.open(DatasetDeleteComponent, {
      context: { name: row.name },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.datasetPoliciesService.deleteDataset(id).subscribe(() => {
            this.notificationsService.success('Dataset successfully deleted', '');
            this.getAllDatasets();
          });
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(DatasetDetailsComponent, {
      context: { dataset: row },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEdit(row);
      } else {
        this.getAllDatasets();
      }
    });
  }
}
