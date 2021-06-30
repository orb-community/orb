import { NgModule } from '@angular/core';

import { PagesModule } from 'app/pages/pages.module';
import { ServicesRoutingModule } from './services.routing.module';

import { LoraComponent } from 'app/pages/services/lora/lora.component';
import { LoraAddComponent } from 'app/pages/services/lora/add/lora.add.component';
import { LoraDetailsComponent } from 'app/pages/services/lora/details/lora.details.component';
import { OpcuaComponent } from 'app/pages/services/opcua/opcua.component';
import { OpcuaAddComponent } from 'app/pages/services/opcua/add/opcua.add.component';
import { OpcuaDetailsComponent } from 'app/pages/services/opcua/details/opcua.details.component';
import { GatewaysComponent } from 'app/pages/services/gateways/gateways.component';
import { GatewaysAddComponent } from 'app/pages/services/gateways/add/gateways.add.component';
import { GatewaysDetailsComponent } from 'app/pages/services/gateways/details/gateways.details.component';
import { GatewaysInfoComponent } from 'app/pages/services/gateways/details/info/gateways.info.component';
import { GatewaysConfigComponent } from 'app/pages/services/gateways/details/config/gateways.config.component';
import { GatewaysXtermComponent } from 'app/pages/services/gateways/details/xterm/gateways.xterm.component';

import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';

@NgModule({
  imports: [
    ServicesRoutingModule,
    PagesModule,
  ],
  declarations: [
    LoraComponent,
    LoraAddComponent,
    LoraDetailsComponent,
    OpcuaComponent,
    OpcuaAddComponent,
    OpcuaDetailsComponent,
    GatewaysComponent,
    GatewaysAddComponent,
    GatewaysDetailsComponent,
    GatewaysInfoComponent,
    GatewaysConfigComponent,
    GatewaysXtermComponent,
  ],
  entryComponents: [
    ConfirmationComponent,
    LoraAddComponent,
    GatewaysAddComponent,
  ],
})
export class ServicesModule { }
