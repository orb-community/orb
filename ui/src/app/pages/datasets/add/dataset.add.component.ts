import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';

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
    private datasetPoliciesService: DatasetPoliciesService,
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
      this.getAvailableAgentGroups();
    });
    !this.datasetLoading && this.getAvailableAgentGroups();
  }

  getAvailableAgentGroups() {
    this.isLoading = true;
    this.datasetPoliciesService.getAvailableAgentGroups().subscribe(agentGroups => {
      this.availableAgentGroups = agentGroups.map(entry => entry.backend);
      this.customSinkSettings = this.availableAgentGroups.reduce((accumulator, curr) => {
        const index = agentGroups.findIndex(entry => entry.backend === curr);
        accumulator[curr] = agentGroups[index].config.map(entry => ({
          type: entry.type,
          label: entry.title,
          prop: entry.name,
          input: entry.input,
          required: entry.required,
        }));
        return accumulator;
      }, {});
      const {name, tags} = !!this.dataset ? this.dataset : {

      } as Dataset;
      this.firstFormGroup = this._formBuilder.group({
        name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      });

      // builds secondFormGroup
      // this.onAgentGroupSelected(backend);

      // this.thirdFormGroup = this._formBuilder.group({
      //   tags: [Object.keys(tags || {}).map(key => ({[key]: tags[key]})),
      //     Validators.minLength(1)],
      //   key: [''],
      //   value: [''],
      // });

      this.isLoading = false;
    });
  }

  goBack() {
    // this.router.navigateByUrl('/pages/datasets/list');
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
