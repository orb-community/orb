import { A11yModule } from '@angular/cdk/a11y';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { NgModule } from '@angular/core';

// Mainflux - Dependencies
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import {
  NbAccordionModule,
  NbAlertModule,
  NbAutocompleteModule,
  NbButtonModule,
  NbCardModule,
  NbCheckboxModule,
  NbDialogService,
  NbFormFieldModule,
  NbInputModule,
  NbListModule,
  NbMenuModule,
  NbSelectModule,
  NbStepperModule,
  NbTabsetModule,
  NbToggleModule,
  NbTooltipModule,
  NbWindowService,
} from '@nebular/theme';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { CommonModule } from 'app/common/common.module';
import { DatasetAddComponent } from 'app/pages/datasets/add/dataset.add.component';
import { DatasetFromComponent } from 'app/pages/datasets/dataset-from/dataset-from.component';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';
import { DatasetDetailsComponent } from 'app/pages/datasets/details/dataset.details.component';
import { DatasetListComponent } from 'app/pages/datasets/list/dataset.list.component';
import { AgentPolicyAddComponent } from 'app/pages/datasets/policies.agent/add/agent.policy.add.component';
import { HandlerPolicyAddComponent } from 'app/pages/datasets/policies.agent/add/handler.policy.add.component';
import { AgentPolicyDeleteComponent } from 'app/pages/datasets/policies.agent/delete/agent.policy.delete.component';
import { AgentPolicyDetailsComponent } from 'app/pages/datasets/policies.agent/details/agent.policy.details.component';
import { AgentPolicyListComponent } from 'app/pages/datasets/policies.agent/list/agent.policy.list.component';
import { AgentPolicyViewComponent } from 'app/pages/datasets/policies.agent/view/agent.policy.view.component';
import { AgentAddComponent } from 'app/pages/fleet/agents/add/agent.add.component';
import { AgentDeleteComponent } from 'app/pages/fleet/agents/delete/agent.delete.component';
import { AgentDetailsComponent } from 'app/pages/fleet/agents/details/agent.details.component';
import { AgentKeyComponent } from 'app/pages/fleet/agents/key/agent.key.component';
import { AgentListComponent } from 'app/pages/fleet/agents/list/agent.list.component';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { AgentGroupAddComponent } from 'app/pages/fleet/groups/add/agent.group.add.component';
import { AgentGroupDeleteComponent } from 'app/pages/fleet/groups/delete/agent.group.delete.component';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { AgentGroupListComponent } from 'app/pages/fleet/groups/list/agent.group.list.component';
import { ShowcaseComponent } from 'app/pages/showcase/showcase.component';
import { SinkAddComponent } from 'app/pages/sinks/add/sink.add.component';
import { SinkDeleteComponent } from 'app/pages/sinks/delete/sink.delete.component';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';

// ORB
import { SinkListComponent } from 'app/pages/sinks/list/sink.list.component';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { SortPipe } from 'app/shared/pipes/sort.pipe';
// Mainflux - Common and Shared
import { SharedModule } from 'app/shared/shared.module';
import { DebounceModule } from 'ngx-debounce';
import { MonacoEditorModule } from 'ngx-monaco-editor';
import { BreadcrumbModule } from 'xng-breadcrumb';
import { ThemeModule } from '../@theme/theme.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { AgentViewComponent } from './fleet/agents/view/agent.view.component';
import { PagesRoutingModule } from './pages-routing.module';
import { PagesComponent } from './pages.component';

@NgModule({
  imports: [
    PagesRoutingModule,
    ThemeModule,
    NbMenuModule,
    DashboardModule,
    SharedModule,
    ClipboardModule,
    CommonModule,
    FormsModule,
    NbButtonModule,
    NbCardModule,
    NbInputModule,
    NbSelectModule,
    NbCheckboxModule,
    NbListModule,
    NbTabsetModule,
    ReactiveFormsModule,
    MatInputModule,
    MatChipsModule,
    MatIconModule,
    BreadcrumbModule,
    NbStepperModule,
    NbFormFieldModule,
    NgxDatatableModule,
    DebounceModule,
    NbTooltipModule,
    NbAlertModule,
    NbAccordionModule,
    MonacoEditorModule,
    NbToggleModule,
    A11yModule,
    NbAutocompleteModule,
    SharedModule,
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
    // Orb
    // Fleet Management
    // Fleet - Agents
    AgentListComponent,
    AgentAddComponent,
    AgentDeleteComponent,
    AgentDetailsComponent,
    AgentKeyComponent,
    AgentMatchComponent,
    AgentViewComponent,
    // Fleet - Agent Groups
    AgentGroupListComponent,
    AgentGroupAddComponent,
    AgentGroupDeleteComponent,
    AgentGroupDetailsComponent,
    // Dataset Explorer
    DatasetAddComponent,
    DatasetListComponent,
    DatasetDeleteComponent,
    DatasetDetailsComponent,
    DatasetFromComponent,
    // Dataset Explorer - Agent Policies
    AgentPolicyAddComponent,
    AgentPolicyDeleteComponent,
    AgentPolicyDetailsComponent,
    AgentPolicyListComponent,
    AgentPolicyViewComponent,
    HandlerPolicyAddComponent,
    // Sink Management
    SinkListComponent,
    SinkAddComponent,
    SinkDetailsComponent,
    SinkDeleteComponent,
    // DEV SHOWCASE
    ShowcaseComponent,
  ],
  providers: [NbDialogService, NbWindowService, SortPipe],
  entryComponents: [ConfirmationComponent],
})
export class PagesModule {}
