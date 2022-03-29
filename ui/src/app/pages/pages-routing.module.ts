import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { environment } from '../../environments/environment';
import { AgentListComponent } from './agent/agent-list/agent-list.component';
import { GroupListComponent } from './groups/group-list/group-list.component';
import { SinkListComponent } from './sinks/sink-list/sink-list.component';
import { DatasetListComponent } from './datasets/dataset-list/dataset-list.component';
import { PolicyListComponent } from './policies/policy-list/policy-list.component';
import { PagesViewComponent } from './pages-view/pages-view.component';


const children: Routes = [
  {
    path: 'home',
    data: { breadcrumb: 'Home' },
    component: DashboardComponent,
  },
  {
    path: 'fleet',
    data: { breadcrumb: 'Fleet Management' },
    children: [
      {
        path: 'agents',
        children: [
          {
            path: '',
            component: AgentListComponent,
            data: { breadcrumb: 'Agents List' },
          },
        ],
      },
      {
        path: 'groups',
        children: [
          {
            path: '',
            component: GroupListComponent,
            data: { breadcrumb: 'Agent Groups List' },
          },
        ],
      },
      {
        path: 'sinks',
        children: [
          {
            path: '',
            component: SinkListComponent,
            data: { breadcrumb: 'Sink Management' },
          },
        ],
      },
      {
        path: 'datasets',
        children: [
          {
            path: 'list',
            component: DatasetListComponent,
            data: { breadcrumb: 'Dataset Explorer' },
          },
          {
            path: 'policies',
            component: PolicyListComponent,
            data: { breadcrumb: 'Policy Management' },
          },
        ],
      },
    ],
  },
];

const DEV_ROUTES = [
  {
    path: 'dev',
  },
];

const routes: Routes = [{
  path: '',
  component: PagesViewComponent,
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

