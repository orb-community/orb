import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';

import { PagesComponent } from './pages.component';
import { environment } from 'environments/environment';

// ORB
// Dataset Explorer
import { AgentPolicyAddComponent } from 'app/pages/datasets/policies.agent/add/agent.policy.add.component';
import { AgentPolicyListComponent } from 'app/pages/datasets/policies.agent/list/agent.policy.list.component';
import { DatasetListComponent } from 'app/pages/datasets/list/dataset.list.component';
// Sink Management
import { SinkListComponent } from 'app/pages/sinks/list/sink.list.component';
import { SinkAddComponent } from 'app/pages/sinks/add/sink.add.component';
// Fleet Management
import { AgentListComponent } from 'app/pages/fleet/agents/list/agent.list.component';
import { AgentAddComponent } from 'app/pages/fleet/agents/add/agent.add.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { AgentViewComponent } from './fleet/agents/view/agent.view.component';
import { AgentGroupListComponent } from 'app/pages/fleet/groups/list/agent.group.list.component';
import { AgentGroupAddComponent } from 'app/pages/fleet/groups/add/agent.group.add.component';
// DEV
import { ShowcaseComponent } from 'app/pages/showcase/showcase.component';
import { DashboardComponent } from 'app/pages/dashboard/dashboard.component';
import { DatasetAddComponent } from 'app/pages/datasets/add/dataset.add.component';
import { ProfileComponent } from './profile/profile.component';
import { AgentPolicyViewComponent } from 'app/pages/datasets/policies.agent/view/agent.policy.view.component';

const children = [
  {
    path: 'home',
    data: {breadcrumb: 'Home'},
    component: DashboardComponent,
  },
  {
    path: 'profile',
    data: {breadcrumb: 'User Profile'},
    component: ProfileComponent,
  },
  {
    path: 'dev',
    component: ShowcaseComponent,
    data: {breadcrumb: 'Library Showcase - DEV'},
  },
  {
    path: 'fleet',
    data: {breadcrumb: {'skip': true}},
    children: [
      {
        path: 'agents',
        children: [
          {
            path: '',
            component: AgentListComponent,
            data: {breadcrumb: 'Agents List'},
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
          {
            path: 'view/:id',
            component: AgentViewComponent,
            data: {breadcrumb: 'Agent View'},
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
    data: {breadcrumb: {'skip': true}},
    children: [
      {
        path: 'list',
        component: DatasetListComponent,
        data: {breadcrumb: 'List'},
      },
      {
        path: 'add',
        component: DatasetAddComponent,
        data: {breadcrumb: 'New Dataset'},
      },
      {
        path: 'edit/:id',
        component: DatasetAddComponent,
        data: {breadcrumb: 'Edit Dataset'},
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
          {
            path: 'view/:id',
            component: AgentPolicyViewComponent,
            data: {breadcrumb: 'View Agent Policy'},
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
