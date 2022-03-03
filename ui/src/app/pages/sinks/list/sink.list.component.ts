import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { SinkDeleteComponent } from 'app/pages/sinks/delete/sink.delete.component';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-sink-list-component',
  templateUrl: './sink.list.component.html',
  styleUrls: ['./sink.list.component.scss'],
})
export class SinkListComponent implements OnInit, AfterViewInit, AfterViewChecked {
  strings = STRINGS.sink;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<Sink>;

  searchPlaceholder = 'Search by name';


  // templates
  @ViewChild('sinkStateTemplateCell') sinkStateTemplateCell: TemplateRef<any>;

  @ViewChild('sinkTagsTemplateCell') sinkTagsTemplateCell: TemplateRef<any>;

  @ViewChild('sinkActionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
      filter: (sink, name) => sink?.name.includes(name),
    },
    {
      id: '1',
      label: 'Tags',
      prop: 'tags',
      selected: false,
      filter: (sink, tag) => Object.entries(sink?.tags)
        .filter(([key, value]) => key?.includes(tag) || (!!value && value as string).includes(tag)).length > 0,
    },
    {
      id: '2',
      label: 'Description',
      prop: 'description',
      selected: false,
      filter: (sink, description) => sink?.description.includes(description),
    },
    {
      id: '3',
      label: 'Type',
      prop: 'backend',
      selected: false,
      filter: (sink, backend) => sink?.backend.includes(backend),
    },
    {
      id: '4',
      label: 'Status',
      prop: 'state',
      selected: false,
      filter: (sink, state) => sink?.state.includes(state),
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
    private sinkService: SinksService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.sinkService.clean();
    this.paginationControls = SinksService.getDefaultPagination();
  }

  ngAfterViewChecked() {
    if (this.table && this.table.recalculate && (this.tableWrapper.nativeElement.clientWidth !== this.currentComponentWidth)) {
      this.currentComponentWidth = this.tableWrapper.nativeElement.clientWidth;
      this.table.recalculate();
    }
  }

  ngOnInit() {
    this.sinkService.clean();
    this.getAllSinks();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        canAutoResize: true,
        resizeable: false,
        flexGrow: 2,
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
        prop: 'backend',
        name: 'Type',
        resizeable: false,
        minWidth: 100,
        flexGrow: 1,
      },
      {
        prop: 'state',
        name: 'Status',
        resizeable: false,
        minWidth: 100,
        flexGrow: 1,
        cellTemplate: this.sinkStateTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        minWidth: 90,
        flexGrow: 3,
        resizeable: false,
        cellTemplate: this.sinkTagsTemplateCell,
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

  getAllSinks(): void {
    this.sinkService.getAllSinks().subscribe(resp => {
      this.paginationControls.data = resp.data;
      this.paginationControls.total = resp.data.length;
      this.paginationControls.offset = resp.offset / resp.limit;
      this.loading = false;
      this.cdr.markForCheck();
    });
  }

  getSinks(pageInfo: NgxDatabalePageInfo = null): void {
    const isFilter = this.paginationControls.name?.length > 0 || this.paginationControls.tags?.length > 0;
    const finalPageInfo = { ...pageInfo };
    finalPageInfo.dir = 'desc';
    finalPageInfo.order = 'name';
    finalPageInfo.limit = this.paginationControls.limit;
    finalPageInfo.offset = pageInfo?.offset * pageInfo?.limit || 0;

    if (isFilter) {
      if (this.paginationControls.name?.length > 0) finalPageInfo.name = this.paginationControls.name;
      if (this.paginationControls.tags?.length > 0) finalPageInfo.tags = this.paginationControls.tags;
    }

    this.loading = true;
    this.sinkService.getSinks(finalPageInfo, isFilter).subscribe(
      (resp: OrbPagination<Sink>) => {
        this.paginationControls.data = resp.data.slice(resp.offset, resp.offset + resp.limit);
        this.paginationControls.total = resp.total;
        this.paginationControls.offset = resp.offset / resp.limit;
        this.loading = false;
        this.cdr.detectChanges();
      },
    );
  }

  onOpenAdd() {
    this.router.navigate(
      ['add'],
      { relativeTo: this.route },
    );
  }

  onOpenEdit(sink: any) {
    this.router.navigate(
      [`edit/${ sink.id }`],
      {
        relativeTo: this.route,
        state: { sink: sink, edit: true },
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
      this.table.rows = this.paginationControls.data.filter(sink => this.filterValue.split(/[\s,.:]+/gm).reduce((prev, curr) => {
        return this.selectedFilter.filter(sink, curr) && prev;
      }, true));
    }

  }

  openDeleteModal(row: any) {
    const { id } = row;
    this.dialogService.open(SinkDeleteComponent, {
      context: { sink: row },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.sinkService.deleteSink(id).subscribe(() => {
            this.getSinks();
            this.notificationsService.success('Sink successfully deleted', '');
          });
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(SinkDetailsComponent, {
      context: { sink: row },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEdit(row);
      } else {
        this.getSinks();
      }
    });
  }

  filterByInactive = (sink) => sink.state === 'inactive';
}
