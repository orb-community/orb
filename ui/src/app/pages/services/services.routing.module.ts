import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { LoraComponent } from 'app/pages/services/lora/lora.component';
import { LoraDetailsComponent } from 'app/pages/services/lora/details/lora.details.component';
import { OpcuaComponent } from 'app/pages/services/opcua/opcua.component';
import { OpcuaDetailsComponent } from 'app/pages/services/opcua/details/opcua.details.component';
import { GatewaysComponent } from 'app/pages/services/gateways/gateways.component';
import { GatewaysDetailsComponent } from 'app/pages/services/gateways/details/gateways.details.component';

const routes: Routes = [
  {
    path: 'lora',
    component: LoraComponent,
  },
  {
    path: 'lora/details/:id',
    component: LoraDetailsComponent,
  },
  {
    path: 'opcua',
    component: OpcuaComponent,
  },
  {
    path: 'opcua/details/:id',
    component: OpcuaDetailsComponent,
  },
  {
    path: 'gateways',
    component: GatewaysComponent,
  },
  {
    path: 'gateways/details/:id',
    component: GatewaysDetailsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ServicesRoutingModule { }
