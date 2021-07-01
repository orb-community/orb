import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { PageFilters, TableConfig, TablePage, MsgFilters } from 'app/common/interfaces/mainflux.interface';
import { Gateway } from 'app/common/interfaces/gateway.interface';
import { GatewaysService } from 'app/common/services/gateways/gateways.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { MessagesService } from 'app/common/services/messages/messages.service';
import { FsService } from 'app/common/services/fs/fs.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { GatewaysAddComponent } from './add/gateways.add.component';

const defSearchBarMs: number = 100;

@Component({
  selector: 'ngx-gateways-component',
  templateUrl: './gateways.component.html',
  styleUrls: ['./gateways.component.scss'],
})
export class GatewaysComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', '', 'Name', 'External ID', 'Messages', 'Last Seen'],
    keys: ['edit', 'delete', 'details', 'name', 'externalID', 'messages', 'seen'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  searchTime = 0;

  constructor(
    private router: Router,
    private gatewaysService: GatewaysService,
    private messagesService: MessagesService,
    private notificationsService: NotificationsService,
    private fsService: FsService,
    private dialogService: NbDialogService,
  ) { }

  ngOnInit() {
    this.getGateways();
  }

  getGateways(name?: string): void {
    this.pageFilters.name = name;
    this.gatewaysService.getGateways(this.pageFilters).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.things,
        };

        this.page.rows.forEach((gw: Gateway) => {
          gw.externalID = gw.metadata.external_id;

          const dataChannID: string = gw.metadata ? gw.metadata.data_channel_id : '';
          const msgFilters: MsgFilters = {
            publisher: gw.id,
          };

          this.messagesService.getMessages(dataChannID, gw.key, msgFilters).subscribe(
            (msgResp: any) => {
              if (msgResp.messages) {
                gw.seen = msgResp.messages[0].time;
                gw.messages = msgResp.total;
              }
            },
          );
        });
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
    this.getGateways();
  }

  onChangeLimit(lim: number) {
    this.pageFilters.limit = lim;
    this.getGateways();
  }

  openAddModal() {
    this.dialogService.open(GatewaysAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          setTimeout(
            () => {
              this.getGateways();
            }, 3000,
          );
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(GatewaysAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getGateways();
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'Gateway' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.gatewaysService.deleteGateway(row).subscribe(
            resp => {
              this.page.rows = this.page.rows.filter((g: Gateway) => g.id !== row.id);
              this.notificationsService.success('Gateway successfully deleted', '');
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

  searchGW(input) {
    const t = new Date().getTime();
    if ((t - this.searchTime) > defSearchBarMs) {
      this.getGateways(input);
      this.searchTime = t;
    }
  }

  onClickSave() {
    this.fsService.exportToCsv('gateways.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
  }
}
