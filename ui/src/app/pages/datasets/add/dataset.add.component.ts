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

  agentGroupFormGroup: FormGroup;

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
    this.detailsFormGroup = this._formBuilder.group({
      name: ['', [Validators.required]],
    });
    this.agentGroupFormGroup = this._formBuilder.group({
      agent_group_id: ['', [Validators.required]],
    });
    this.policyFormGroup = this._formBuilder.group({
      agent_policy_id: ['', [Validators.required]],
    });
    this.sinkFormGroup = this._formBuilder.group({
      selected_sink: ['', [Validators.required]],
    });
    this.dataset = this.router.getCurrentNavigation().extras.state?.dataset as Dataset || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.datasetID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.datasetID;
    this.datasetLoading = this.isEdit;

    !!this.datasetID && datasetPoliciesService.getDatasetById(this.datasetID).subscribe(resp => {
      this.dataset = resp;
      this.datasetLoading = false;
    });

    // when editing, do not change agent group or policy
    if (!this.isEdit) {
      this.getAvailableAgentGroups();
      this.getAvailableAgentPolicies();
      this.getAvailableSinks();
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
        this.availableSinks = resp.data.filter(sink => !this.selectedSinks.includes(sink));
        this.isLoading = false;
      });
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/list');
  }

  onAgentGroupSelected(agentGroup: any) {
    this.agentGroupFormGroup.controls.agent_group_id.patchValue(agentGroup);
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

  }
}
