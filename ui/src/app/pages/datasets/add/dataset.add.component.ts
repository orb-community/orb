import { Component, ElementRef, ViewChild } from '@angular/core';

import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { STRINGS } from '../../../../assets/text/strings';

const CONFIG = {
  SINKS: 'SINKS',
  GROUPS: 'GROUPS',
  POLICIES: 'POLICIES',
  DATASET: 'DATASET',
};

@Component({
  selector: 'ngx-dataset-add-component',
  templateUrl: './dataset.add.component.html',
  styleUrls: ['./dataset.add.component.scss'],
})
export class DatasetAddComponent {
  @ViewChild('sinkSelLead') sinkSelLead: ElementRef;

  // page vars
  strings = { stepper: STRINGS.stepper };

  // stepper form groups
  detailsFormGroup: FormGroup;

  sinkFormGroup: FormGroup;

  selectedGroup: AgentGroup;

  selectedPolicy: AgentPolicy;

  // stores user selected sinks
  selectedSinks: { id: string; name?: string }[] = [];

  isEdit: boolean;

  // #load controls
  loading = Object.entries(CONFIG).reduce((acc, [value]) => {
    acc[value] = false;
    return acc;
  }, {});

  datasetID: string;

  dataset: Dataset;

  availableAgentGroups: AgentGroup[] = [];

  availableAgentPolicies: AgentPolicy[] = [];

  availableSinks: Sink[] = [];

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
    this.readyForms();

    this.getDatasetAvailableConfigList();

    this.dataset =
      (this.router.getCurrentNavigation().extras.state?.dataset as Dataset) ||
      null;
    this.datasetID = this.route.snapshot.paramMap.get('id');
    this.isEdit = !!this.datasetID;
    this.loading[CONFIG.DATASET] = this.isEdit;

    !!this.datasetID &&
      datasetPoliciesService
        .getDatasetById(this.datasetID)
        .subscribe((resp) => {
          this.dataset = resp;
          this.loading[CONFIG.DATASET] = false;
          this.updateForms();
        });

    this.updateForms();
  }

  readyForms() {
    const { name, sink_ids } = (this.dataset = {
      name: '',
      agent_group_id: '',
      agent_policy_id: '',
      sink_ids: [],
    } as Dataset);

    this.selectedSinks = sink_ids.map<{ id: string; name: string }>((id) => {
      return {
        id,
        name:
          this.availableSinks.length > 0
            ? this.availableSinks.find((sink) => sink.id === id)?.name
            : '',
      };
    });

    this.detailsFormGroup = this._formBuilder.group({
      name: [
        name,
        [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')],
      ],
    });
    this.sinkFormGroup = this._formBuilder.group({
      selected_sink: ['', [Validators.minLength(1)]],
    });
  }

  updateForms() {
    const {
      name,
      agent_group_id,
      agent_policy_id,
      sink_ids,
    } = (this.dataset = {
      name: '',
      agent_group_id: '',
      agent_policy_id: '',
      sink_ids: [],
      ...this.dataset,
    } as Dataset);

    this.selectedSinks = sink_ids.map<{ id: string; name: string }>((id) => {
      return {
        id,
        name:
          this.availableSinks.length > 0
            ? this.availableSinks.find((sink) => sink.id === id)?.name
            : '',
      };
    });
    if (this.availableSinks.length > 0 && this.selectedSinks.length > 0)
      this.availableSinks = this.availableSinks.filter(
        (sink) => !this.selectedSinks.includes({ id: sink.id }),
      );

    this.loading[CONFIG.SINKS] = true;
    this.getAvailableSinks().catch((reason) =>
      console.warn(`Couldn't fetch available sinks. Reason: ${reason}`),
    );

    this.selectedGroup =
      !!agent_group_id &&
      this.availableAgentGroups.find((agent) => agent.id === agent_group_id);
    this.selectedPolicy =
      !!agent_policy_id &&
      this.availableAgentPolicies.find(
        (policy) => policy.id === agent_policy_id,
      );

    this.detailsFormGroup.controls.name.patchValue(name);
    this.detailsFormGroup.updateValueAndValidity();
  }

  isLoading() {
    return Object.values<boolean>(this.loading).reduce(
      (prev, curr) => prev && curr,
    );
  }

  getDatasetAvailableConfigList() {
    Promise.all([
      this.getAvailableAgentGroups(),
      this.getAvailableAgentPolicies(),
      this.getAvailableSinks(),
    ])
      .then(
        (value) => {
          if (this.isEdit && this.dataset) {
            this.updateForms();
          }
        },
        (reason) =>
          console.warn(
            `Cannot retrieve available configurations - reason: ${JSON.parse(
              reason,
            )}`,
          ),
      )
      .catch((reason) => {
        console.warn(
          `Cannot retrieve backend data - reason: ${JSON.parse(reason)}`,
        );
      });
  }

  getAvailableAgentGroups() {
    return new Promise((resolve) => {
      this.loading[CONFIG.GROUPS] = true;
      this.agentGroupsService
        .getAllAgentGroups()
        .subscribe((resp: AgentGroup[]) => {
          this.availableAgentGroups = resp;
          this.loading[CONFIG.GROUPS] = false;

          resolve(this.availableAgentGroups);
        });
    });
  }

  getAvailableAgentPolicies() {
    return new Promise((resolve) => {
      this.loading[CONFIG.POLICIES] = true;

      this.agentPoliciesService
        .getAllAgentPolicies()
        .subscribe((resp: AgentPolicy[]) => {
          this.availableAgentPolicies = resp;
          this.loading[CONFIG.POLICIES] = false;

          resolve(this.availableAgentPolicies);
        });
    });
  }

  getAvailableSinks() {
    return new Promise((resolve) => {
      this.loading[CONFIG.SINKS] = true;
      this.sinksService.getAllSinks().subscribe((resp: Sink[]) => {
        this.selectedSinks.forEach((sink) => {
          sink.name = resp.find(
            (anotherSink) => anotherSink.id === sink.id,
          ).name;
        });
        const sinkIDMap = this.selectedSinks.map((sink) => sink.id);
        this.availableSinks = resp.filter(
          (sink) => !sinkIDMap.includes(sink.id),
        );
        this.loading[CONFIG.SINKS] = false;

        resolve(this.availableSinks);
      });
    });
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/list');
  }

  onAddSink() {
    const sink = this.sinkFormGroup.controls.selected_sink.value;
    this.selectedSinks.push(sink);
    this.sinkFormGroup.controls.selected_sink.reset('');
    this.sinkSelLead.nativeElement.focus();
    this.getAvailableSinks();
  }

  onRemoveSink(sink: any) {
    this.selectedSinks.splice(this.selectedSinks.indexOf(sink), 1);
    this.getAvailableSinks();
  }

  onPolicySelected(policy) {
    this.selectedPolicy = policy;
    this.agentPoliciesService
      .getAgentPolicyById(policy.id)
      .subscribe((fullPolicy) => {
        this.selectedPolicy = { ...policy, ...fullPolicy };
      });
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFormGroup.controls.name.value,
      agent_group_id: this.selectedGroup.id,
      agent_policy_id: this.selectedPolicy.id,
      sink_ids: this.selectedSinks.map((sink) => sink.id),
    } as Dataset;
    if (this.isEdit) {
      // updating existing dataset
      this.datasetPoliciesService
        .editDataset({ ...payload, id: this.datasetID })
        .subscribe(() => {
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
