import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import {
  ColumnMode,
  DatatableComponent,
  TableColumn,
} from '@swimlane/ngx-datatable';
import { SinkDeleteComponent } from 'app/pages/sinks/delete/sink.delete.component';
import {
  Sink,
  SinkBackends,
  SinkStates,
} from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { combineLatest, Observable } from 'rxjs';
import { map, startWith } from 'rxjs/operators';
import { OrbService } from 'app/common/services/orb.service';
import {
  FilterOption,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';
import { FilterService } from 'app/common/services/filter.service';

@Component({
  selector: 'ngx-sink-list-component',
  templateUrl: './sink.list.component.html',
  styleUrls: ['./sink.list.component.scss'],
})
export class SinkListComponent implements AfterViewInit, AfterViewChecked {
  strings = STRINGS.sink;

  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  // templates
  @ViewChild('sinkNameTemplateCell') sinkNameTemplateCell: TemplateRef<any>;

  @ViewChild('sinkStateTemplateCell') sinkStateTemplateCell: TemplateRef<any>;

  @ViewChild('sinkTagsTemplateCell') sinkTagsTemplateCell: TemplateRef<any>;

  @ViewChild('sinkActionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

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
    this.filters$ = this.filters.getFilters().pipe(startWith([]));

    this.filteredSinks$ = combineLatest([this.sinks$, this.filters$]).pipe(
      map(([agents, _filters]) => {
        let filtered = agents;
        _filters.forEach((_filter) => {
          filtered = filtered.filter((value) => {
            const paramValue = _filter.param;
            const result = _filter.filter(value, paramValue);
            return result;
          });
        });

        return filtered;
      }),
    );

    this.filterOptions = [
      {
        name: 'Name',
        prop: 'name',
        filter: (sink: Sink, name: string) => {
          return sink.name?.includes(name);
        },
        type: FilterTypes.Input,
      },
      {
        name: 'Tags',
        prop: 'tags',
        filter: (sink: Sink, tag: string) => {
          const values = Object.entries(sink.tags)
            .map((entry) => `${entry[0]}: ${entry[1]}`);
          return values.reduce((acc, val) => {
            acc = acc || val.includes(tag.trim());
            return acc;
          }, false);
        },
        autoSuggestion: orb.getSinksTags(),
        type: FilterTypes.AutoComplete,
      },
      {
        name: 'Status',
        prop: 'state',
        filter: (sink: Sink, states: string[]) => {
          return states.reduce((prev, cur) => {
            return sink.state === cur || prev;
          }, false);
        },
        type: FilterTypes.MultiSelect,
        options: Object.values(SinkStates).map((value) => value as string),
      },
      {
        name: 'Backend',
        prop: 'backend',
        filter: (sink: Sink, backends: string[]) => {
          return backends.reduce((prev, cur) => {
            return sink.backend === cur || prev;
          }, false);
        },
        type: FilterTypes.MultiSelect,
        options: Object.values(SinkBackends).map((value) => value as string),
      },
    ];
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
