import { NgModule } from '@angular/core';
import { NbCardModule, NbIconModule } from '@nebular/theme';

import { ThemeModule } from '../../@theme/theme.module';
import { DashboardComponent } from './dashboard.component';
import { BreadcrumbModule } from 'xng-breadcrumb';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    NbIconModule,
    BreadcrumbModule,
  ],
  declarations: [
    DashboardComponent,
  ],
})
export class DashboardModule {
}
