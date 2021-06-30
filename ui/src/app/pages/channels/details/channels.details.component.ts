import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { ChannelsService } from 'app/common/services/channels/channels.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { Channel, Thing, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';


@Component({
  selector: 'ngx-channels-details-component',
  templateUrl: './channels.details.component.html',
  styleUrls: ['./channels.details.component.scss'],
})
export class ChannelsDetailsComponent implements OnInit {
  channel: Channel = {};
  thingKey = '';

  tableConfig: TableConfig = {
    colNames: ['Name', 'Thing ID'],
    keys: ['name', 'id', 'checkbox'],
  };

  connThingsPage: TablePage = {};
  disconnThingsPage: TablePage = {};

  thingsToConnect: string[] = [];
  thingsToDisconnect: string[] = [];

  editorMetadata = '';

  constructor(
    private route: ActivatedRoute,
    private channelsService: ChannelsService,
    private notificationsService: NotificationsService,
  ) {}

  ngOnInit() {
    const chanID = this.route.snapshot.paramMap.get('id');

    this.channelsService.getChannel(chanID).subscribe(
      (ch: Channel) => {
        this.channel = ch;
        this.updateConnections();
      },
    );
  }

  onEdit() {
    if (this.editorMetadata !== '') {
      try {
        this.channel.metadata = JSON.parse(this.editorMetadata);
      } catch (e) {
        this.notificationsService.error('Wrong metadata format', '');
        return;
      }
    }

    this.channelsService.editChannel(this.channel).subscribe(
      resp => {
        this.notificationsService.success('Channel metadata successfully edited', '');
      },
    );
  }

  onConnect() {
    if (this.thingsToConnect.length > 0) {
      this.channelsService.connectThings([this.channel.id], this.thingsToConnect).subscribe(
        resp => {
          this.updateConnections();
          this.notificationsService.success('Thing(s) successfully connected', '');
        },
      );
    } else {
      this.notificationsService.warn('Thing(s) must be provided', '');
    }
  }

  onDisconnect() {
    this.thingsToDisconnect.forEach(thingID => {
      this.channelsService.disconnectThing(this.channel.id, thingID).subscribe(
        resp => {
          this.updateConnections();
          this.notificationsService.success('Thing successfully disconnected', '');
        },
      );
    });
  }

  updateConnections() {
    this.thingsToConnect = [];
    this.thingsToDisconnect = [];
    this.findConnectedThings();
    this.findDisconnectedThings();
  }

  findConnectedThings(offset?: number, limit?: number) {
    this.channelsService.connectedThings(this.channel.id, offset, limit).subscribe(
      (resp: any) => {
        this.connThingsPage = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.things,
        };
        if (this.connThingsPage.rows.length > 0) {
          const thing: Thing = this.connThingsPage.rows[0];
          this.thingKey = thing.key;
        }
      },
    );
  }

  findDisconnectedThings(offset?: number, limit?: number) {
    this.channelsService.disconnectedThings(this.channel.id, offset, limit).subscribe(
      (resp: any) => {
        this.disconnThingsPage = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.things,
        };
      },
    );
  }

  onChangeLimit(limit: number) {
    this.findConnectedThings(0, limit);
  }

  onChangeLimitDisconn(limit: number) {
    this.findDisconnectedThings(0, limit);
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      const offset = this.connThingsPage.offset - this.connThingsPage.limit;
      this.findConnectedThings(offset, this.connThingsPage.limit);
    }
    if (dir === 'next') {
      const offset = this.connThingsPage.offset + this.connThingsPage.limit;
      this.findConnectedThings(offset, this.connThingsPage.limit);
    }
  }

  onChangePageDisconn(dir: any) {
    if (dir === 'prev') {
      const offset = this.disconnThingsPage.offset - this.disconnThingsPage.limit;
      this.findDisconnectedThings(offset, this.disconnThingsPage.limit);
    }
    if (dir === 'next') {
      const offset = this.disconnThingsPage.offset + this.disconnThingsPage.limit;
      this.findDisconnectedThings(offset, this.disconnThingsPage.limit);
    }
  }

  onCheckboxConns(row: any) {
    const index = this.thingsToConnect.indexOf(row.id);
    (index > -1) ? this.thingsToConnect.splice(index, 1) : this.thingsToConnect.push(row.id);
  }

  onCheckboxDisconns(row: any) {
    const index = this.thingsToDisconnect.indexOf(row.id);
    (index > -1) ? this.thingsToDisconnect.splice(index, 1) : this.thingsToDisconnect.push(row.id);
  }
}
