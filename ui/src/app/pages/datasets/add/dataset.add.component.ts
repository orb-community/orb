import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup } from '@angular/forms';
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
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  thirdFormGroup: FormGroup;

  dataset: Dataset;

  datasetID: string;

  availableAgentGroups = [];

  availableAgentPolicies = [];

  availableSinks = [];

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

    // when editing, do not change agent group or policy
    if (!this.isEdit) {
      this.getAvailableAgentGroups();
      this.getAvailableAgentPolicies();
      this.getAvailableSinks();
    }

  }

  getAvailableAgentGroups() {
    this.isLoading = true;
    // todo how to request with infinite limit?
    const pageInfo = { ...AgentGroupsService.getDefaultPagination(), limit: 999 };
    // todo check if any filter should be applied for available Agent Groups
    this.agentGroupsService
      .getAgentGroups(pageInfo, false)
      .subscribe((resp: OrbPagination<AgentGroup>) => {
        this.availableAgentGroups = resp.data;
        this.isLoading = false;
      });
  }

  getAvailableAgentPolicies() {
    this.isLoading = true;
    // todo how to request with infinite limit?
    const pageInfo = { ...AgentGroupsService.getDefaultPagination(), limit: 999 };
    // todo check if any filter should be applied for available Agent Groups
    this.agentPoliciesService
      .getAgentsPolicies(pageInfo, false)
      .subscribe((resp: OrbPagination<AgentPolicy>) => {
        this.availableAgentPolicies = resp.data;
        this.isLoading = false;
      });
  }

  getAvailableSinks() {
    this.isLoading = true;
    // todo how to request with infinite limit?
    const pageInfo = { ...AgentGroupsService.getDefaultPagination(), limit: 999 };
    // todo check if any filter should be applied for available Agent Groups
    this.agentPoliciesService
      .getAgentsPolicies(pageInfo, false)
      .subscribe((resp: OrbPagination<Sink>) => {
        this.availableSinks = resp.data;
        this.isLoading = false;
      });
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/list');
  }

  onFormSubmit() {

  }

  onAgentGroupSelected(selectedValue) {

  }

  onAddSink() {

  }

  onRemoveSink(sink: any) {

  }
}
