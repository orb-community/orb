import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { User, PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { UserGroupsService } from 'app/common/services/users/groups.service';
import { FsService } from 'app/common/services/fs/fs.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { UserGroupsAddComponent } from 'app/pages/user-groups/add/user-groups.add.component';

const defFreq: number = 100;

@Component({
  selector: 'ngx-user-groups-component',
  templateUrl: './user-groups.component.html',
  styleUrls: ['./user-groups.component.scss'],
})
export class UserGroupsComponent implements OnInit {
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
    private userGroupsService: UserGroupsService,
    private fsService: FsService,
    private notificationsService: NotificationsService,
  ) { }

  ngOnInit() {
    // Fetch all User Groups
    this.getGroups();
  }

  getGroups(name?: string): void {
    this.userGroupsService.getGroups(this.page.offset, this.page.limit, name).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.groups,
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
    this.getGroups();
  }

  onChangeLimit(limit: number) {
    this.pageFilters.limit = limit;
    this.getGroups();
  }

  openAddModal() {
    this.dialogService.open(UserGroupsAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getGroups();
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(UserGroupsAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getGroups();
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'User Group' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.userGroupsService.deleteGroup(row.id).subscribe(
            resp => {
              this.page.rows = this.page.rows.filter((u: User) => u.id !== row.id);
              this.notificationsService.success('User Group successfully deleted', '');
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

  searcUserGroupsbyName(input) {
    const t = new Date().getTime();
    if ((t - this.searchFreq) > defFreq) {
      this.getGroups(input);
      this.searchFreq = t;
    }
  }

  onClickSave() {
    this.fsService.exportToCsv('mfx_user_groups.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
  }
}
