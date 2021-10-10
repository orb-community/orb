import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-dataset-add-component',
  templateUrl: './dataset.add.component.html',
  styleUrls: ['./dataset.add.component.scss'],
})
export class DatasetAddComponent {
  // stepper vars
  detailsFormGroup: FormGroup;

  agentFormGroup: FormGroup;

  policyFormGroup: FormGroup;

  sinkFormGroup: FormGroup;

  dataset: Dataset;

  datasetID: string;

  availableAgentGroups = [];

  availableAgentPolicies = [];

  availableSinks = [];

  selectedSinks = [];

  isEdit: boolean;

  isLoading = false;

  datasetLoading = false;

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentPoliciesService: AgentPoliciesService,
    private datasetPoliciesService: DatasetPoliciesService,
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.dataset = this.router.getCurrentNavigation().extras.state?.dataset as Dataset || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.datasetID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.datasetID;
    this.datasetLoading = this.isEdit;

    !!this.datasetID && datasetPoliciesService.getDatasetById(this.datasetID).subscribe(resp => {
      this.dataset = resp;
      this.datasetLoading = false;
    });

    const { name, agent_group_id, agent_policy_id, sink_ids } = !!this.dataset ? this.dataset : {
      name: '',
      agent_group_id: { name: '', id: '' },
      agent_policy_id: { name: '', id: '' },
      sink_ids: [],
    };

    this.selectedSinks = sink_ids;

    this.detailsFormGroup = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
    });
    this.agentFormGroup = this._formBuilder.group({
      agent_group_id: [agent_group_id, [Validators.required]],
    });
    this.policyFormGroup = this._formBuilder.group({
      agent_policy_id: [agent_policy_id, [Validators.required]],
    });
    this.sinkFormGroup = this._formBuilder.group({
      selected_sink: ['', [Validators.required]],
    });
    // when editing, do not change agent group or policy
    if (!this.isEdit) {
      this.getAvailableAgentGroups();
      this.getAvailableAgentPolicies();
      this.getAvailableSinks();
    } else {
      this.agentFormGroup.controls.agent_group_id.disable();
      this.policyFormGroup.controls.agent_policy_id.disable();
    }

  }

  getAvailableAgentGroups() {
    this.isLoading = true;
    const pageInfo = { ...AgentGroupsService.getDefaultPagination(), limit: 100 };
    this.agentGroupsService
      .getAgentGroups(pageInfo, false)
      .subscribe((resp: OrbPagination<AgentGroup>) => {
        this.availableAgentGroups = resp.data;
        this.isLoading = false;
      });
  }

  getAvailableAgentPolicies() {
    this.isLoading = true;
    const pageInfo = { ...AgentPoliciesService.getDefaultPagination(), limit: 100 };
    this.agentPoliciesService
      .getAgentsPolicies(pageInfo, false)
      .subscribe((resp: OrbPagination<AgentPolicy>) => {
        this.availableAgentPolicies = resp.data;
        this.isLoading = false;
      });
  }

  getAvailableSinks() {
    this.isLoading = true;
    const pageInfo = { ...SinksService.getDefaultPagination(), limit: 100 };
    this.sinksService
      .getSinks(pageInfo, false)
      .subscribe((resp: OrbPagination<Sink>) => {
        const sinkIDMap = this.selectedSinks.map(sink => sink.id);
        this.availableSinks = resp.data.filter(sink => !sinkIDMap.includes(sink.id));
        this.isLoading = false;
      });
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/list');
  }

  onAgentGroupSelected(agentGroup: any) {
    this.agentFormGroup.controls.agent_group_id.patchValue(agentGroup);
  }

  onAddSink() {
    const sink = this.sinkFormGroup.controls.selected_sink.value;
    this.selectedSinks.push(sink);
    this.sinkFormGroup.controls.selected_sink.reset('');
    this.getAvailableSinks();
  }

  onRemoveSink(sink: any) {
    this.selectedSinks.splice(this.selectedSinks.indexOf(sink), 1);
    this.getAvailableSinks();
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFormGroup.controls.name.value,
      agent_group_id: this.agentFormGroup.controls.agent_group_id.value.id,
      agent_policy_id: this.policyFormGroup.controls.agent_policy_id.value.id,
      sink_ids: this.selectedSinks.map(sink => sink.id),
    } as Dataset;
    if (this.isEdit) {
      // updating existing dataset
      this.datasetPoliciesService.editDataset({ ...payload, id: this.datasetID }).subscribe(() => {
        this.notificationsService.success('Dataset successfully updated', '');
        this.goBack();
      });
    } else {
      this.datasetPoliciesService.addDataset(payload).subscribe(() => {
        this.notificationsService.success('Dataset successfully created', '');
        this.goBack();
      });
    }
  }
}
