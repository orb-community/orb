import { NgModule } from '@angular/core';
import {
  NbButtonModule,
  NbCardModule,
  NbCheckboxModule,
  NbDialogService,
  NbFormFieldModule,
  NbInputModule,
  NbListModule,
  NbMenuModule,
  NbSelectModule,
  NbStepperModule,
  NbTabsetModule,
  NbWindowService,
} from '@nebular/theme';
import { ThemeModule } from '../@theme/theme.module';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';

// Mainflux - Dependencies
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Mainflux - Common and Shared
import { SharedModule } from 'app/shared/shared.module';
import { CommonModule } from 'app/common/common.module';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';

// ORB
import { DatasetsComponent } from 'app/pages/datasets/datasets.component';
import { DatasetsAddComponent } from 'app/pages/datasets/add/datasets.add.component';
import { DatasetsDetailsComponent } from 'app/pages/datasets/details/datasets.details.component';
import { FleetsComponent } from 'app/pages/fleets/fleets.component';
import { FleetsAddComponent } from 'app/pages/fleets/add/fleets.add.component';
import { FleetsDetailsComponent } from 'app/pages/fleets/details/fleets.details.component';
import { SinksComponent } from 'app/pages/sinks/sinks.component';
import { SinksAddComponent } from 'app/pages/sinks/add/sinks.add.component';
import { SinksDetailsComponent } from 'app/pages/sinks/details/sinks.details.component';
import { SinksDeleteComponent } from 'app/pages/sinks/delete/sinks.delete.component';
import { MatInputModule } from '@angular/material/input';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { BreadcrumbModule } from 'xng-breadcrumb';
import { AgentGroupsComponent } from 'app/pages/agent-groups/agent.groups.component';
import { AgentGroupAddComponent } from 'app/pages/agent-groups/add/agent.group.add.component';
import { AgentGroupDetailsComponent } from 'app/pages/agent-groups/details/agent.group.details.component';
import { AgentGroupDeleteComponent } from 'app/pages/agent-groups/delete/agent.group.delete.component';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { ShowcaseComponent } from 'app/pages/showcase/showcase.component';
import { DebounceModule } from 'ngx-debounce';

@NgModule({
  imports: [
    PagesRoutingModule,
    ThemeModule,
    NbMenuModule,
    DashboardModule,
    SharedModule,
    CommonModule,
    FormsModule,
    NbButtonModule,
    NbCardModule,
    NbInputModule,
    NbSelectModule,
    NbCheckboxModule,
    NbListModule,
    NbTabsetModule,
    ReactiveFormsModule,
    MatInputModule,
    MatChipsModule,
    MatIconModule,
    BreadcrumbModule,
    NbStepperModule,
    NbFormFieldModule,
    NgxDatatableModule,
    DebounceModule,
  ],
  exports: [
    SharedModule,
    CommonModule,
    FormsModule,
    NbButtonModule,
    NbCardModule,
    NbInputModule,
    NbSelectModule,
    NbCheckboxModule,
    NbListModule,
  ],
  declarations: [
    PagesComponent,
    // Orb
    // Agent Group Management
    AgentGroupsComponent,
    AgentGroupAddComponent,
    AgentGroupDetailsComponent,
    AgentGroupDeleteComponent,
    // Dataset Explorer
    DatasetsComponent,
    DatasetsAddComponent,
    DatasetsDetailsComponent,
    // Fleet Management
    FleetsComponent,
    FleetsAddComponent,
    FleetsDetailsComponent,
    // Sink Management
    SinksComponent,
    SinksAddComponent,
    SinksDetailsComponent,
    SinksDeleteComponent,
    ShowcaseComponent,
  ],
  providers: [
    NbDialogService,
    NbWindowService,
  ],
  entryComponents: [
    ConfirmationComponent,
  ],
})
export class PagesModule {
}
