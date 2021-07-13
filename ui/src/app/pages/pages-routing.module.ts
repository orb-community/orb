import {RouterModule, Routes} from '@angular/router';
import {NgModule} from '@angular/core';

import {PagesComponent} from './pages.component';
import {DashboardComponent} from './dashboard/dashboard.component';
import {ProfileComponent} from './profile/profile.component';

import {environment} from 'environments/environment';

// Mainflux - User Groups
import {UserGroupsComponent} from './user-groups/user-groups.component';
import {UserGroupsDetailsComponent} from './user-groups/details/user-groups.details.component';
// Mainflux - User
import {UsersComponent} from './users/users.component';
import {UsersDetailsComponent} from './users/details/users.details.component';
// Mainflux - Things
import {ThingsComponent} from './things/things.component';
import {ThingsDetailsComponent} from './things/details/things.details.component';
// Mainflux - Channels
import {ChannelsComponent} from './channels/channels.component';
import {ChannelsDetailsComponent} from './channels/details/channels.details.component';
// Mainflux - Twins
import {TwinsComponent} from './twins/twins.component';
import {TwinsDetailsComponent} from './twins/details/twins.details.component';
import {TwinsStatesComponent} from './twins/states/twins.states.component';
import {TwinsDefinitionsComponent} from './twins/definitions/twins.definitions.component';
// ORB
// Agent Group Management
import { AgentGroupsManagementComponent } from './agent-groups-management/agent-groups-management.component';
import { AgentGroupsManagementDetailsComponent } from './agent-groups-management/details/agent-groups-management.details.component';
// Dataset Explorer
import {DatasetExplorerComponent} from 'app/pages/dataset-explorer/dataset-explorer.component';
import {DatasetExplorerDetailsComponent} from 'app/pages/dataset-explorer/details/dataset-explorer.details.component';
// Fleet Management
import {FleetManagementComponent} from 'app/pages/fleet-management/fleet-management.component';
import {FleetManagementDetailsComponent} from 'app/pages/fleet-management/details/fleet-management.details.component';
// Sink Management
import { SinkManagementComponent } from './sink-management/sink-management.component';
import { SinkManagementDetailsComponent } from './sink-management/details/sink-management.details.component';


const children = environment.production ?
    [
        {
            path: 'home',
            component: DashboardComponent,
        },
        {
            path: '',
            redirectTo: 'home',
            pathMatch: 'full',
        },
        {
            path: 'things',
            component: ThingsComponent,
        },
        {
            path: 'things/details/:id',
            component: ThingsDetailsComponent,
        },
        {
            path: 'channels',
            component: ChannelsComponent,
        },
        {
            path: 'channels/details/:id',
            component: ChannelsDetailsComponent,
        },
        {
            path: 'twins/states/:id',
            component: TwinsStatesComponent,
        },
        {
            path: 'twins/definitions/:id',
            component: TwinsDefinitionsComponent,
        },
        {
            path: 'profile',
            component: ProfileComponent,
        },
        {
            path: 'agent-groups-management',
            component: AgentGroupsManagementComponent,
        },
        {
            path: 'agent-groups-management/:id',
            component: AgentGroupsManagementDetailsComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinkManagementDetailsComponent,
        },
        {
            path: 'dataset-explorer',
            component: DatasetExplorerComponent,
        },
        {
            path: 'dataset-explorer/:id',
            component: DatasetExplorerDetailsComponent,
        },
        {
            path: 'fleet-management',
            component: FleetManagementComponent,
        },
        {
            path: 'fleet-management/:id',
            component: FleetManagementDetailsComponent,
        },
        {
            path: 'sink-management',
            component: SinkManagementComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinkManagementDetailsComponent,
        },

    ] : [
        {
            path: 'home',
            component: DashboardComponent,
        },
        {
            path: '',
            redirectTo: 'home',
            pathMatch: 'full',
        },
        {
            path: 'users/groups',
            component: UserGroupsComponent,
        },
        {
            path: 'users/groups/details/:id',
            component: UserGroupsDetailsComponent,
        },
        {
            path: 'users',
            component: UsersComponent,
        },
        {
            path: 'users/details/:id',
            component: UsersDetailsComponent,
        },
        {
            path: 'things',
            component: ThingsComponent,
        },
        {
            path: 'things/details/:id',
            component: ThingsDetailsComponent,
        },
        {
            path: 'channels',
            component: ChannelsComponent,
        },
        {
            path: 'channels/details/:id',
            component: ChannelsDetailsComponent,
        },
        {
            path: 'twins',
            component: TwinsComponent,
        },
        {
            path: 'twins/details/:id',
            component: TwinsDetailsComponent,
        },
        {
            path: 'twins/states/:id',
            component: TwinsStatesComponent,
        },
        {
            path: 'twins/definitions/:id',
            component: TwinsDefinitionsComponent,
        },
        {
            path: 'profile',
            component: ProfileComponent,
        },
        {
            path: 'services',
            loadChildren: () => import('./services/services.module')
                .then(m => m.ServicesModule),
        },
        {
            path: 'agent-groups-management',
            component: AgentGroupsManagementComponent,
        },
        {
            path: 'agent-groups-management/:id',
            component: AgentGroupsManagementDetailsComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinkManagementDetailsComponent,
        },
        {
            path: 'dataset-explorer',
            component: DatasetExplorerComponent,
        },
        {
            path: 'dataset-explorer/:id',
            component: DatasetExplorerDetailsComponent,
        },
        {
            path: 'fleet-management',
            component: FleetManagementComponent,
        },
        {
            path: 'fleet-management/:id',
            component: FleetManagementDetailsComponent,
        },
        {
            path: 'sink-management',
            component: SinkManagementComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinkManagementDetailsComponent,
        },
    ];

const routes: Routes = [{
    path: '',
    component: PagesComponent,
    children: children,
}];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class PagesRoutingModule {
}
