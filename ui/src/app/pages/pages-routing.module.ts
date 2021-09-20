import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';

import { PagesComponent } from './pages.component';
import { environment } from 'environments/environment';

// ORB
// Agent Group Management
// Dataset Explorer
import { DatasetsComponent } from 'app/pages/datasets/datasets.component';
// Sink Management
import { SinkListComponent } from 'app/pages/sinks/list/sink.list.component';
import { SinkAddComponent } from 'app/pages/sinks/add/sink.add.component';
import { ShowcaseComponent } from 'app/pages/showcase/showcase.component';
import { AgentListComponent } from 'app/pages/fleet/agents/list/agent.list.component';
import { AgentAddComponent } from 'app/pages/fleet/agents/add/agent.add.component';
import { AgentGroupListComponent } from 'app/pages/fleet/groups/list/agent.group.list.component';
import { AgentGroupAddComponent } from 'app/pages/fleet/groups/add/agent.group.add.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';

const children = [
  {
    path: 'home',
    redirectTo: 'sinks',
    // component: DashboardComponent,
  },
  {
    path: 'dev',
    component: ShowcaseComponent,
    data: {breadcrumb: 'Library Showcase - DEV'},
  },
  {
    path: 'fleet',
    data: {breadcrumb: 'Fleet Management'},
    children: [
      {
        path: 'agents',
        children: [
          {
            path: '',
            component: AgentListComponent,
            data: {breadcrumb: 'Agent List'},
          },
          {
            path: 'add',
            component: AgentAddComponent,
            data: {breadcrumb: 'New Agent'},
          },
          {
            path: 'edit/:id',
            component: AgentAddComponent,
            data: {breadcrumb: 'Edit Agent'},
          },
          {
            path: 'details/:id',
            component: AgentDetailsComponent,
            data: {breadcrumb: 'Agent Detail'},
          },
        ],
      },
      {
        path: 'groups',
        children: [
          {
            path: '',
            component: AgentGroupListComponent,
            data: {breadcrumb: 'Agent Groups List'},
          },
          {
            path: 'add',
            component: AgentGroupAddComponent,
            data: {breadcrumb: 'New Agent Group'},
          },
          {
            path: 'edit/:id',
            component: AgentGroupAddComponent,
            data: {breadcrumb: 'Edit Agent Group'},
          },
        ],
      },
    ],
  },
  {
    path: 'sinks',
    children: [
      {
        path: '',
        component: SinkListComponent,
        data: {breadcrumb: 'Sink Management'},
      },
      {
        path: 'add',
        component: SinkAddComponent,
        data: {breadcrumb: 'New Sink'},
      },
      {
        path: 'edit/:id',
        component: SinkAddComponent,
        data: {breadcrumb: 'Edit Sink'},
      },
    ],
  },
  {
    path: 'datasets',
    component: DatasetsComponent,
    data: {breadcrumb: 'Datasets Management'},
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
