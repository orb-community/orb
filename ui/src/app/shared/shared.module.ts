import { NgModule } from '@angular/core';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';

import { ThemeModule } from 'app/@theme/theme.module';
import {
  NbAccordionModule,
  NbButtonModule,
  NbCardModule,
  NbCheckboxModule,
  NbDatepickerModule,
  NbDialogModule,
  NbDialogService,
  NbIconModule,
  NbInputModule,
  NbListModule,
  NbSelectModule,
} from '@nebular/theme';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { MapModule } from './components/map/map.module';
import { ConfirmationComponent } from './components/confirmation/confirmation.component';
import { ChartModule } from './components/chart/chart.module';
import { MessageMonitorComponent } from './components/message-monitor/message-monitor.component';
import { MessageValuePipe } from './pipes/message-value.pipe';
import { ToMillisecsPipe } from './pipes/time.pipe';
import { TableComponent } from './components/table/table.component';
import { PaginationComponent } from './components/pagination/pagination.component';
import { TaglistChipPipe } from 'app/shared/pipes/taglist-chip.pipe';
import { TagColorPipe } from 'app/shared/pipes/tag-color.pipe';
import { TagChipPipe } from 'app/shared/pipes/tag-chip.pipe';
import { ValidTagInputDirective } from 'app/shared/directives/valid-tag-input.directive';
import { AdvancedOptionsPipe } from 'app/shared/pipes/advanced-options.pipe';
import { PrettyJsonPipe} from 'app/shared/pipes/pretty-json.pipe';
import { TagControlComponent } from './components/orb/tag-control/tag-control.component';
import { AgentInformationComponent } from './components/orb/agent/information/agent-information.component';
import { AgentCapabilitiesComponent } from './components/orb/agent/capabilities/agent-capabilities.component';
import {
  AgentPoliciesDatasetsComponent,
} from './components/orb/agent/policies-datasets/agent-policies-datasets.component';
import { AgentGroupsComponent } from './components/orb/agent/groups/agent-groups.component';
import { AgentProvisioningComponent } from './components/orb/agent/provisioning/agent-provisioning.component';
import { ClipboardModule } from '@angular/cdk/clipboard';

@NgModule({
  imports: [
    ThemeModule,
    NbButtonModule,
    NbCardModule,
    NbDialogModule,
    MapModule,
    ChartModule,
    NbSelectModule,
    NbDatepickerModule,
    NbInputModule,
    NbAccordionModule,
    NbListModule,
    FormsModule,
    NbIconModule,
    NbCheckboxModule,
    MatChipsModule,
    MatIconModule,
    ReactiveFormsModule,
    ClipboardModule,
  ],
  declarations: [
    ConfirmationComponent,
    MessageMonitorComponent,
    MessageValuePipe,
    ToMillisecsPipe,
    TableComponent,
    PaginationComponent,
    AdvancedOptionsPipe,
    TagColorPipe,
    TagChipPipe,
    TaglistChipPipe,
    ValidTagInputDirective,
    PrettyJsonPipe,
    AgentInformationComponent,
    AgentCapabilitiesComponent,
    AgentPoliciesDatasetsComponent,
    AgentGroupsComponent,
    AgentProvisioningComponent,
    TagControlComponent,
  ],
  exports: [
    ThemeModule,
    NbCardModule,
    NbIconModule,
    MapModule,
    ChartModule,
    ConfirmationComponent,
    MessageMonitorComponent,
    TableComponent,
    PaginationComponent,
    AdvancedOptionsPipe,
    TagColorPipe,
    TagChipPipe,
    TaglistChipPipe,
    ValidTagInputDirective,
    PrettyJsonPipe,
    TagControlComponent,
    AgentInformationComponent,
    AgentCapabilitiesComponent,
    AgentPoliciesDatasetsComponent,
    AgentGroupsComponent,
    AgentProvisioningComponent,
  ],
  providers: [
    MessageValuePipe,
    ToMillisecsPipe,
    AdvancedOptionsPipe,
    TagColorPipe,
    TagChipPipe,
    TaglistChipPipe,
    ValidTagInputDirective,
  ],
})

export class SharedModule {
}
