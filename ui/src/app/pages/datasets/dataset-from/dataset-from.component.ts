import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { NbDialogRef, NbDialogService } from '@nebular/theme';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';
import { Observable, of } from 'rxjs';

const CONFIG = {
  SINKS: 'SINKS',
  GROUPS: 'GROUPS',
  POLICIES: 'POLICIES',
  DATASET: 'DATASET',
};

@Component({
  selector: 'ngx-dataset-from',
  templateUrl: './dataset-from.component.html',
  styleUrls: ['./dataset-from.component.scss'],
})
export class DatasetFromComponent implements OnInit {
  @Input()
  dataset: Dataset;

  @Input()
  policy: AgentPolicy;

  @Input()
  group: AgentGroup;

  isEdit: boolean;

  selectedGroup: string;

  selectedPolicy: string;
  sinkIDs: string[];
  availableAgentGroups: AgentGroup[];
  filteredAgentGroups$: Observable<AgentGroup[]>;
  availableAgentPolicies: AgentPolicy[];
  availableSinks: Sink[];
  unselectedSinks: Sink[];
  form: FormGroup;
  loading = Object.entries(CONFIG).reduce((acc, [value]) => {
    acc[value] = false;
    return acc;
  }, {});

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentPoliciesService: AgentPoliciesService,
    private datasetService: DatasetPoliciesService,
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private fb: FormBuilder,
    private dialogRef: NbDialogRef<DatasetFromComponent>,
    private dialogService: NbDialogService,
    private cdr: ChangeDetectorRef,
  ) {
    this.isEdit = false;
    this.availableAgentGroups = [];
    this.filteredAgentGroups$ = of(this.availableAgentGroups);
    this.availableAgentPolicies = [];
    this.availableSinks = [];
    this._selectedSinks = [];
    this.unselectedSinks = [];
    this.sinkIDs = [];

    this.getDatasetAvailableConfigList();

    this.readyForms();
  }

  private _selectedSinks: Sink[];

  get selectedSinks(): Sink[] {
    return this._selectedSinks;
  }

  // #load controls

  set selectedSinks(sinks: Sink[]) {
    this._selectedSinks = sinks;
    this.sinkIDs = sinks.map((sink) => sink.id);
    this.form.controls.sink_ids.patchValue(this.sinkIDs);
    this.form.controls.sink_ids.markAsDirty();
    this.cdr.markForCheck();
    this.updateUnselectedSinks();
  }

  readyForms() {
    const { name, agent_policy_id, agent_group_id, sink_ids } =
      this?.dataset ||
      ({
        name: '',
        agent_group_id: '',
        agent_policy_id: '',
        sink_ids: [],
      } as Dataset);

    this.form = this.fb.group({
      name: [
        name,
        [
          Validators.required,
          Validators.pattern(
            // https://github.com/ns1labs/orb/wiki/Architecture:-Common-Patterns#name-labels
            // anything starting with alpha chars or underscore followed by any
            // number of alphanumeric chars, dash '-' or underscore '_'. e.g.:
            // valid: my_name, _name0, name__anything invalid: 0something, 000,
            // 0_bla
            '^[a-zA-Z_][a-zA-Z0-9_-]*$',
          ),
        ],
      ],
      agent_policy_id: [agent_policy_id, [Validators.required]],
      agent_group_id: [agent_group_id, [Validators.required]],
      agent_group_name: [null, [this.groupNameValidator]],
      sink_ids: [sink_ids],
    });
  }

  groupNameValidator = (): ValidatorFn => {
    return (control: AbstractControl) =>
      this.availableAgentGroups.filter((agent) => agent.name === control.value)
        .length === 0
        ? { noMatch: 'Select a valid agent' }
        : null;
  }

  updateFormSelectedAgentGroupId(groupName) {
    const group = this.availableAgentGroups.filter(
      (agent) => agent.name === groupName,
    );
    let id;
    if (group.length > 0) {
      id = group[0].id;
    }
    this.form.patchValue({ agent_group_id: id });
    this.cdr.markForCheck();
  }

  updateFormSelectedAgentGroupName(groupId) {
    const group = this.availableAgentGroups.filter(
      (agent) => agent.id === groupId,
    );
    let name;
    if (group.length > 0) {
      name = group[0].name;
    }
    this.form.patchValue({ agent_group_name: name });
    this.cdr.markForCheck();
  }

  onChangeGroupName(event) {
    const value = event.currentTarget.value;
    this.onFilterGroup(value);
  }

  onSelectChangeGroupName(event) {
    this.onFilterGroup(event);
  }

  onFilterGroup(value) {
    this.filteredAgentGroups$ = of(this.filter(value));
  }

  ngOnInit(): void {
    if (!!this.group) {
      this.selectedGroup = this.group.id;
      this.form.patchValue({ agent_group_id: this.group.id });
      this.form.controls.agent_group_id.disable();
    }
    if (!!this.policy) {
      this.selectedPolicy = this.policy.id;
      this.form.patchValue({ agent_policy_id: this.policy.id });
      this.form.controls.agent_policy_id.disable();
    }
    if (!!this.dataset) {
      const { name, agent_group_id, agent_policy_id, sink_ids } = this.dataset;
      this.selectedGroup = agent_group_id;
      this.selectedSinks = this.availableSinks.filter((sink) =>
        sink_ids.includes(sink.id),
      );
      this.selectedPolicy = agent_policy_id;
      this.form.patchValue({ name, agent_group_id, agent_policy_id, sink_ids });
      this.isEdit = true;
      this.form.controls.agent_group_id.disable();
      this.form.controls.agent_policy_id.disable();

      this.unselectedSinks = this.availableSinks.filter(
        (sink) => !this._selectedSinks.includes(sink),
      );
    }
  }

  updateUnselectedSinks() {
    this.unselectedSinks = this.availableSinks.filter(
      (sink) => !this._selectedSinks.includes(sink),
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
          // console.log('warning');
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
          this.availableAgentGroups = resp.sort((a, b) =>
            a.name > b.name ? -1 : 1,
          );
          this.filteredAgentGroups$ = of(this.availableAgentGroups);
          this.loading[CONFIG.GROUPS] = false;

          if (this.dataset?.agent_group_id) {
            this.updateFormSelectedAgentGroupName(this.dataset.agent_group_id);
          }
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
        this._selectedSinks.forEach((sink) => {
          sink.name = resp.find(
            (anotherSink) => anotherSink.id === sink.id,
          ).name;
        });

        this.availableSinks = resp;
        this.updateUnselectedSinks();

        this.loading[CONFIG.SINKS] = false;

        resolve(this.availableSinks);
      });
    });
  }

  isLoading() {
    return Object.values<boolean>(this.loading).reduce(
      (prev, curr) => prev && curr,
    );
  }

  onFormSubmit() {
    const payload = {
      name: this.form.controls.name.value,
      agent_group_id: this.form.controls.agent_group_id.value,
      agent_policy_id: this.form.controls.agent_policy_id.value,
      sink_ids: this._selectedSinks.map((sink) => sink.id),
    } as Dataset;
    if (this.isEdit) {
      // updating existing dataset
      this.datasetService
        .editDataset({ ...payload, id: this.dataset.id })
        .subscribe(() => {
          this.notificationsService.success('Dataset successfully updated', '');
          this.dialogRef.close('edited');
        });
    } else {
      this.datasetService.addDataset(payload).subscribe(() => {
        this.notificationsService.success('Dataset successfully created', '');
        this.dialogRef.close('created');
      });
    }
  }

  onDelete() {
    this.dialogService
      .open(DatasetDeleteComponent, {
        context: { name: this.dataset.name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.datasetService.deleteDataset(this.dataset.id).subscribe(() => {
            this.notificationsService.success(
              'Dataset successfully deleted',
              '',
            );
            this.dialogRef.close('deleted');
          });
        }
      });
  }

  onClose() {
    this.dialogRef.close('canceled');
  }

  private filter(value: string): AgentGroup[] {
    let filtered;
    if (value === '') {
      filtered = this.availableAgentGroups;
    } else {
      filtered = this.availableAgentGroups.filter((group) =>
        group.name.includes(value),
      );
    }
    this.updateFormSelectedAgentGroupId(value);

    return filtered;
  }
}
