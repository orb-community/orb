import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { UsersService } from 'app/common/services/users/users.service';
import { FsService } from 'app/common/services/fs/fs.service';
import { UsersAddComponent } from 'app/pages/users/add/users.add.component';

const defFreq: number = 100;

@Component({
  selector: 'ngx-users-component',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'],
})
export class UsersComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', 'Email', 'ID'],
    keys: ['edit', 'details', 'email', 'id'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  searchFreq = 0;

  constructor(
    private router: Router,
    private dialogService: NbDialogService,
    private usersService: UsersService,
    private fsService: FsService,
  ) { }

  ngOnInit() {
    // Fetch all Users
    this.getUsers();
  }

  getUsers(email?: string): void {
    this.usersService.getUsers(this.page.offset, this.page.limit, email).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.users,
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
    this.getUsers();
  }

  onChangeLimit(limit: number) {
    this.pageFilters.limit = limit;
    this.getUsers();
  }

  openAddModal() {
    this.dialogService.open(UsersAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getUsers();
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(UsersAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getUsers();
        }
      },
    );
  }

  onOpenDetails(row: any) {
    if (row.id) {
      this.router.navigate([`${this.router.routerState.snapshot.url}/details/${row.id}`]);
    }
  }

  searchUsersbyEmail(input) {
    const t = new Date().getTime();
    if ((t - this.searchFreq) > defFreq) {
      this.getUsers(input);
      this.searchFreq = t;
    }
  }

  onClickSave() {
    this.fsService.exportToCsv('mfx_users.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
  }
}
