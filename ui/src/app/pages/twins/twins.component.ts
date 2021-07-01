import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { Twin, PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { TwinsService } from 'app/common/services/twins/twins.service';
import { FsService } from 'app/common/services/fs/fs.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { TwinsAddComponent } from './add/twins.add.component';

@Component({
  selector: 'ngx-twins',
  templateUrl: './twins.component.html',
  styleUrls: ['./twins.component.scss'],
})
export class TwinsComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', '', 'Name', 'Created', 'Updated', 'Revision'],
    keys: ['edit', 'delete', 'details', 'name', 'created', 'updated', 'revision'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  searchTime = 0;

  constructor(
    private router: Router,
    private dialogService: NbDialogService,
    private twinsService: TwinsService,
    private fsService: FsService,
  ) { }

  ngOnInit() {
    this.getTwins();
  }

  getTwins(): void {
    this.twinsService.getTwins(this.pageFilters).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.twins,
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
    this.getTwins();
  }

  onChangeLimit(lim: number) {
    this.pageFilters.limit = lim;
    this.getTwins();
  }

  openAddModal() {
    this.dialogService.open(TwinsAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getTwins();
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(TwinsAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getTwins();
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'Twin' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.twinsService.deleteTwin(row.id).subscribe(
            resp => {
              this.page.rows = this.page.rows.filter((t: Twin) => t.id !== row.id);
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

  onEditConfirm(event): void {
    // close edit row
    event.confirm.resolve();

    this.twinsService.editTwin(event.newData).subscribe();
  }

  onSaveFile() {
    this.fsService.exportToCsv('twins.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
  }

  searchTwin(input) {
  }
}
