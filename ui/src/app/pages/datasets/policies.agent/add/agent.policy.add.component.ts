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
    this.agentPolicy = this.router.getCurrentNavigation().extras.state?.agentPolicy as AgentPolicy || null;
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.route.snapshot.paramMap.get('agentPolicy') as AgentPolicy;

    this.isEdit = !!this.agentPolicyID;
    this.isLoading[CONFIG.AGENT_POLICY] = this.isEdit;
    !!this.agentPolicyID && agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(resp => {
      this.agentPolicy = resp;
      this.isLoading[CONFIG.AGENT_POLICY] = false;
      this.readyForms();
    });

    this.readyForms();
  }

  readyForms() {
    const { name, description, backend } = this.agentPolicy
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

    this.handlers = Object.entries(this.agentPolicy.policy.handlers.modules)
      .map(([key, handler]) => ({ name: key, type: handler.config.type, ...handler }));

    this.detailsFG = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      description: [description],
      backend: [{ value: backend, disabled: backend !== '' }, [Validators.required]],
    });
    this.tapFG = this._formBuilder.group({
      'selected_tap': ['', Validators.required],
      'input_type': ['', Validators.required],
    });

    this.handlerSelectorFG = this._formBuilder.group({ 'selected_handler': [''] });
    this.dynamicHandlerConfigFG = this._formBuilder.group({});

    this.getBackendsList();
  }

  getBackendsList() {
    this.isLoading[CONFIG.BACKEND] = true;
    this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
      this.availableBackends = !!backends['data'] && backends['data'].reduce((acc, curr) => {
        acc[curr.backend] = curr;
        return acc;
      }, {});

      if (this.isLoading[CONFIG.AGENT_POLICY] === false) {
        this.onBackendSelected(this.agentPolicy.backend);
      }

      this.isLoading[CONFIG.BACKEND] = false;
    });
  }

  onBackendSelected(selectedBackend) {
    this.backend = this.availableBackends[selectedBackend];
    this.backend.config = {};

    // todo hardcoded for pktvisor
    this.getBackendData();
  }

  getBackendData() {
    Promise.all([this.getTaps(), this.getInputs(), this.getHandlers()])
      .then(value => {
        if (this.isEdit && this.agentPolicy) {
          const selected_tap = this.agentPolicy.policy.input.tap;
          this.tapFG.patchValue({ selected_tap });
          this.tapFG.controls.selected_tap.disable();
          this.onTapSelected(selected_tap);
          this.handlers = Object.entries(this.agentPolicy.policy.handlers.modules)
            .map(([key, handler]) => ({ ...handler, name: key, type: handler.config.type }));
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
          this.availableTaps = !!taps['data'] && taps['data'].reduce((acc, curr) => {
            acc[curr.name] = curr;
            return acc;
          }, {});

          this.isLoading[CONFIG.TAPS] = false;

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

  getInputs() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.INPUTS] = true;
      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'inputs'])
        .subscribe(inputs => {
          this.availableInputs = !!inputs['data'] && inputs['data'];

          this.isLoading[CONFIG.INPUTS] = false;

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

  getHandlers() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.HANDLERS] = true;

      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
        .subscribe(handlers => {
          this.availableHandlers = !!handlers['data'] && handlers['data'];

          this.handlerSelectorFG = this._formBuilder.group({
            'selected_handler': ['', [Validators.required]],
            'label': ['', [Validators.required]],
          });

          this.isLoading[CONFIG.HANDLERS] = false;
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

    const { config, filter } = !!this.liveHandler ? this.liveHandler : { config: {}, filter: {} };

    const dynamicControls = Object.entries(config).reduce((controls, [key, value]) => {
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

  onHandlerRemoved(selectedHandler) {
    delete this.handlers[selectedHandler];
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
          filter: Object.keys(this.inputFilterFG.controls)
            .map(key => ({ [key]: this.inputConfigFG.controls[key].value }))
            .reduce((acc, curr) => {
              for (const [key, value] of Object.entries(curr)) {
                if (!!value && value !== '') acc[key] = value;
              }
              return acc;
            }, {}),
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
