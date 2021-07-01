import { NgModule } from '@angular/core';

import { MapComponent } from './leaflet/map.leaflet.component';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';

@NgModule({
  imports: [
    LeafletModule.forRoot(),
    LeafletDrawModule.forRoot(),
  ],
  declarations: [
    MapComponent,
  ],
  exports: [
    MapComponent,
  ],
})
export class MapModule { }
