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
    NbDialogModule, NbFormFieldModule,
    NbIconModule,
    NbInputModule,
    NbListModule,
    NbSelectModule, NbTooltipModule,
} from '@nebular/theme';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { ConfirmationComponent } from './components/confirmation/confirmation.component';
import { MessageValuePipe } from './pipes/message-value.pipe';
import { ToMillisecsPipe } from './pipes/time.pipe';
import { TableComponent } from './components/table/table.component';
import { PaginationComponent } from './components/pagination/pagination.component';
import { TaglistChipPipe } from 'app/shared/pipes/taglist-chip.pipe';
import { TagColorPipe } from 'app/shared/pipes/tag-color.pipe';
import { TagChipPipe } from 'app/shared/pipes/tag-chip.pipe';
import { ValidTagInputDirective } from 'app/shared/directives/valid-tag-input.directive';
import { AdvancedOptionsPipe } from 'app/shared/pipes/advanced-options.pipe';
import { PrettyJsonPipe } from 'app/shared/pipes/pretty-json.pipe';
import { TagControlComponent } from './components/orb/tag-control/tag-control.component';
import { AgentInformationComponent } from './components/orb/agent/agent-information/agent-information.component';
import { AgentCapabilitiesComponent } from './components/orb/agent/agent-capabilities/agent-capabilities.component';
import {
  AgentPoliciesDatasetsComponent,
} from './components/orb/agent/agent-policies-datasets/agent-policies-datasets.component';
import { AgentGroupsComponent } from './components/orb/agent/agent-groups/agent-groups.component';
import { AgentProvisioningComponent } from './components/orb/agent/agent-provisioning/agent-provisioning.component';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { TagDisplayComponent } from './components/orb/tag-display/tag-display.component';
import { MatTooltipModule } from '@angular/material/tooltip';
import { PolicyDetailsComponent } from './components/orb/policy/policy-details/policy-details.component';
import { PolicyInterfaceComponent } from './components/orb/policy/policy-interface/policy-interface.component';
import { PolicyDatasetsComponent } from './components/orb/policy/policy-datasets/policy-datasets.component';
import { GroupedAgentsComponent } from './components/orb/dataset/grouped-agents/grouped-agents.component';
import { PrettyYamlPipe } from './pipes/pretty-yaml.pipe';
import { SinkControlComponent } from './components/orb/sink-control/sink-control.component';
import { PolicyGroupsComponent } from './components/orb/policy/policy-groups/policy-groups.component';
import { SinkDisplayComponent } from 'app/shared/components/orb/sink-display/sink-display.component';
import { SinkDisplayListComponent } from 'app/shared/components/orb/sink/sink-display/sink-display-list.component';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { MonacoEditorModule } from 'ngx-monaco-editor';

@NgModule({
  imports: [
    ThemeModule,
    NbButtonModule,
    NbCardModule,
    NbDialogModule,
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
    MatTooltipModule,
    NgxDatatableModule,
    NbTooltipModule,
    NbFormFieldModule,
    MonacoEditorModule,
  ],
  declarations: [
    ConfirmationComponent,
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
    TagDisplayComponent,
    PolicyDetailsComponent,
    PolicyInterfaceComponent,
    PolicyDatasetsComponent,
    PolicyGroupsComponent,
    GroupedAgentsComponent,
    PrettyYamlPipe,
    SinkControlComponent,
    SinkDisplayComponent,
    SinkDisplayListComponent,
  ],
  exports: [
    ThemeModule,
    NbCardModule,
    NbIconModule,
    ConfirmationComponent,
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
    TagDisplayComponent,
    PolicyDetailsComponent,
    PolicyInterfaceComponent,
    PolicyDatasetsComponent,
    GroupedAgentsComponent,
    PolicyGroupsComponent,
    PrettyYamlPipe,
    SinkControlComponent,
    SinkDisplayComponent,
    SinkDisplayListComponent,
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
