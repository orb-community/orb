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
import { AgentsComponent } from 'app/pages/agents/agents.component';
import { AgentsDetailsComponent } from 'app/pages/agents/details/agents.details.component';
// Dataset Explorer
import {DatasetsComponent} from 'app/pages/datasets/datasets.component';
import {DatasetsDetailsComponent} from 'app/pages/datasets/details/datasets.details.component';
// Fleet Management
import {FleetsComponent} from 'app/pages/fleets/fleets.component';
import {FleetsDetailsComponent} from 'app/pages/fleets/details/fleets.details.component';
// Sink Management
import { SinksComponent } from 'app/pages/sinks/sinks.component';
import { SinksDetailsComponent } from 'app/pages/sinks/details/sinks.details.component';


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
            component: AgentsComponent,
        },
        {
            path: 'agent-groups-management/:id',
            component: AgentsDetailsComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinksDetailsComponent,
        },
        {
            path: 'dataset-explorer',
            component: DatasetsComponent,
        },
        {
            path: 'dataset-explorer/:id',
            component: DatasetsDetailsComponent,
        },
        {
            path: 'fleet-management',
            component: FleetsComponent,
        },
        {
            path: 'fleet-management/:id',
            component: FleetsDetailsComponent,
        },
        {
            path: 'sink-management',
            component: SinksComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinksDetailsComponent,
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
            component: AgentsComponent,
        },
        {
            path: 'agent-groups-management/:id',
            component: AgentsDetailsComponent,
        },
        {
            path: 'sink-management/:id',
            component: SinksDetailsComponent,
        },
        {
            path: 'dataset-explorer',
            component: DatasetsComponent,
        },
        {
            path: 'dataset-explorer/:id',
            component: DatasetsDetailsComponent,
        },
        {
            path: 'fleet-management',
            component: FleetsComponent,
        },
        {
            path: 'fleet-management/:id',
            component: FleetsDetailsComponent,
        },
        {
            path: 'sink-management',
            component: SinksComponent,
        },
        {
            path: 'sink-management/details/:id',
            component: SinksDetailsComponent,
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
