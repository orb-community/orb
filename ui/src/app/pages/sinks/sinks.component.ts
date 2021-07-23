import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { User, PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { FsService } from 'app/common/services/fs/fs.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { SinksAddComponent } from 'app/pages/sinks/add/sinks.add.component';
import { SinksService } from 'app/common/services/sinks/sinks.service';

const defFreq: number = 100;

@Component({
  selector: 'ngx-sinks-component',
  templateUrl: './sinks.component.html',
  styleUrls: ['./sinks.component.scss'],
})
export class SinksComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', '', 'Name', 'Description', 'ID'],
    keys: ['edit', 'delete', 'details', 'name', 'description', 'id'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  searchFreq = 0;

  constructor(
    private router: Router,
    private dialogService: NbDialogService,
    private sinkService: SinksService,
    private fsService: FsService,
    private notificationsService: NotificationsService,
  ) { }

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
          rows: resp,
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

  openAddModal() {
    this.dialogService.open(SinksAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getSinks();
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(SinksAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getSinks();
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'Sink Management' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.sinkService.deleteSink(row.id).subscribe(
            resp => {
              this.page.rows = this.page.rows.filter((u: User) => u.id !== row.id);
              this.notificationsService.success('Sink Item successfully deleted', '');
            },
          );
        }
      },
    );
  }

  onOpenDetails(row: any) {
    if (row.id) {
      this.router.navigate([`${this.router.routerState.snapshot.url}/details/${row.id}`]);
    }
  }

  searchSinkItembyName(input) {
    const t = new Date().getTime();
    if ((t - this.searchFreq) > defFreq) {
      this.getSinks(input);
      this.searchFreq = t;
    }
  }

  onClickSave() {
    this.fsService.exportToCsv('sink_items.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
  }
}
