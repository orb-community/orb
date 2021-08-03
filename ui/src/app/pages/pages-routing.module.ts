import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';

import { PagesComponent } from './pages.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ProfileComponent } from './profile/profile.component';

import { environment } from 'environments/environment';

// ORB
// Agent Group Management
import { AgentsComponent } from 'app/pages/agents/agents.component';
// Dataset Explorer
import { DatasetsComponent } from 'app/pages/datasets/datasets.component';
// Fleet Management
import { FleetsComponent } from 'app/pages/fleets/fleets.component';
// Sink Management
import { SinksComponent } from 'app/pages/sinks/sinks.component';
import { SinksAddComponent } from 'app/pages/sinks/add/sinks.add.component';

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
      path: 'profile',
      component: ProfileComponent,
    },
    {
      path: 'agents',
      component: AgentsComponent,
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
    },
    {
      path: 'sinks/add',
      component: SinksAddComponent,
    },
    {
      path: 'sinks/edit/:id',
      component: SinksAddComponent,
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
      path: 'profile',
      component: ProfileComponent,
    },
    {
      path: 'agents',
      component: AgentsComponent,
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
    },
    {
      path: 'sinks/add',
      component: SinksAddComponent,
    },
    {
      path: 'sinks/edit/:id',
      component: SinksAddComponent,
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
