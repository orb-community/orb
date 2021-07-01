import { Component, Input, OnChanges } from '@angular/core';
import { Gateway } from 'app/common/interfaces/gateway.interface';
import { Config, ConfigContent, Route, ConfigUpdate } from 'app/common/interfaces/bootstrap.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { environment } from 'environments/environment';
import { BootstrapService } from 'app/common/services/bootstrap/bootstrap.service';

@Component({
  selector: 'ngx-gateways-config',
  templateUrl: './gateways.config.component.html',
  styleUrls: ['./gateways.config.component.scss'],
})
export class GatewaysConfigComponent implements OnChanges {
  @Input() gateway: Gateway;

  content: ConfigContent = {
    log_level: '',
    http_port: '',
    mqtt_url: '',
    edgex_url: '',
    nats_url: '',
    export_config: {
      file: `${environment.exportConfigFile}`,
      mqtt : {},
      exp: {},
      routes: Array<Route>(2),
    },
  };

  constructor(
    private bootstrapService: BootstrapService,
    private notificationsService: NotificationsService,
  ) {}

  ngOnChanges() {
    if (this.gateway.metadata.external_key === '') {
      return;
    }

    this.bootstrapService.getConfig(this.gateway).subscribe(
      resp => {
        const cfg = <Config>resp;
        this.content = JSON.parse(cfg.content);
      },
      err => {
        this.notificationsService.error(
          'Failed to get bootstrap configuration',
          `Error: ${err.status} - ${err.statusText}`);
      },
    );
  }

  submit() {
    this.content.export_config.mqtt.host = `tcp://${this.content.mqtt_url}`;
    this.content.export_config.exp.nats = `nats://${this.content.nats_url}`;
    const configUpdate: ConfigUpdate = {
      content: JSON.stringify(this.content),
      name: this.gateway.name,
    };

    this.bootstrapService.updateConfig(configUpdate, this.gateway).subscribe(
      resp => {
        this.notificationsService.success('Bootstrap configuration updated', '');
      },
      err => {
        this.notificationsService.error(
          'Failed to update bootstrap configuration',
          `Error: ${err.status} - ${err.statusText}`);
      },
    );
  }
}
