import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { ActivatedRoute, Router } from '@angular/router';
import {
  ColumnMode,
  DatatableComponent,
  TableColumn,
} from '@swimlane/ngx-datatable';
import {
  filterMultiSelect,
  FilterOption,
  filterSubstr,
  filterTags,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';
import {
  Sink,
  SinkBackends,
  SinkStates,
} from 'app/common/interfaces/orb/sink.interface';
import { FilterService } from 'app/common/services/filter.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkDeleteComponent } from 'app/pages/sinks/delete/sink.delete.component';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';
import { STRINGS } from 'assets/text/strings';
import { Observable } from 'rxjs';

@Component({
  selector: 'ngx-sink-list-component',
  templateUrl: './sink.list.component.html',
  styleUrls: ['./sink.list.component.scss'],
})
export class SinkListComponent implements AfterViewInit, AfterViewChecked, OnDestroy {
  strings = STRINGS.sink;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  // templates
  @ViewChild('sinkNameTemplateCell') sinkNameTemplateCell: TemplateRef<any>;

  @ViewChild('sinkStateTemplateCell') sinkStateTemplateCell: TemplateRef<any>;

  @ViewChild('sinkTagsTemplateCell') sinkTagsTemplateCell: TemplateRef<any>;

  @ViewChild('sinkActionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  sinks$: Observable<Sink[]>;
  filterOptions: FilterOption[];
  filters$!: Observable<FilterOption[]>;
  filteredSinks$: Observable<Sink[]>;

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private notificationsService: NotificationsService,
    private sinkService: SinksService,
    private route: ActivatedRoute,
    private router: Router,
    private orb: OrbService,
    private filters: FilterService,
  ) {
    this.sinks$ = this.orb.getSinkListView();
    this.filters$ = this.filters.getFilters();

    this.filterOptions = [
      {
        name: 'Name',
        prop: 'name',
        filter: filterSubstr,
        type: FilterTypes.Input,
      },
      {
        name: 'Tags',
        prop: 'tags',
        filter: filterTags,
        autoSuggestion: orb.getSinksTags(),
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Status',
        prop: 'state',
        filter: filterMultiSelect,
        type: FilterTypes.MultiSelect,
        options: Object.values(SinkStates).map((value) => value as string),
      },
      {
        name: 'Backend',
        prop: 'backend',
        filter: filterMultiSelect,
        type: FilterTypes.MultiSelect,
        options: Object.values(SinkBackends).map((value) => value as string),
      },
    ];

    this.filteredSinks$ = this.filters.createFilteredList()(
      this.sinks$,
      this.filters$,
      this.filterOptions,
    );
  }

  ngOnDestroy(): void {
    this.orb.killPolling.next();
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

  ngAfterViewInit() {
    this.orb.refreshNow();
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        canAutoResize: true,
        resizeable: false,
        flexGrow: 2,
        minWidth: 150,
        cellTemplate: this.sinkNameTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: false,
        minWidth: 150,
        flexGrow: 2,
        cellTemplate: this.sinkNameTemplateCell,
      },
      {
        prop: 'backend',
        name: 'Type',
        resizeable: false,
        minWidth: 120,
        flexGrow: 1,
        cellTemplate: this.sinkNameTemplateCell,
      },
      {
        prop: 'state',
        name: 'Status',
        resizeable: false,
        flexGrow: 1,
        cellTemplate: this.sinkStateTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        flexGrow: 2,
        resizeable: false,
        cellTemplate: this.sinkTagsTemplateCell,
        comparator: (a, b) =>
          Object.entries(a)
            .map(([key, value]) => `${key}:${value}`)
            .join(',')
            .localeCompare(
              Object.entries(b)
                .map(([key, value]) => `${key}:${value}`)
                .join(','),
            ),
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 150,
        resizeable: false,
        sortable: false,
        flexGrow: 2,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
  }

  onOpenAdd() {
    this.router.navigate(['add'], { relativeTo: this.route });
  }

  onOpenEdit(sink: any) {
    this.router.navigate([`edit/${sink.id}`], {
      relativeTo: this.route,
      state: { sink: sink, edit: true },
    });
  }

  openDeleteModal(row: any) {
    const { id } = row;
    this.dialogService
      .open(SinkDeleteComponent, {
        context: { sink: row },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.sinkService.deleteSink(id).subscribe(() => {
            this.notificationsService.success('Sink successfully deleted', '');
            this.orb.refreshNow();
          });
        }
      });
  }

  openDetailsModal(row: any) {
    this.dialogService
      .open(SinkDetailsComponent, {
        context: { sink: row },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEdit(row);
        }
      });
  }

  filterByInactive = (sink) => sink.state === 'inactive';
}
