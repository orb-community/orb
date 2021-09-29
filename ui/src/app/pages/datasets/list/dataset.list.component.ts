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
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Debounce } from 'app/shared/decorators/utils';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { ActivatedRoute, Router } from '@angular/router';

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

  filterSelectedIndex = '0';

  // templates
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
      label: 'Tags',
      prop: 'tags',
      selected: false,
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  constructor(
    private cdr: ChangeDetectorRef,
    private route: ActivatedRoute,
    private router: Router,
    private datasetPoliciesService: DatasetPoliciesService,
  ) {
    this.datasetPoliciesService.clean();
    this.paginationControls = SinksService.getDefaultPagination();
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
    this.datasetPoliciesService.clean();
    this.getDatasets();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        flexGrow: 1,
        minWidth: 90,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: false,
        minWidth: 100,
        flexGrow: 2,
      },
      {
        prop: 'ts_created',
        name: 'Date Created',
        minWidth: 90,
        flexGrow: 2,
        resizeable: false,
      },
      {
        prop: 'ts_created',
        name: 'Date Last Received',
        minWidth: 90,
        flexGrow: 2,
        resizeable: false,
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


  @Debounce(500)
  getDatasets(pageInfo: NgxDatabalePageInfo = null): void {
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
    this.datasetPoliciesService.getDatasetPolicies(pageInfo, isFilter).subscribe(
      (resp: OrbPagination<Dataset>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo.offset;
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

  onOpenEdit(dataset: any) {
  }

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${ this.tableFilters[selectedIndex].label }`;
  }

  openDeleteModal(row: any) {
  }

  openDetailsModal(row: any) {
  }

  searchDatasetItemByName(input) {
    this.getDatasets({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }

  filterByInactive = (sink) => sink.state === 'inactive';
}
