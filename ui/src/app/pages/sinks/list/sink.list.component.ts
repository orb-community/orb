import { AfterViewInit, ChangeDetectorRef, Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Debounce } from 'app/shared/decorators/utils';
import { SinkDeleteComponent } from 'app/pages/sinks/delete/sink.delete.component';

@Component({
  selector: 'ngx-sink-list-component',
  templateUrl: './sink.list.component.html',
  styleUrls: ['./sink.list.component.scss'],
})
export class SinkListComponent implements OnInit, AfterViewInit {
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

  @ViewChild('sinkStateTemplateCell') sinkStateTemplateCell: TemplateRef<any>;
  @ViewChild('sinkTagsTemplateCell') sinkTagsTemplateCell: TemplateRef<any>;
  @ViewChild('sinkActionsTemplateCell') actionsTemplateCell: TemplateRef<any>;


  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private sinkService: SinksService,
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
    this.router.navigate(
      ['add'],
      {relativeTo: this.route},
    );
  }

  onOpenEdit(sink: any) {
    this.router.navigate(
      [`edit/${sink.id}`],
      {
        relativeTo: this.route,
        state: {sink: sink, edit: true},
      },
    );
  }

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${this.tableFilters[selectedIndex].label}`;
  }

  openDeleteModal(row: any) {
    const {id} = row;
    this.dialogService.open(SinkDeleteComponent, {
      context: {sink: row},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.sinkService.deleteSink(id).subscribe(() => this.getSinks());
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(SinkDetailsComponent, {
      context: {sink: row},
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

  searchSinkItemByName(input) {
    this.getSinks({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }

  filterByInactive = (sink) => sink.state === 'inactive';
}
