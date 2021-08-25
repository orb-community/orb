import {AfterViewInit, ChangeDetectorRef, Component, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {NbDialogService} from '@nebular/theme';

import {
  DropdownFilterItem,
} from 'app/common/interfaces/mainflux.interface';
import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {SinksService} from 'app/common/services/sinks/sinks.service';
import {SinksDetailsComponent} from 'app/pages/sinks/details/sinks.details.component';
import {SinksDeleteComponent} from 'app/pages/sinks/delete/sinks.delete.component';
import {ActivatedRoute, Router} from '@angular/router';
import {STRINGS} from 'assets/text/strings';
import {ColumnMode, TableColumn} from '@swimlane/ngx-datatable';
import {NgxDatabalePageInfo, OrbPagination} from 'app/common/interfaces/orb/pagination';
import {AgentGroup} from 'app/common/interfaces/orb/agent.group.interface';
import {Debounce} from 'app/shared/decorators/utils';

@Component({
  selector: 'ngx-sinks-component',
  templateUrl: './sinks.component.html',
  styleUrls: ['./sinks.component.scss'],
})
export class SinksComponent implements OnInit, AfterViewInit {
  strings = STRINGS.sink;

  columnMode = ColumnMode;
  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<AgentGroup>;

  searchPlaceholder = 'Search by name';
  filterSelectedIndex = '0';

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

  // templates

  @ViewChild('sinkStatusTemplateCell') sinkStatusTemplateCell: TemplateRef<any>;
  @ViewChild('sinkTagsTemplateCell') sinkTagsTemplateCell: TemplateRef<any>;
  @ViewChild('sinkActionsTemplateCell') actionsTemplateCell: TemplateRef<any>;


  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private sinkService: SinksService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.sinkService.clean();
    this.paginationControls = SinksService.getDefaultPagination();
  }

  ngOnInit() {
    this.sinkService.clean();
    this.getSinks();
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
        prop: 'backend',
        name: 'Type',
        resizeable: false,
        minWidth: 100,
        flexGrow: 1,
      },
      {
        prop: 'status',
        name: 'Status',
        resizeable: false,
        minWidth: 100,
        flexGrow: 1,
        cellTemplate: this.sinkStatusTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        minWidth: 90,
        flexGrow: 3,
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


  @Debounce(400)
  getSinks(pageInfo: NgxDatabalePageInfo = null): void {
    const isFilter = pageInfo === null;
    if (isFilter) {
      pageInfo = {
        offset: this.paginationControls.offset,
        limit: this.paginationControls.limit,
      };
      if (this.paginationControls.name?.length > 0) pageInfo.name = this.paginationControls.name;
      if (this.paginationControls.tags?.length > 0) pageInfo.tags = this.paginationControls.tags;
    }

    this.loading = true;
    this.sinkService.getSinks(pageInfo, isFilter).subscribe(
      (resp: OrbPagination<AgentGroup>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo.offset;
        this.loading = false;
      },
    );
  }

  onOpenAdd() {
    this.router.navigate(['../sinks/add'], {
      relativeTo: this.route,
      state: {edit: false},
    });
  }

  onOpenEdit(row: any) {
    this.router.navigate(['../sinks/edit'], {
      relativeTo: this.route,
      queryParams: {id: row.id},
      state: {sink: row, edit: true},
    });
  }

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${this.tableFilters[selectedIndex].label}`;
  }

  openDeleteModal(row: any) {
    const {name, id} = row;
    this.dialogService.open(SinksDeleteComponent, {
      context: {sink: {name, id}},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.sinkService.deleteSink(row.id).subscribe(
            () => {
              // this.page.rows = this.page.rows.filter((u: User) => u.id !== row.id);
              this.notificationsService.success('Sink Item successfully deleted', '');
            },
          );
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(SinksDetailsComponent, {
      context: {sink: row},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getSinks();
        }
      },
    );
  }

  searchSinkItemByName(input) {
    this.getSinks({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }

  filterByInactive = (sink) => sink.status === 'inactive';
}
