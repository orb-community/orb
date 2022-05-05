import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NbDialogRef, NbDialogService } from '@nebular/theme';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';

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

  get selectedSinks(): Sink[] {
    return this._selectedSinks;
  }

  set selectedSinks(sinks: Sink[]) {
    this._selectedSinks = sinks;
    this.sinkIDs = sinks.map(sink => sink.id);
    this.form.controls.sink_ids.patchValue(this.sinkIDs);
    this.form.controls.sink_ids.markAsDirty();
    this.cdr.markForCheck();
    this.updateUnselectedSinks();
  }

  private _selectedSinks: Sink[];

  sinkIDs: string[];

  availableAgentGroups: AgentGroup[];

  availableAgentPolicies: AgentPolicy[];

  availableSinks: Sink[];

  unselectedSinks: Sink[];

  form: FormGroup;

  loading = Object.entries(CONFIG)
    .reduce((acc, [value]) => {
      acc[value] = false;
      return acc;
    }, {});

  // #load controls

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
    this.availableAgentPolicies = [];
    this.availableSinks = [];
    this._selectedSinks = [];
    this.unselectedSinks = [];
    this.sinkIDs = [];

    this.getDatasetAvailableConfigList();

    this.readyForms();
  }

  readyForms() {
    const {
      name,
      agent_policy_id,
      agent_group_id,
      sink_ids,
    } = this?.dataset || {
      name: '',
      agent_group_id: '',
      agent_policy_id: '',
      sink_ids: [],
    } as Dataset;

    this.form = this.fb.group({
      name: [
        name, [
          Validators.required, Validators.pattern(
            // https://github.com/ns1labs/orb/wiki/Architecture:-Common-Patterns#name-labels
            // anything starting with alpha chars or underscore followed by any
            // number of alphanumeric chars, dash '-' or underscore '_'. e.g.:
            // valid: my_name, _name0, name__anything invalid: 0something, 000,
            // 0_bla
            '^[a-zA-Z_][a-zA-Z0-9_-]*$'),
        ],
      ],
      agent_policy_id: [
        agent_policy_id, [Validators.required],
      ],
      agent_group_id: [
        agent_group_id, [Validators.required],
      ],
      sink_ids: [sink_ids],
    });
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
      this.selectedSinks = this.availableSinks.filter(
        sink => sink_ids.includes(sink.id));
      this.selectedPolicy = agent_policy_id;
      this.form.patchValue({ name, agent_group_id, agent_policy_id, sink_ids });
      this.isEdit = true;
      this.form.controls.agent_group_id.disable();
      this.form.controls.agent_policy_id.disable();

      this.unselectedSinks = this.availableSinks.filter(
        sink => !this._selectedSinks.includes(sink));
    }
  }

  updateUnselectedSinks() {
    this.unselectedSinks = this.availableSinks.filter(
      sink => !this._selectedSinks.includes(sink));
  }

  getDatasetAvailableConfigList() {
    Promise.all([
        this.getAvailableAgentGroups(),
        this.getAvailableAgentPolicies(), this.getAvailableSinks(),
      ])
      .then(value => {
        // console.log('warning');
      }, reason => console.warn(
        `Cannot retrieve available configurations - reason: ${JSON.parse(
          reason)}`))
      .catch(reason => {
        console.warn(
          `Cannot retrieve backend data - reason: ${JSON.parse(reason)}`);
      });
  }

  getAvailableAgentGroups() {
    return new Promise((resolve) => {
      this.loading[CONFIG.GROUPS] = true;
      this.agentGroupsService.getAllAgentGroups()
        .subscribe((resp: OrbPagination<AgentGroup>) => {
          this.availableAgentGroups = resp.data;
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
        .subscribe((resp: OrbPagination<AgentPolicy>) => {
          this.availableAgentPolicies = resp.data;
          this.loading[CONFIG.POLICIES] = false;

          resolve(this.availableAgentPolicies);
        });
    });
  }

  getAvailableSinks() {
    return new Promise((resolve) => {
      this.loading[CONFIG.SINKS] = true;
      const pageInfo = { ...SinksService.getDefaultPagination(), limit: 100 };
      this.sinksService
        .getSinks(pageInfo, false)
        .subscribe((resp: OrbPagination<Sink>) => {
          this._selectedSinks.forEach((sink) => {
            sink.name = resp.data.find(
              anotherSink => anotherSink.id === sink.id).name;
          });

          this.availableSinks = resp.data;
          this.updateUnselectedSinks();

          this.loading[CONFIG.SINKS] = false;

          resolve(this.availableSinks);
        });
    });
  }

  isLoading() {
    return Object.values<boolean>(this.loading)
      .reduce((prev, curr) => prev && curr);
  }

  onFormSubmit() {
    const payload = {
      name: this.form.controls.name.value,
      agent_group_id: this.form.controls.agent_group_id.value,
      agent_policy_id: this.form.controls.agent_policy_id.value,
      sink_ids: this._selectedSinks.map(sink => sink.id),
    } as Dataset;
    if (this.isEdit) {
      // updating existing dataset
      this.datasetService.editDataset({ ...payload, id: this.dataset.id })
        .subscribe(() => {
          this.notificationsService.success('Dataset successfully updated',
            '');
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
    this.dialogService.open(DatasetDeleteComponent, {
      context: { name: this.dataset.name },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.datasetService.deleteDataset(this.dataset.id).subscribe(() => {
            this.notificationsService.success('Dataset successfully deleted',
              '');
            this.dialogRef.close('deleted');
          });
        }
      },
    );
  }

  onClose() {
    this.dialogRef.close('canceled');
  }

}
