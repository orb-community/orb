import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { environment } from 'environments/environment';
import { Gateway } from 'app/common/interfaces/gateway.interface';
import { Config, ConfigContent, ConfigUpdate } from 'app/common/interfaces/bootstrap.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ThingsService } from 'app/common/services/things/things.service';

@Injectable()
export class BootstrapService {
  content: ConfigContent = {
    log_level: 'debug',
    http_port: '9000',
    mqtt_url: 'localhost:1883',
    edgex_url: 'http://localhost:48090/api/v1/',
    nats_url: 'localhost:4222',
    export_config:  {
      file: `${environment.exportConfigFile}`,
      exp: {
        port: '8170',
      },
      mqtt:  {},
      routes: [{}, {}],
    },
  };

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
    private thingsService: ThingsService,
  ) { }

  addConfig(gw: Gateway) {
    // Boostrap
    this.content.export_config.mqtt.channel = gw.metadata.export_channel_id;
    this.content.export_config.mqtt.username = gw.id;
    this.content.export_config.routes[0].mqtt_topic = `channels/${gw.metadata.export_channel_id}/messages`;
    this.content.export_config.routes[1].mqtt_topic = `channels/${gw.metadata.export_channel_id}/messages`;
    this.content.export_config.mqtt.password = gw.key;
    this.content.export_config.exp.nats = this.content.nats_url;


    const config: Config = {
      thing_id: gw.id,
      thing_key: gw.key,
      channels: [gw.metadata.ctrl_channel_id, gw.metadata.data_channel_id],
      external_id: gw.metadata.external_id,
      external_key: gw.metadata.external_key,
      content: JSON.stringify(this.content),
      state: 0,
    };

    return this.http.post(environment.bootstrapConfigsUrl, config, { observe: 'response' })
      .map(
        resp => {
          const cfg_id: string = resp.headers.get('location').replace('/things/configs/', '');
          gw.metadata.cfg_id = cfg_id;
          this.thingsService.editThing(gw).subscribe(
            respEdit => {
              this.notificationsService.success('Gateway successfully bootstrapped', '');
            },
            errEdit => {
              this.notificationsService.error(
                'Failed to add config ID to GW metadata',
                `Error: ${errEdit.status} - ${errEdit.statusText}`);
            },
          );
        },
        err => {
          this.notificationsService.error(
            'Failed to add bootstrap config to gateway',
            `Error: ${err.status} - ${err.statusText}`);
        },
      );
  }

  getConfig(gateway: Gateway) {
    const headers = new HttpHeaders({
      'Authorization': gateway.metadata.external_key,
    });

    return this.http.get(`${environment.bootstrapUrl}/${gateway.metadata.external_id}`, { headers: headers });
  }

  updateConfig(configUpdate: ConfigUpdate, gateway: Gateway) {
    return this.http.put(`${environment.bootstrapConfigsUrl}/${gateway.id}`, configUpdate, { observe: 'response' });
  }
}
