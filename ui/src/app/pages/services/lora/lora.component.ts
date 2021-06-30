import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { LoraService } from 'app/common/services/lora/lora.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { MessagesService } from 'app/common/services/messages/messages.service';
import { FsService } from 'app/common/services/fs/fs.service';
import { LoraDevice } from 'app/common/interfaces/lora.interface';
import { MsgFilters, PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { LoraAddComponent } from './add/lora.add.component';

const defSearchBarMs: number = 100;

@Component({
  selector: 'ngx-lora-component',
  templateUrl: './lora.component.html',
  styleUrls: ['./lora.component.scss'],
})
export class LoraComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', '', 'Name', 'Application ID', 'Device EUI', 'Messages', 'Last Seen'],
    keys: ['edit', 'delete', 'details', 'name', 'appID', 'devEUI', 'messages', 'seen'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  searchTime = 0;

  constructor(
    private router: Router,
    private loraService: LoraService,
    private messagesService: MessagesService,
    private notificationsService: NotificationsService,
    private fsService: FsService,
    private dialogService: NbDialogService,
  ) { }

  ngOnInit() {
    this.getLoraDevices();
  }

  getLoraDevices(name?: string): void {
    this.loraService.getDevices(this.pageFilters.offset, this.pageFilters.limit, name).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.things,
        };

        this.page.rows.forEach((lora: LoraDevice) => {
          if (lora.metadata.lora !== undefined) {
            lora.devEUI = lora.metadata.lora.dev_eui;
            lora.appID = lora.metadata.lora.app_id;

            const chanID: string = lora.metadata.channel_id;
            const msgFilters: MsgFilters = {
              publisher: lora.id,
            };

            this.messagesService.getMessages(chanID, lora.key, msgFilters).subscribe(
              (msgResp: any) => {
                if (msgResp.messages) {
                  lora.seen = msgResp.messages[0].time;
                  lora.messages = msgResp.total;
                }
              },
            );
          }
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
    this.getLoraDevices();
  }

  onChangeLimit(lim: number) {
    this.pageFilters.limit = lim;
    this.getLoraDevices();
  }

  openAddModal() {
    this.dialogService.open(LoraAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getLoraDevices();
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(LoraAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getLoraDevices();
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'LoRa Device' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.loraService.deleteDevice(row).subscribe(
            resp => {
              this.page.rows = this.page.rows.filter((t: LoraDevice) => t.id !== row.id);
              this.notificationsService.success('LoRa device successfully deleted', '');
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

  searchLora(input) {
    const t = new Date().getTime();
    if ((t - this.searchTime) > defSearchBarMs) {
      this.getLoraDevices(input);
      this.searchTime = t;
    }
  }

  onClickSave() {
    this.fsService.exportToCsv('lora_devices.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
  }
}
