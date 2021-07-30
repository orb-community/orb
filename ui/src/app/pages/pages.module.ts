import { NgModule } from '@angular/core';
import {
  NbMenuModule,
  NbDialogService,
  NbWindowService,
  NbButtonModule,
  NbCardModule,
  NbInputModule,
  NbSelectModule,
  NbCheckboxModule,
  NbListModule,
  NbTabsetModule,
} from '@nebular/theme';
import { ThemeModule } from '../@theme/theme.module';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';

// Mainflux - Dependencies
import { FormsModule } from '@angular/forms';
// Mainflux - Common and Shared
import { SharedModule } from 'app/shared/shared.module';
import { CommonModule } from 'app/common/common.module';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
// Mainflux - User Groups
import { UserGroupsComponent } from 'app/pages/user-groups/user-groups.component';
import { UserGroupsAddComponent } from 'app/pages/user-groups/add/user-groups.add.component';
import { UserGroupsDetailsComponent } from 'app/pages/user-groups/details/user-groups.details.component';
// Mainflux - User
import { UsersComponent } from 'app/pages/users/users.component';
import { UsersAddComponent } from 'app/pages/users/add/users.add.component';
import { UsersDetailsComponent } from 'app/pages/users/details/users.details.component';
// Mainflux - Things
import { ThingsComponent } from 'app/pages/things/things.component';
import { ThingsAddComponent } from 'app/pages/things/add/things.add.component';
import { ThingsDetailsComponent } from 'app/pages/things/details/things.details.component';
// Mainflux - Channels
import { ChannelsComponent } from 'app/pages/channels/channels.component';
import { ChannelsAddComponent } from 'app/pages/channels/add/channels.add.component';
import { ChannelsDetailsComponent } from 'app/pages/channels/details/channels.details.component';
// Mainflux - Twins
import { TwinsComponent } from './twins/twins.component';
import { TwinsAddComponent } from './twins/add/twins.add.component';
import { TwinsDetailsComponent } from './twins/details/twins.details.component';
import { TwinsStatesComponent } from './twins/states/twins.states.component';
import { TwinsDefinitionsComponent } from './twins/definitions/twins.definitions.component';
// ORB SPECIFIC
// ORB - dataset-explorer
import { DatasetsComponent } from 'app/pages/datasets/datasets.component';
import { DatasetsAddComponent } from 'app/pages/datasets/add/datasets.add.component';
import { DatasetsDetailsComponent } from 'app/pages/datasets/details/datasets.details.component';
// ORB - fleet-management
import { FleetsComponent } from 'app/pages/fleets/fleets.component';
import { FleetsAddComponent } from 'app/pages/fleets/add/fleets.add.component';
import { FleetsDetailsComponent } from 'app/pages/fleets/details/fleets.details.component';
import { SinksComponent } from 'app/pages/sinks/sinks.component';
import { SinksAddComponent } from 'app/pages/sinks/add/sinks.add.component';
import { SinksDetailsComponent } from 'app/pages/sinks/details/sinks.details.component';
import { SinksDeleteComponent } from 'app/pages/sinks/delete/sinks.delete.component';
import { AgentsComponent } from 'app/pages/agents/agents.component';
import { AgentsAddComponent } from 'app/pages/agents/add/agents.add.component';
import { AgentsDetailsComponent } from 'app/pages/agents/details/agents.details.component';

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
    // User Groups
    UserGroupsComponent,
    UserGroupsAddComponent,
    UserGroupsDetailsComponent,
    // Users
    UsersComponent,
    UsersAddComponent,
    UsersDetailsComponent,
    // Things
    ThingsComponent,
    ThingsAddComponent,
    ThingsDetailsComponent,
    // Channels
    ChannelsComponent,
    ChannelsAddComponent,
    ChannelsDetailsComponent,
    // Twins
    TwinsComponent,
    TwinsAddComponent,
    TwinsDetailsComponent,
    TwinsStatesComponent,
    TwinsDefinitionsComponent,
    // Orb
    // Agent Group Management
    AgentsComponent,
    AgentsAddComponent,
    AgentsDetailsComponent,
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
