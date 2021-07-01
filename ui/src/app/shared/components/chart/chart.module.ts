import { NgModule } from '@angular/core';
import { ThemeModule } from 'app/@theme/theme.module';
import { ChartsModule } from 'ng2-charts';

import { ChartComponent } from './chart.component';

@NgModule({
  imports: [
    ThemeModule,
    ChartsModule,
  ],
  declarations: [
    ChartComponent,
  ],
  exports: [
    ChartComponent,
  ],
})
export class ChartModule { }
