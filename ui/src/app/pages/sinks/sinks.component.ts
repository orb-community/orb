import { Component, OnInit } from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import {
  DropdownFilterItem,
  PageFilters,
  TableConfig,
  TablePage,
  User,
} from 'app/common/interfaces/mainflux.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinksDetailsComponent } from 'app/pages/sinks/details/sinks.details.component';
import { SinksDeleteComponent } from 'app/pages/sinks/delete/sinks.delete.component';
import { Router } from '@angular/router';

const defFreq: number = 100;

@Component({
  selector: 'ngx-sinks-component',
  templateUrl: './sinks.component.html',
  styleUrls: ['./sinks.component.scss'],
})
export class SinksComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['Name', 'Description', 'Type', 'Status', 'Tags', 'orb-sink-add'],
    keys: ['name', 'description', 'type', 'status', 'tags', 'orb-action-hover'],
  };
  page: TablePage = {
    limit: 10,
  };
  pageFilters: PageFilters = {
    offset: 0,
    order: 'id',
    dir: 'desc',
    name: '',
  };
  tableFilters: DropdownFilterItem[];

  searchFreq = 0;

  constructor(
    private dialogService: NbDialogService,
    private sinkService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
  ) {
    this.tableFilters = this.tableConfig.colNames.map((name, index) => ({
      id: index.toString(),
      name,
      order: 'asc',
      selected: false,
    })).filter((filter) => (!filter.name.startsWith('orb-')));
  }

  ngOnInit() {
    // Fetch all sinks
    this.getSinks();
  }

  getSinks(name?: string): void {
    this.pageFilters.name = name;
    this.sinkService.getSinks(this.pageFilters).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.sinks,
        };
      },
    );
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      this.pageFilters.offset = this.page.offset - this.page.limit;
    }
    if (dir === 'next') {
      this.pageFilters.offset = this.page.offset + this.page.limit;
    }
    this.getSinks();
  }

  onChangeLimit(limit: number) {
    this.pageFilters.limit = limit;
    this.getSinks();
  }

  onOpenAdd() {
    this.router.navigate([`${this.router.routerState.snapshot.url}/add`]);
  }

  onOpenEdit() {
    // this.dialogService.open(SinksAddComponent, {context: {action: 'Edit'}}).onClose.subscribe(
    //   confirm => {
    //     if (confirm) {
    //       this.getSinks();
    //     }
    //   },
    // );
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
              this.page.rows = this.page.rows.filter((u: User) => u.id !== row.id);
              this.notificationsService.success('Sink Item successfully deleted', '');
            },
          );
        }
      },
    );
  }

  openDetailsModal(row: any) {
    const {name, description, backend, config, ts_created, id} = row;

    this.dialogService.open(SinksDetailsComponent, {
      context: {sink: {id, name, description, backend, config, ts_created}},
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
    const t = new Date().getTime();
    if ((t - this.searchFreq) > defFreq) {
      this.getSinks(input);
      this.searchFreq = t;
    }
  }

  filterByInactive = (sink) => sink.status === 'inactive';

}
