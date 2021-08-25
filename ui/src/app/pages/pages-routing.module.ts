import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';

import { PagesComponent } from './pages.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ProfileComponent } from './profile/profile.component';
import { environment } from 'environments/environment';

// ORB
// Agent Group Management
import { AgentGroupsComponent } from 'app/pages/agent-groups/agent.groups.component';
// Dataset Explorer
import { DatasetsComponent } from 'app/pages/datasets/datasets.component';
// Fleet Management
import { FleetsComponent } from 'app/pages/fleets/fleets.component';
// Sink Management
import { SinksComponent } from 'app/pages/sinks/sinks.component';
import { SinksAddComponent } from 'app/pages/sinks/add/sinks.add.component';
import { AgentGroupAddComponent } from 'app/pages/agent-groups/add/agent.group.add.component';
import { ShowcaseComponent } from 'app/pages/showcase/showcase.component';

const children = [
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
    path: 'dev',
    component: ShowcaseComponent,
  },
  {
    path: 'profile',
    component: ProfileComponent,
  },
  {
    path: 'agents',
    component: AgentGroupsComponent,
    data: {breadcrumb: 'Agent Groups'},
  },
  {
    path: 'agents/add',
    component: AgentGroupAddComponent,
    data: {breadcrumb: 'New'},
  },
  {
    path: 'agents/edit',
    component: AgentGroupAddComponent,
    data: {breadcrumb: 'Edit'},
  },
  {
    path: 'datasets',
    component: DatasetsComponent,
  },
  {
    path: 'fleets',
    component: FleetsComponent,
  },
  {
    path: 'sinks',
    component: SinksComponent,
    data: {breadcrumb: 'Sink Management'},
  },
  {
    path: 'sinks/add',
    component: SinksAddComponent,
    data: {breadcrumb: 'New'},
  },
  {
    path: 'sinks/edit',
    component: SinksAddComponent,
    data: {breadcrumb: 'Edit'},
  },
];


const DEV_ROUTES = [
   {
    path: 'dev',
    component: ShowcaseComponent,
  },
];

const routes: Routes = [{
  path: '',
  component: PagesComponent,
   children: [
    ...children,
    ...environment.production ? [] : DEV_ROUTES,
  ],
}];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PagesRoutingModule {
}
