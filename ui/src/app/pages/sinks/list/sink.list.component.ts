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
  FilterOption, filterString,
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
import { Observable, Subscription } from 'rxjs';
import { DeleteSelectedComponent } from 'app/shared/components/delete/delete.selected.component';

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

  selected: any[] = [];

  private sinksSubscription: Subscription;

  // templates
  @ViewChild('sinkNameTemplateCell') sinkNameTemplateCell: TemplateRef<any>;

  @ViewChild('sinkStateTemplateCell') sinkStateTemplateCell: TemplateRef<any>;

  @ViewChild('sinkTagsTemplateCell') sinkTagsTemplateCell: TemplateRef<any>;

  @ViewChild('sinkActionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateCell') checkboxTemplateCell: TemplateRef<any>;

  @ViewChild('checkboxTemplateHeader') checkboxTemplateHeader: TemplateRef<any>;

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

  contextMenuRow: any;

  showContextMenu: boolean;
  menuPositionLeft: number;
  menuPositionTop: number;

  sinkContextMenu = [
    {icon: 'search-outline', action: 'openview'},
    {icon: 'edit-outline', action: 'openview'},
    {icon: 'trash-outline', action: 'opendelete'},
  ];

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
    this.selected = [];
    this.sinks$ = this.orb.getSinkListView();
    this.filters$ = this.filters.getFilters();

    this.filterOptions = [
      {
        name: 'Name',
        prop: 'name',
        filter: filterString,
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
      {
        name: 'Description',
        prop: 'description',
        filter: filterString,
        type: FilterTypes.Input,
      },
    ];

    this.filteredSinks$ = this.filters.createFilteredList()(
      this.sinks$,
      this.filters$,
      this.filterOptions,
    );
    this.showContextMenu = false;
  }

  ngOnDestroy(): void {
    if (this.sinksSubscription) {
      this.sinksSubscription.unsubscribe();
    }
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
    this.columns = [
      {
        name: '',
        prop: 'checkbox',
        width: 62,
        minWidth: 62,
        canAutoResize: false,
        resizeable: false,
        sortable: false,
        cellTemplate: this.checkboxTemplateCell,
        headerTemplate: this.checkboxTemplateHeader,
      },
      {
        prop: 'name',
        name: 'Name',
        canAutoResize: true,
        resizeable: true,
        width: 250,
        minWidth: 150,
        cellTemplate: this.sinkNameTemplateCell,
      },
      {
        prop: 'state',
        name: 'Status',
        resizeable: true,
        width: 160,
        cellTemplate: this.sinkStateTemplateCell,
      },
      {
        prop: 'backend',
        name: 'Backend',
        resizeable: true,
        minWidth: 120,
        width: 160,
        cellTemplate: this.sinkNameTemplateCell,
      },
      {
        prop: 'description',
        name: 'Description',
        resizeable: true,
        minWidth: 150,
        width: 350,
        cellTemplate: this.sinkNameTemplateCell,
      },
      {
        prop: 'tags',
        name: 'Tags',
        width: 370,
        resizeable: true,
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
        minWidth: 140,
        resizeable: true,
        sortable: false,
        width: 140,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
  }

  onTableContextMenu(event) {
    event.event.preventDefault();
    event.event.stopPropagation();
    if (event.type === 'body') {
      this.contextMenuRow = {
        objectType: 'sink',
        ...event.content,
      };
      this.menuPositionLeft = event.event.clientX;
      this.menuPositionTop = event.event.clientY;
      this.showContextMenu = true;
    }
  }
  handleContextClick() {
    if (this.showContextMenu) {
      this.showContextMenu = false;
    }
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

  onOpenView(sink: any) {
    this.router.navigate([`view/${sink.id}`], {
      relativeTo: this.route,
      state: { sink: sink },
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
  onOpenDeleteSelected() {
    const selected = this.selected;
    const elementName = 'Sinks';
    this.dialogService
      .open(DeleteSelectedComponent, {
        context: { selected, elementName },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.deleteSelectedSinks();
          this.selected = [];
          this.orb.refreshNow();
        }
      });
  }

  deleteSelectedSinks() {
    this.selected.forEach((sink) => {
      this.sinkService.deleteSink(sink.id).subscribe();
    });
    this.notificationsService.success('All selected Sinks delete requests succeeded', '');
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

  public onCheckboxChange(event: any, row: any): void {
    const sinkSelected = {
      id: row.id,
      name: row.name,
      state: row.state,
    };
    if (this.getChecked(row) === false) {
      this.selected.push(sinkSelected);
    } else {
      for (let i = 0; i < this.selected.length; i++) {
        if (this.selected[i].id === row.id) {
          this.selected.splice(i, 1);
          break;
        }
      }
    }
  }

  public getChecked(row: any): boolean {
    const item = this.selected.filter((e) => e.id === row.id);
    return item.length > 0 ? true : false;
  }
  onHeaderCheckboxChange(event: any) {
    if (event.target.checked && this.filteredSinks$)  {
      this.sinksSubscription = this.filteredSinks$.subscribe(rows => {
        this.selected = [];
        rows.forEach(row => {
          const sinkSelected = {
            id: row.id,
            name: row.name,
            state: row.state,
          };
          this.selected.push(sinkSelected);
        });
      });
    } else {
      if (this.sinksSubscription) {
        this.sinksSubscription.unsubscribe();
      }
      this.selected = [];
    }
  }
}
