import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { DynamicFormConfig } from 'app/common/interfaces/orb/dynamic.form.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { PolicyTap } from 'app/common/interfaces/orb/policy/policy.tap.interface';

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
  // #forms
  // agent policy general information - name, desc, backend
  detailsFG: FormGroup;

  // selected tap, input_type
  tapFG: FormGroup;

  // dynamic input config
  inputConfigFG: FormGroup;

  // dynamic input filter config
  inputFilterFG: FormGroup;

  // handlers
  handlerSelectorFG: FormGroup;

  dynamicHandlerConfigFG: FormGroup;

  // #key inputs holders
  // selected backend object
  backend: { [propName: string]: any };

  // selected tap object
  tap: PolicyTap;

  // selected input object
  input: {
    version?: string,
    config?: DynamicFormConfig,
    filter?: DynamicFormConfig,
  };

  // holds selected handler conf.
  // handler template currently selected, to be edited by user and then added to the handlers list or discarded
  liveHandler: {
    version?: string,
    config?: DynamicFormConfig,
    filter?: DynamicFormConfig,
    type?: string,
  };

  // holds all handlers added by user
  modules: {
    [propName: string]: {
      name?: string,
      type?: string,
      config?: { [propName: string]: {} | any },
    },
  } = {};

  // #services responses
  // hold info retrieved
  availableBackends: {
    [propName: string]: {
      backend: string,
      description: string,
    },
  };

  availableTaps: { [propName: string]: PolicyTap };

  availableInputs: {
    [propName: string]: {
      version?: string,
      config?: DynamicFormConfig,
      filter?: DynamicFormConfig,
    },
  };

  availableHandlers: {
    [propName: string]: {
      version?: string,
      config?: DynamicFormConfig,
      filter?: DynamicFormConfig,
      metrics?: DynamicFormConfig,
      metrics_groups?: DynamicFormConfig,
    },
  } = {};

  // #if edit
  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  isEdit: boolean;

  // #load controls
  isLoading = Object.entries(CONFIG)
    .reduce((acc, [value]) => {
      acc[value] = false;
      return acc;
    }, {});

  constructor(
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.newAgent();
    this.isEdit = !!this.agentPolicyID;

    this.readyForms();

    Promise.all([
      this.isEdit ? this.retrieveAgentPolicy() : Promise.resolve(),
      this.getBackendsList(),
    ]).catch(reason => console.warn(`Couldn't fetch data. Reason: ${ reason }`))
      .then(() => this.updateForms())
      .catch((reason) => console.warn(`Couldn't fetch ${ this.agentPolicy?.backend } data. Reason: ${ reason }`));
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
    return new Promise(resolve => {
      this.agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(policy => {
        const {
          name,
          description,
          backend,
          policy: {
            input: {
              tap,
              input_type,
            },
            handlers: {
              modules,
            },
          },
        } = policy;
        this.agentPolicy = {
          ...this.agentPolicy,
          name,
          description,
          backend,
          policy: {
            input: {
              tap,
              input_type,
            },
            handlers: {
              modules,
            },
          },
        };
        this.isLoading[CONFIG.AGENT_POLICY] = false;
        resolve(policy);
      });
    });
  }

  readyForms() {
    const {
      name,
      description,
      backend,
      policy: {
        input: {
          tap,
          input_type,
        },
        handlers: {
          modules,
        },
      },
    } = this.agentPolicy;

    this.modules = modules;

    this.detailsFG = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      description: [description],
      backend: [{ value: backend, disabled: backend !== '' }, [Validators.required]],
    });
    this.tapFG = this._formBuilder.group({
      selected_tap: [tap, Validators.required],
      input_type: [input_type, Validators.required],
    });

    this.handlerSelectorFG = this._formBuilder.group({
      'selected_handler': ['', [Validators.required]],
      'label': ['', [Validators.required]],
    });

    this.dynamicHandlerConfigFG = this._formBuilder.group({});
  }

  updateForms() {
    const {
      name,
      description,
      backend,
      policy: {
        handlers: {
          modules,
        },
      },
    } = this.agentPolicy;

    this.detailsFG.patchValue({ name, description, backend });

    this.modules = modules;

    this.dynamicHandlerConfigFG = this._formBuilder.group({});

    this.onBackendSelected(backend).catch(reason => console.warn(`${ reason }`));


  }

  getBackendsList() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.BACKEND] = true;
      this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
        this.availableBackends = !!backends && backends.reduce((acc, curr) => {
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
        resolve();
      });
    });
  }

  getBackendData() {
    return Promise.all([this.getTaps(), this.getInputs(), this.getHandlers()])
      .then(value => {
        if (this.isEdit && this.agentPolicy) {
          const selected_tap = this.agentPolicy.policy.input.tap;
          this.tapFG.patchValue({ selected_tap }, { emitEvent: true });
          this.onTapSelected(selected_tap);
          this.tapFG.controls.selected_tap.disable();
        }

      }, reason => console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`))
      .catch(reason => {
        console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`);
      });
  }

  getTaps() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.TAPS] = true;
      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'taps'])
        .subscribe(taps => {
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
    const { input_type, config_predefined } = this.tap;

    this.tap['config'] = {
      ...config_predefined,
      ...input.config,
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
      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'inputs'])
        .subscribe(inputs => {
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
    const { config: agentConfig, filter: agentFilter } = !!this.isEdit ? this.agentPolicy.policy?.input : null;
    // tap config values, cannot be overridden if set
    const preConfig = this.tap.config_predefined;

    // populate form controls for config
    const inputConfDynamicCtrl = Object.entries(inputConfig)
      .reduce((acc, [key, input]) => {
        const value = agentConfig?.[key] || '';
        if (!preConfig.includes(key)) {
          acc[key] = [
            value,
            [!!input?.props?.required && input.props.required === true ? Validators.required : Validators.nullValidator],
          ];
        }
        return acc;
      }, {});

    this.inputConfigFG = Object.keys(inputConfDynamicCtrl).length > 0 ? this._formBuilder.group(inputConfDynamicCtrl) : null;

    const inputFilterDynamicCtrl = Object.entries(filterConfig)
      .reduce((acc, [key, input]) => {
        const value = !!agentConfig?.[key] ? agentConfig[key] : '';
        // const disabled = !!preConfig?.[key];
        if (!preConfig.includes(key)) {
          acc[key] = [
            value,
            [!!input?.props?.required && input.props.required === true ? Validators.required : Validators.nullValidator],
          ];
        }
        return acc;
      }, {});

    this.inputFilterFG = Object.keys(inputFilterDynamicCtrl).length > 0 ? this._formBuilder.group(inputFilterDynamicCtrl) : null;

  }

  getHandlers() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.HANDLERS] = true;

      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
        .subscribe(handlers => {
          this.availableHandlers = handlers || {};

          this.isLoading[CONFIG.HANDLERS] = false;
          resolve(handlers);
        });
    });
  }

  onHandlerSelected(selectedHandler) {
    if (this.dynamicHandlerConfigFG) {
      this.dynamicHandlerConfigFG = null;
    }

    this.liveHandler = selectedHandler !== '' && !!this.availableHandlers[selectedHandler] ?
      { ...this.availableHandlers[selectedHandler], type: selectedHandler } : null;

    const { config } = this.liveHandler || { config: {} };

    const dynamicControls = Object.entries(config || {}).reduce((controls, [key]) => {
      controls[key] = ['', [Validators.required]];
      return controls;
    }, {});

    this.dynamicHandlerConfigFG = Object.keys(dynamicControls).length > 0 ? this._formBuilder.group(dynamicControls) : null;
  }

  checkValidName() {
    const { policy: { handlers: { modules } } } = this.agentPolicy;
    const { value } = this.handlerSelectorFG.controls.label;
    return !(value === '' || Object.keys(modules || {}).find(name => value === name));
  }

  onHandlerAdded() {
    let config = {};

    if (this.dynamicHandlerConfigFG !== null) {
      config = Object.entries(this.dynamicHandlerConfigFG.controls)
        .reduce((acc, [key, control]) => {
          acc[key] = control.value;
          return acc;
        }, {});
    }

    const handlerName = this.handlerSelectorFG.controls.label.value;
    this.modules[handlerName] = ({
      type: this.liveHandler.type,
      config,
    });

    this.handlerSelectorFG.reset({
      selected_handler: { value: '', disabled: false },
      label: { value: '', disabled: false },
    });
    this.onHandlerSelected('');
  }

  onHandlerRemoved(handlerName) {
    delete this.modules[handlerName];
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFG.controls.name.value,
      description: this.detailsFG.controls.description.value,
      backend: this.detailsFG.controls.backend.value,
      tags: {},
      version: !!this.isEdit && !!this.agentPolicy.version && this.agentPolicy.version || 1,
      policy: {
        kind: 'collection',
        input: {
          tap: this.tap.name,
          input_type: this.tapFG.controls.input_type.value,
          ...Object.entries(this.inputConfigFG.controls)
            .map(([key, control]) => ({ [key]: control.value }))
            .reduce((acc, curr) => {
              for (const [key, value] of Object.entries(curr)) {
                if (!!value && value !== '') acc.config[key] = value;
              }
              return Object.keys(acc.config).length > 0 ? acc : null;
            }, {config: {}}),
          // filter: Object.keys(this.inputFilterFG.controls)
          //   .map(key => ({ [key]: this.inputConfigFG.controls[key].value }))
          //   .reduce((acc, curr) => {
          //     for (const [key, value] of Object.entries(curr)) {
          //       if (!!value && value !== '') acc[key] = value;
          //     }
          //     return acc;
          //   }, {}),
        },
        handlers: {
          modules: Object.entries(this.modules).reduce((acc, [key, value]) => {
            const {type, config} = value;
            acc[key] = {
              type: type,
            };
            if (Object.keys(config || {}).length > 0) acc[key][config] = config;
            return acc;
          }, {}),
        },
      },
    } as AgentPolicy;

    if (this.isEdit) {
      // updating existing sink
      this.agentPoliciesService.editAgentPolicy({ ...payload, id: this.agentPolicyID }).subscribe(() => {
        this.notificationsService.success('Agent Policy successfully updated', '');
        this.goBack();
      });
    } else {
      this.agentPoliciesService.addAgentPolicy(payload).subscribe(() => {
        this.notificationsService.success('Agent Policy successfully created', '');
        this.goBack();
      });
    }
  }
}
