import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DashboardComponent } from './dashboard/dashboard.component';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatCardModule } from '@angular/material/card';
import { MatMenuModule } from '@angular/material/menu';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { LayoutModule } from '@angular/cdk/layout';
import { RouterModule } from '@angular/router';
import { PagesRoutingModule } from './pages-routing.module';
import { SharedModule } from '../shared/shared.module';
import { AgentListComponent } from './agent/agent-list/agent-list.component';
import { GroupListComponent } from './groups/group-list/group-list.component';
import { SinkListComponent } from './sinks/sink-list/sink-list.component';
import { DatasetListComponent } from './datasets/dataset-list/dataset-list.component';
import { PolicyListComponent } from './policies/policy-list/policy-list.component';
import { PagesViewComponent } from './pages-view/pages-view.component';


@NgModule({
  declarations: [
    DashboardComponent,
    AgentListComponent,
    GroupListComponent,
    SinkListComponent,
    DatasetListComponent,
    PolicyListComponent,
    PagesViewComponent,
  ],
  imports: [
    CommonModule,
    MatGridListModule,
    MatCardModule,
    MatMenuModule,
    MatIconModule,
    MatButtonModule,
    LayoutModule,
    RouterModule,
    SharedModule,
    PagesRoutingModule,
  ],
})
export class PagesModule {
}
