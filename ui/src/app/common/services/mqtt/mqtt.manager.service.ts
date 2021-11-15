import { Injectable, EventEmitter } from '@angular/core';
import { MqttService, IMqttMessage, MqttConnectionState } from 'ngx-mqtt';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

import { SenMLRec } from 'app/common/interfaces/mainflux.interface';
import { Subscription } from 'rxjs';

@Injectable()
export class MqttManagerService {
  public message: SenMLRec;
  messageChange = new EventEmitter();
  connectChange = new EventEmitter();
  subscriptions: Subscription[] = new Array();

  constructor(
    private mqttService: MqttService,
    private notificationsService: NotificationsService,
  ) {
  }

  init(username: string, password: string, channel: string) {
    const connSub = this.mqttService.onConnect.subscribe(
      resp => {
        this.notificationsService.success('Connected to MQTT broker', '');
        this.subscribe(channel);
        this.connectChange.emit(MqttConnectionState.CONNECTED);
      },
      err => {
        this.notificationsService.error(
          'Failed to connect to MQTT broker.',
          `Error: ${err.status} - ${err.statusText}`);
      },
    );
    this.subscriptions.push(connSub);

    const closeSub = this.mqttService.onClose.subscribe(
      resp => {
        this.notificationsService.success('Disconnected from MQTT broker', '');
        this.connectChange.emit(MqttConnectionState.CLOSED);
        this.subscriptions.forEach(
          sub => { sub.unsubscribe(); },
        );
        this.subscriptions = new Array();
      },
      err => {
        this.notificationsService.error(
          'Failed to disconnect from MQTT broker.',
          `Error: ${err.status} - ${err.statusText}`);
      },
    );
    this.subscriptions.push(closeSub);

    this.connect(username, password);
  }

  connect(username: string, password: string) {
    this.mqttService.connect({
      username: username,
      password: password,
    });
  }

  disconnect() {
    this.mqttService.disconnect();
  }

  createTopic(channel: string) {
    return `channels/${channel}/messages`;
  }

  createPayload(baseName: string, name: string, valueString: string) {
    return `[{"bn":"${baseName}:", "n":"${name}", "vs":"${valueString}"}]`;
  }

  publish(channel: string, bn: string, n: string, vs: string) {
    const topic = `${this.createTopic(channel)}/req`;
    const payload = this.createPayload(bn, n, vs);
    this.mqttService.publish(topic, payload).subscribe();
  }

  publishToService(channel: string, svc: string, bn: string, n: string, vs: string) {
    const topic = `${this.createTopic(channel)}/services/${svc}`;
    const payload = this.createPayload(bn, n, vs);
    this.mqttService.publish(topic, payload).subscribe();
  }

  subscribe(channel: string) {
    const topic = `${this.createTopic(channel)}/res`;
    const topicSub = this.mqttService.observe(topic).subscribe(
      (message: IMqttMessage) => {
        const pl = message.payload.toString();
        this.message = <SenMLRec>JSON.parse(pl)[0];
        this.messageChange.emit(this.message);
      },
    );
    this.subscriptions.push(topicSub);
  }
}
