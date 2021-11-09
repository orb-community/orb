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
  // todo rename
  onebigform: FormGroup;
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
  tap: PolicyTap = { name: '' };

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
  handlers: {
    name: string,
    type: string,
    config: { [propName: string]: {} | any },
  }[] = [];

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
    },
  };

  // #if edit
  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  isEdit: boolean;

  // #load controls
  loadControls = Object.entries(CONFIG)
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
    this.newAgentPolicy();
    this.readyForms();

    this.agentPolicy = this.router.getCurrentNavigation().extras.state?.agentPolicy as AgentPolicy || null;
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.route.snapshot.paramMap.get('agentPolicy') as AgentPolicy;

    this.isEdit = !!this.agentPolicyID;
    this.loadControls[CONFIG.AGENT_POLICY] = this.isEdit;
    !!this.agentPolicyID && agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(resp => {
      this.agentPolicy = resp;
      this.loadControls[CONFIG.AGENT_POLICY] = false;
    });

    this.getBackendsList().then((backends) => {

    }).catch(reason => console.warn(`Couldn't retrieve available backends. Reason: ${ reason }`));
  }

  readyForms() {
    // todo this is pktvisor specific
    this.onebigform = this._formBuilder.group({
      name: [null, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      description: [null],
      backend: [null, [Validators.required]],
      policy: this._formBuilder.group({
        input: this._formBuilder.group({
          input_type: [null, [Validators.required]],
          tap: [null, [Validators.required]],
          config: this._formBuilder.group({}),
        }, [Validators.required]),
        handlers: this._formBuilder.group({
          modules: this._formBuilder.group({}),
        }),
      }),
    });

    this.handlerSelectorFG = this._formBuilder.group({ 'selected_handler': [''] });
    this.dynamicHandlerConfigFG = this._formBuilder.group({});
  }

  updateForms() {
    const {
      name,
      description,
      backend,
      policy: {
        input: {
          tap,
          input_type,
          config,
        },
        handlers: {
          modules,
        },
      },
    } = this.agentPolicy
      = {
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
      ...this.agentPolicy,
    } as AgentPolicy;

    this.onebigform.patchValue({
      name,
      description,
      backend,
      policy: {
        input: {
          tap,
          input_type,
          config,
        },
        handlers: {
          modules,
        },
      },
    }, {emitEvent: true});
  }

  retrieveEditAgentPolicy() {
    return new Promise((resolve) => {
      this.loadControls[CONFIG.AGENT_POLICY] = true;
      this.agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(agentPolicy => {
        this.loadControls[CONFIG.AGENT_POLICY] = false;
        resolve(agentPolicy as AgentPolicy);
      });
    });
  }

  newAgentPolicy() {
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

  getBackendsList() {
    return new Promise((resolve) => {
      this.loadControls[CONFIG.BACKEND] = true;
      this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
        this.availableBackends = !!backends && backends.reduce((acc, curr) => {
          acc[curr.backend] = curr;
          return acc;
        }, {});

        this.loadControls[CONFIG.BACKEND] = false;
        resolve(this.availableBackends);
      });
    });
  }

  onBackendSelected(selectedBackend) {
    this.backend['config'] = {};

    // todo hardcoded for pktvisor
    this.getBackendData(selectedBackend.backend);
  }

  isLoading() {
    return Object.values<boolean>(this.loadControls).reduce((prev, curr) => prev && curr);
  }

  getBackendData(backendName) {
    // TODO pktvisor specific
    Promise.all([this.getTaps(backendName), this.getInputs(backendName), this.getHandlers(backendName)])
      .then(value => {
        if (this.isEdit && this.agentPolicy) {
          const selected_tap = this.agentPolicy.policy.input.tap;
          this.tapFG.patchValue({ selected_tap });
          this.tapFG.controls.selected_tap.disable();
          this.onTapSelected(selected_tap);
          this.handlers = Object.entries(this.agentPolicy.policy.handlers.modules)
            .map(([key, handler]) => ({ ...handler, name: key, type: handler.config.type }));
          this.updateForms();
        }
      }, reason => console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`))
      .catch(reason => {
        console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`);
      });
  }

  getTaps(backend) {
    return new Promise((resolve) => {
      this.loadControls[CONFIG.TAPS] = true;
      this.agentPoliciesService.getBackendConfig([backend, 'taps'])
        .subscribe(taps => {
          this.availableTaps = !!taps && taps.reduce((acc, curr) => {
            acc[curr.name] = curr;
            return acc;
          }, {});

          this.loadControls[CONFIG.TAPS] = false;

          resolve(this.availableTaps);
        });
    });
  }

  onTapSelected(selectedTap) {
    this.tap = this.availableTaps[selectedTap];

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

  getInputs(backend) {
    return new Promise((resolve) => {
      this.loadControls[CONFIG.INPUTS] = true;
      this.agentPoliciesService.getBackendConfig([backend, 'inputs'])
        .subscribe(inputs => {
          this.availableInputs = !!inputs && inputs;

          this.loadControls[CONFIG.INPUTS] = false;

          resolve(this.availableInputs);
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
    const agentConfig = !!this.isEdit ? this.agentPolicy.policy?.input?.config : null;
    // tap config values, cannot be overridden if set
    const preConfig = this.tap.config_predefined;

    if (this.isEdit === false) {
      this.agentPolicy.policy = { input: { config: {} } };
    }

    // populate form controls for config
    const inputConfDynamicCtrl = Object.entries(inputConfig)
      .reduce((acc, [key, input]) => {
        const value = !!agentConfig?.[key] ? agentConfig[key] : '';
        if (!preConfig.includes(key)) {
          acc[key] = [
            { value },
            [!!input?.props?.required && input.props.required === true ? Validators.required : Validators.nullValidator],
          ];
        }
        return acc;
      }, {});

    this.inputConfigFG = Object.keys(inputConfDynamicCtrl).length > 0 ? this._formBuilder.group(inputConfDynamicCtrl) : null;

    const inputFilterDynamicCtrl = Object.entries(filterConfig)
      .reduce((acc, [key, input]) => {
        const value = !!agentConfig?.[key] ? agentConfig[key] : '';
        const disabled = !!preConfig?.[key];
        acc[key] = [
          { value, disabled },
          [!!input?.props?.required && input.props.required === true ? Validators.required : Validators.nullValidator],
        ];

        return acc;
      }, {});

    this.inputFilterFG = Object.keys(inputFilterDynamicCtrl).length > 0 ? this._formBuilder.group(inputFilterDynamicCtrl) : null;

  }

  getHandlers(backend) {
    return new Promise((resolve) => {
      this.loadControls[CONFIG.HANDLERS] = true;

      this.agentPoliciesService.getBackendConfig([backend, 'handlers'])
        .subscribe(handlers => {
          this.availableHandlers = !!handlers && handlers;

          this.handlerSelectorFG = this._formBuilder.group({
            'selected_handler': ['', [Validators.required]],
            'label': ['', [Validators.required]],
          });

          this.loadControls[CONFIG.HANDLERS] = false;
          resolve(this.availableBackends);
        });
    });
  }

  onHandlerSelected(selectedHandler) {
    if (this.dynamicHandlerConfigFG) {
      this.dynamicHandlerConfigFG = null;
    }

    this.liveHandler = selectedHandler !== '' && !!this.availableHandlers[selectedHandler] ?
      { ...this.availableHandlers[selectedHandler], type: selectedHandler } : null;

    const { config } = !!this.liveHandler ? this.liveHandler : { config: {} };

    const dynamicControls = Object.entries(config || {}).reduce((controls, [key]) => {
      controls[key] = ['', [Validators.required]];
      return controls;
    }, {});

    this.dynamicHandlerConfigFG = Object.keys(dynamicControls).length > 0 ? this._formBuilder.group(dynamicControls) : null;
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
    this.handlers.push({
      name: handlerName,
      type: this.liveHandler.type,
      config,
    });

    this.handlerSelectorFG.reset({
      selected_handler: { value: '', disabled: false },
      label: { value: '', disabled: false },
    });
    this.onHandlerSelected('');
  }

  onHandlerRemoved(index) {
    this.handlers.splice(index, 1);
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
          tap: this.tapFG.controls.selected_tap.value,
          input_type: this.tapFG.controls.input_type.value,
          config: Object.keys(this.inputConfigFG.controls)
            .map(key => ({ [key]: this.inputConfigFG.controls[key].value }))
            .reduce((acc, curr) => {
              for (const [key, value] of Object.entries(curr)) {
                if (!!value && value !== '') acc[key] = value;
              }
              return acc;
            }, {}),
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
          modules: this.handlers.reduce((acc, handler) => {
            acc[handler.name] = {
              ...(Object.keys(handler.config).length > 0 ? { config: handler.config } : {}),
              type: handler.type,
            };
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
