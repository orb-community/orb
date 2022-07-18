import { ClipboardModule } from '@angular/cdk/clipboard';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import {
  NbAccordionModule,
  NbAutocompleteModule,
  NbButtonModule,
  NbCardModule,
  NbCheckboxModule,
  NbDatepickerModule,
  NbDialogModule,
  NbFormFieldModule,
  NbIconModule,
  NbInputModule,
  NbListModule,
  NbSelectModule,
  NbTooltipModule,
} from '@nebular/theme';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';

import { ThemeModule } from 'app/@theme/theme.module';
import { SinkDisplayComponent } from 'app/shared/components/orb/sink-display/sink-display.component';
import { SinkDisplayListComponent } from 'app/shared/components/orb/sink/sink-display/sink-display-list.component';
import { ValidTagInputDirective } from 'app/shared/directives/valid-tag-input.directive';
import { AdvancedOptionsPipe } from 'app/shared/pipes/advanced-options.pipe';
import { PrettyJsonPipe } from 'app/shared/pipes/pretty-json.pipe';
import { SortPipe } from 'app/shared/pipes/sort.pipe';
import { TagChipPipe } from 'app/shared/pipes/tag-chip.pipe';
import { TagColorPipe } from 'app/shared/pipes/tag-color.pipe';
import { TaglistChipPipe } from 'app/shared/pipes/taglist-chip.pipe';
import { MonacoEditorModule } from 'ngx-monaco-editor';
import { ConfirmationComponent } from './components/confirmation/confirmation.component';
import { FilterComponent } from './components/filter/filter.component';
import { AgentCapabilitiesComponent } from './components/orb/agent/agent-capabilities/agent-capabilities.component';
import { AgentGroupsComponent } from './components/orb/agent/agent-groups/agent-groups.component';
import { AgentInformationComponent } from './components/orb/agent/agent-information/agent-information.component';
import { AgentPoliciesDatasetsComponent } from './components/orb/agent/agent-policies-datasets/agent-policies-datasets.component';
import { AgentProvisioningComponent } from './components/orb/agent/agent-provisioning/agent-provisioning.component';
import { CombinedTagComponent } from './components/orb/combined-tag/combined-tag.component';
import { GroupedAgentsComponent } from './components/orb/dataset/grouped-agents/grouped-agents.component';
import { PolicyDatasetsComponent } from './components/orb/policy/policy-datasets/policy-datasets.component';
import { PolicyDetailsComponent } from './components/orb/policy/policy-details/policy-details.component';
import { PolicyGroupsComponent } from './components/orb/policy/policy-groups/policy-groups.component';
import { PolicyInterfaceComponent } from './components/orb/policy/policy-interface/policy-interface.component';
import { SinkControlComponent } from './components/orb/sink-control/sink-control.component';
import { TagControlComponent } from './components/orb/tag-control/tag-control.component';
import { TagDisplayComponent } from './components/orb/tag-display/tag-display.component';
import { PaginationComponent } from './components/pagination/pagination.component';
import { TableComponent } from './components/table/table.component';
import { MessageValuePipe } from './pipes/message-value.pipe';
import { PrettyYamlPipe } from './pipes/pretty-yaml.pipe';
import { ToMillisecsPipe } from './pipes/time.pipe';
import { PollControlComponent } from './components/poll-control/poll-control.component';

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
    MatInputModule,
    MatButtonModule,
    MatFormFieldModule,
    MatSelectModule,
    MatAutocompleteModule,
    MatCheckboxModule,
    MatToolbarModule,
    NbAutocompleteModule,
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
    CombinedTagComponent,
    SortPipe,
    FilterComponent,
    PollControlComponent,
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
    CombinedTagComponent,
    SortPipe,
    FilterComponent,
    PollControlComponent,
  ],
  providers: [
    MessageValuePipe,
    ToMillisecsPipe,
    AdvancedOptionsPipe,
    TagColorPipe,
    TagChipPipe,
    TaglistChipPipe,
    ValidTagInputDirective,
    SortPipe,
  ],
})
export class SharedModule {}
