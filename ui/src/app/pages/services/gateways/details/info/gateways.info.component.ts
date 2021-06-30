import { Component, Input } from '@angular/core';

import { Gateway } from 'app/common/interfaces/gateway.interface';

@Component({
  selector: 'ngx-gateways-info',
  templateUrl: './gateways.info.component.html',
  styleUrls: ['./gateways.info.component.scss'],

})
export class GatewaysInfoComponent {
  @Input() gateway: Gateway;

  constructor(
  ) { }
}
