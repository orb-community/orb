import { Component, ViewChild } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { DynamicFormConfig } from 'app/common/interfaces/orb/dynamic.form.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { PolicyTap } from 'app/common/interfaces/orb/policy/policy.tap.interface';
import { NbDialogService } from '@nebular/theme';
import { HandlerPolicyAddComponent } from 'app/pages/datasets/policies.agent/add/handler.policy.add.component';
import { STRINGS } from '../../../../../assets/text/strings';

const CONFIG = {
  TAPS: 'TAPS',
  BACKEND: 'BACKEND',
  INPUTS: 'INPUTS',
  HANDLERS: 'HANDLERS',
  AGENT_POLICY: 'AGENT_POLICY',
};

@Component({
  selector: 'ngx-agent-policy-add-component',
  templateUrl: './agent.policy.add.component.html',
  styleUrls: ['./agent.policy.add.component.scss'],
})
export class AgentPolicyAddComponent {
  strings = { stepper: STRINGS.stepper };

  // #forms
  // agent policy general information - name, desc, backend
  detailsFG: FormGroup;

  // selected tap, input_type
  tapFG: FormGroup;

  // dynamic input config
  inputConfigFG: FormGroup;

  // dynamic input filter config
  inputFilterFG: FormGroup;

  // #key inputs holders
  // selected backend object
  backend: { [propName: string]: any };

  // selected tap object
  tap: PolicyTap;

  // selected input object
  input: {
    version?: string;
    config?: DynamicFormConfig;
    filter?: DynamicFormConfig;
  };

  // holds all handlers added by user
  modules: {
    [propName: string]: {
      name?: string;
      type?: string;
      config?: { [propName: string]: {} | any };
      filter?: { [propName: string]: {} | any };
    };
  } = {};

  // #services responses
  // hold info retrieved
  availableBackends: {
    [propName: string]: {
      backend: string;
      description: string;
    };
  };

  availableTaps: { [propName: string]: PolicyTap };

  availableInputs: {
    [propName: string]: {
      version?: string;
      config?: DynamicFormConfig;
      filter?: DynamicFormConfig;
    };
  };

  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  @ViewChild('editorComponent')
  editor;

  isEdit: boolean;

  editorOptions = {
    theme: 'vs-dark',
    language: 'yaml',
    automaticLayout: true,
    glyphMargin: false,
    folding: true,
    // Undocumented see https://github.com/Microsoft/vscode/issues/30795#issuecomment-410998882
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
  };

  code = `handlers:
  modules:
    default_dns:
      type: dns
    default_net:
      type: net
input:
  input_type: pcap
  tap: default_pcap
kind: collection`;

  // is config specified wizard mode or in YAML or JSON
  isWizard = true;

  // format definition
  format = 'yaml';

  // #load controls
  isLoading = Object.entries(CONFIG).reduce((acc, [value]) => {
    acc[value] = false;
    return acc;
  }, {}) as { [propName: string]: boolean };

  constructor(
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
    private dialogService: NbDialogService,
  ) {
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.newAgent();
    this.isEdit = !!this.agentPolicyID;

    this.readyForms();

    Promise.all([
      this.isEdit ? this.retrieveAgentPolicy() : Promise.resolve(),
      this.getBackendsList(),
    ])
      .catch((reason) => console.warn(`Couldn't fetch data. Reason: ${reason}`))
      .then(() => this.updateForms())
      .catch((reason) =>
        console.warn(
          `Couldn't fetch ${this.agentPolicy?.backend} data. Reason: ${reason}`,
        ),
      );
  }

  resizeComponents() {
    const timeoutId = setTimeout(() => {
      window.dispatchEvent(new Event('resize'));
      clearTimeout(timeoutId);
    }, 50);
    !!this.editor?.layout && this.editor.layout();
  }

  newAgent() {
    return {
      name: '',
      description: '',
      backend: 'pktvisor',
      tags: {},
      version: 1,
      policy: {
        kind: 'collection',
        input: {
          config: {},
          tap: '',
          input_type: '',
        },
        handlers: {
          modules: {},
        },
      },
    } as AgentPolicy;
  }

  retrieveAgentPolicy() {
    return new Promise((resolve) => {
      this.agentPoliciesService
        .getAgentPolicyById(this.agentPolicyID)
        .subscribe((policy) => {
          this.agentPolicy = policy;
          this.isLoading[CONFIG.AGENT_POLICY] = false;
          resolve(policy);
        });
    });
  }

  isLoadComplete() {
    return !Object.values(this.isLoading).reduce(
      (prev, curr) => prev || curr,
      false,
    );
  }

  readyForms() {
    const {
      name: name,
      description,
      backend,
      policy_data,
      policy: {
        input: { tap, input_type },
        handlers: { modules },
      },
    } = this.agentPolicy;

    if (policy_data) {
      this.code = policy_data;
    }

    this.modules = modules;

    this.detailsFG = this._formBuilder.group({
      name: [
        name,
        [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')],
      ],
      description: [description],
      backend: [
        { value: backend, disabled: backend !== '' },
        [Validators.required],
      ],
    });
    this.tapFG = this._formBuilder.group({
      selected_tap: [tap, Validators.required],
      input_type: [input_type, Validators.required],
    });
  }

  updateForms() {
    const {
      name: name,
      description,
      backend,
      format,
      policy_data,
      policy: { handlers },
    } = this.agentPolicy;

    const wizard = format !== this.format;

    if (policy_data) {
      this.isWizard = false;
      this.code = policy_data;
    }

    this.detailsFG.patchValue({ name, description, backend });

    this.modules = handlers?.modules || {};

    if (wizard) {
      this.onBackendSelected(backend).catch((reason) =>
        console.warn(`${reason}`),
      );
    }
  }

  getBackendsList() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.BACKEND] = true;
      this.agentPoliciesService.getAvailableBackends().subscribe((backends) => {
        this.availableBackends =
          !!backends &&
          backends.reduce((acc, curr) => {
            acc[curr.backend] = curr;
            return acc;
          }, {});

        this.isLoading[CONFIG.BACKEND] = false;

        resolve(backends);
      });
    });
  }

  onBackendSelected(selectedBackend) {
    return new Promise((resolve) => {
      this.backend = this.availableBackends[selectedBackend];
      this.backend['config'] = {};

      // todo hardcoded for pktvisor
      this.getBackendData().then(() => {
        resolve(null);
      });
    });
  }

  getBackendData() {
    return Promise.all([this.getTaps(), this.getInputs()])
      .then(
        (value) => {
          if (this.isEdit && this.agentPolicy && this.isWizard) {
            const selected_tap = this.agentPolicy.policy.input.tap;
            this.tapFG.patchValue({ selected_tap }, { emitEvent: true });
            this.onTapSelected(selected_tap);
            this.tapFG.controls.selected_tap.disable();
          }
        },
        (reason) =>
          console.warn(
            `Cannot retrieve backend data - reason: ${JSON.parse(reason)}`,
          ),
      )
      .catch((reason) => {
        console.warn(
          `Cannot retrieve backend data - reason: ${JSON.parse(reason)}`,
        );
      });
  }

  getTaps() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.TAPS] = true;
      this.agentPoliciesService
        .getBackendConfig([this.backend.backend, 'taps'])
        .subscribe((taps) => {
          this.availableTaps = taps.reduce((acc, curr) => {
            acc[curr.name] = curr;
            return acc;
          }, {});

          this.isLoading[CONFIG.TAPS] = false;

          resolve(taps);
        });
    });
  }

  onTapSelected(selectedTap) {
    this.tap = this.availableTaps[selectedTap];
    this.tapFG.controls.selected_tap.patchValue(selectedTap);

    const { input } = this.agentPolicy.policy;
    const { input_type, config_predefined, filter_predefined } = this.tap;

    this.tap.config = {
      ...config_predefined,
      ...input.config,
    };

    this.tap.filter = {
      ...filter_predefined,
      ...input.filter,
    };

    if (input_type) {
      this.onInputSelected(input_type);
    } else {
      this.input = null;
      this.tapFG.controls.input_type.reset('');
    }
  }

  getInputs() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.INPUTS] = true;
      this.agentPoliciesService
        .getBackendConfig([this.backend.backend, 'inputs'])
        .subscribe((inputs) => {
          this.availableInputs = !!inputs && inputs;

          this.isLoading[CONFIG.INPUTS] = false;

          resolve(inputs);
        });
    });
  }

  onInputSelected(input_type) {
    // TODO version here
    this.input = this.availableInputs[input_type]['1.0'];

    this.tapFG.patchValue({ input_type });

    // input type config model
    const { config: inputConfig, filter: filterConfig } = this.input;
    // if editing, some values might not be overrideable any longer, all should be prefilled in form
    const {
      config: agentConfig,
      filter: agentFilter,
    } = this.agentPolicy.policy.input;
    // tap config values, cannot be overridden if set
    const {
      config_predefined: preConfig,
      filter_predefined: preFilter,
    } = this.tap;

    // populate form controls for config
    const inputConfDynamicCtrl = Object.entries(inputConfig).reduce(
      (acc, [key, input]) => {
        const value = agentConfig?.[key] || '';
        if (!preConfig?.includes(key)) {
          acc[key] = [
            value,
            [
              !!input?.props?.required && input.props.required === true
                ? Validators.required
                : Validators.nullValidator,
            ],
          ];
        }
        return acc;
      },
      {},
    );

    this.inputConfigFG =
      Object.keys(inputConfDynamicCtrl).length > 0
        ? this._formBuilder.group(inputConfDynamicCtrl)
        : null;

    const inputFilterDynamicCtrl = Object.entries(filterConfig).reduce(
      (acc, [key, input]) => {
        const value = !!agentFilter?.[key] ? agentFilter[key] : '';
        // const disabled = !!preConfig?.[key];
        if (!preFilter?.includes(key)) {
          acc[key] = [
            value,
            [
              !!input?.props?.required && input.props.required === true
                ? Validators.required
                : Validators.nullValidator,
            ],
          ];
        }
        return acc;
      },
      {},
    );

    this.inputFilterFG =
      Object.keys(inputFilterDynamicCtrl).length > 0
        ? this._formBuilder.group(inputFilterDynamicCtrl)
        : null;
  }

  addHandler() {
    this.dialogService
      .open(HandlerPolicyAddComponent, {
        context: {
          backend: this.backend,
          modules: this.modules,
        },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((handler) => {
        // save handler to the policy being created/edited
        if (handler) {
          this.onHandlerAdded(handler);
        }
      });
  }

  onHandlerAdded(handler) {
    const { config, filter, type, name } = handler;

    this.modules[name] = {
      type,
      config,
      filter,
    };
  }

  onHandlerRemoved(name) {
    delete this.modules[name];
  }

  hasModules() {
    return Object.keys(this.modules).length > 0;
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  viewPolicy(id) {
    this.router.navigateByUrl(`/pages/datasets/policies/view/${id}`);
  }

  onYAMLSubmit() {
    const payload = {
      name: this.detailsFG.controls.name.value,
      description: this.detailsFG.controls.description.value,
      backend: this.detailsFG.controls.backend.value,
      format: this.format,
      policy_data: this.code,
      version:
        (!!this.isEdit &&
          !!this.agentPolicy.version &&
          this.agentPolicy.version) ||
        1,
    };

    this.submit(payload);
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFG.controls.name.value,
      description: this.detailsFG.controls.description.value,
      backend: this.detailsFG.controls.backend.value,
      tags: {},
      version:
        (!!this.isEdit &&
          !!this.agentPolicy.version &&
          this.agentPolicy.version) ||
        1,
      policy: {
        kind: 'collection',
        input: {
          tap: this.tap.name,
          input_type: this.tapFG.controls.input_type.value,
          ...Object.entries(this.inputConfigFG.controls)
            .map(([key, control]) => ({ [key]: control.value }))
            .reduce(
              (acc, curr) => {
                for (const [key, value] of Object.entries(curr)) {
                  if (!!value && value !== '') acc.config[key] = value;
                }
                return acc;
              },
              { config: {} },
            ),
          ...Object.entries(this.inputFilterFG.controls)
            .map(([key, control]) => ({ [key]: control.value }))
            .reduce(
              (acc, curr) => {
                for (const [key, value] of Object.entries(curr)) {
                  if (!!value && value !== '') acc.filter[key] = value;
                }
                return acc;
              },
              { filter: {} },
            ),
        },
        handlers: {
          modules: Object.entries(this.modules).reduce((acc, [key, value]) => {
            const { type, config, filter } = value;
            acc[key] = {
              type,
              config,
              filter,
            };
            if (Object.keys(config || {}).length > 0) acc[key][config] = config;
            return acc;
          }, {}),
        },
      },
    } as AgentPolicy;

    if (Object.keys(payload.policy?.input?.config).length <= 0)
      delete payload.policy.input.config;
    if (Object.keys(payload.policy?.input?.filter).length <= 0)
      delete payload.policy.input.filter;

    this.submit(payload);
  }

  submit(payload) {
    if (this.isEdit) {
      // updating existing sink
      this.agentPoliciesService
        .editAgentPolicy({ ...payload, id: this.agentPolicyID })
        .subscribe(() => {
          this.notificationsService.success(
            'Agent Policy successfully updated',
            '',
          );
          this.viewPolicy(this.agentPolicyID);
        });
    } else {
      this.agentPoliciesService.addAgentPolicy(payload).subscribe((next) => {
        this.notificationsService.success(
          'Agent Policy successfully created',
          '',
        );
        this.viewPolicy(next.id);
      });
    }
  }
}
