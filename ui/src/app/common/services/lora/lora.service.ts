import { Injectable } from '@angular/core';

import { LoraDevice } from 'app/common/interfaces/lora.interface';
import { ThingsService } from 'app/common/services/things/things.service';
import { ChannelsService } from 'app/common/services/channels/channels.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

const defLimit: number = 20;
const typeLora = 'lora';
const typeLoraApp = 'loraApp';

@Injectable()
export class LoraService {
  constructor(
    private thingsService: ThingsService,
    private channelsService: ChannelsService,
    private notificationsService: NotificationsService,
  ) { }

  getDevice(id: string) {
    return this.thingsService.getThing(id);
  }

  getDevices(offset?: number, limit?: number, name?: string) {
    const filters = {
      offset: offset || 0,
      limit: limit || defLimit,
      type: typeLora,
      name: name,
    };

    return this.thingsService.getThings(filters);
  }

  getChannels(offset: number, limit: number) {
    const filters = {
      offset: offset || 0,
      limit: limit || defLimit,
      type: typeLora,
    };

    return this.channelsService.getChannels(filters);
  }

  addDevice(row: LoraDevice) {
    const filters = {
      offset: 0,
      limit: 1,
      type: typeLora,
      metadata: `{"app_id": "${row.appID}"}`,
    };

    // Check if a channel exist for row appID
    return this.channelsService.getChannels(filters).map(
      (resp: any) => {
        if (resp.total === 0) {
          const chanReq = {
            name: `${typeLoraApp}-${row.appID}`,
            metadata: {
              type: typeLora,
              lora: {
                app_id: row.appID,
              },
            },
          };

          this.channelsService.addChannel(chanReq).subscribe(
            respChan => {
              const chanID = respChan.headers.get('location').replace('/channels/', '');
              this.addAndConnect(chanID, row);
            },
          );
        } else {
          const chanID = resp.channels[0].id;
          this.addAndConnect(chanID, row);
        }
      },
    );
  }

  addAndConnect(chanID: string, row: LoraDevice) {
    const devReq: LoraDevice = {
      name: row.name,
      metadata: {
        type: typeLora,
        channel_id: chanID,
        lora: {
          dev_eui: row.devEUI,
          app_id: row.appID,
        },
      },
    };

    this.thingsService.addThing(devReq).subscribe(
      respThing => {
        const thingID = respThing.headers.get('location').replace('/things/', '');

        this.channelsService.connectThing(chanID, thingID).subscribe(
          respCon => {
            this.notificationsService.success('LoRa Device successfully created', '');
          },
          err => {
            this.thingsService.deleteThing(thingID).subscribe();
            this.channelsService.deleteChannel(chanID).subscribe();
          },
        );
      },
      err => {
        this.channelsService.deleteChannel(chanID).subscribe();
      },
    );
  }

  editDevice(row: LoraDevice) {
    const devReq: LoraDevice = {
      id: row.id,
      name: row.name,
      metadata: row.metadata,
    };

    devReq.metadata.lora = {
        dev_eui: row.devEUI,
        app_id: row.appID,
    };

    return this.thingsService.editThing(row).map(
      resp => {
        this.notificationsService.success('LoRa Device successfully edited', '');
      },
    );
  }

  deleteDevice(loraDev: LoraDevice) {
    const channelID = loraDev.metadata.channel_id;
    return this.channelsService.deleteChannel(channelID).map(
      () => {
        this.thingsService.deleteThing(loraDev.id).subscribe(
          resp => {
            this.notificationsService.success('LoRa Device successfully deleted', '');
          },
        );
      },
    );
  }
}
