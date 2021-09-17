import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';

import { PagesComponent } from './pages.component';
import { environment } from 'environments/environment';

// ORB
// Dataset Explorer
import { AgentPolicyListComponent } from 'app/pages/datasets/policies.agent/list/agent.policy.list.component';
// Sink Management
import { SinkListComponent } from 'app/pages/sinks/list/sink.list.component';
import { SinkAddComponent } from 'app/pages/sinks/add/sink.add.component';
// Fleet Management
import { AgentListComponent } from 'app/pages/fleet/agents/list/agent.list.component';
import { AgentAddComponent } from 'app/pages/fleet/agents/add/agent.add.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { AgentGroupListComponent } from 'app/pages/fleet/groups/list/agent.group.list.component';
import { AgentGroupAddComponent } from 'app/pages/fleet/groups/add/agent.group.add.component';
// DEV
import { ShowcaseComponent } from 'app/pages/showcase/showcase.component';
import { DatasetListComponent } from 'app/pages/datasets/list/dataset.list.component';
import { AgentPolicyAddComponent } from 'app/pages/datasets/policies.agent/add/agent.policy.add.component';

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
    data: {breadcrumb: 'Datasets Explorer'},
    children: [
      {
        path: 'list',
        component: DatasetListComponent,
        data: {breadcrumb: 'List'},
      },
      {
        path: 'policies',
        children: [
          {
            path: '',
            component: AgentPolicyListComponent,
            data: {breadcrumb: 'Policy Management'},
          },
          {
            path: 'add',
            component: AgentPolicyAddComponent,
            data: {breadcrumb: 'New Agent Policy'},
          },
          {
            path: 'edit/:id',
            component: AgentPolicyAddComponent,
            data: {breadcrumb: 'Edit Agent Policy'},
          },
        ],
      },
    ],
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
