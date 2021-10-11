import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';

@Component({
  selector: 'ngx-agent-policy-add-component',
  templateUrl: './agent.policy.add.component.html',
  styleUrls: ['./agent.policy.add.component.scss'],
})
export class AgentPolicyAddComponent {
  // #forms
  // agent policy general information - name, desc, backend
  detailsFormGroup: FormGroup;

  // selected tap, input_type
  tapFormGroup: FormGroup;

  // dynamic input config
  inputFormGroup: FormGroup;

  // handlers
  handlerSelectorFormGroup: FormGroup;

  dynamicHandlerConfigFormGroup: FormGroup;

  // #key inputs holders
  // selected backend object
  backend: { [propName: string]: any };

  // selected tap object
  tap: { name: string, input_type: string, config_predefined: string[], agents: { total: number } };

  // selected input object
  input: { [propName: string]: any };

  // holds selected handler conf.
  // handler template currently selected, to be edited by user and then added to the handlers list or discarded
  liveHandler: { [propName: string]: any };

  // holds all handlers added by user
  handlers: { name: string, handler: { [propName: string]: any } }[] = [];

  // #services responses
  // hold info retrieved
  availableBackends: { [propName: string]: any }[];

  availableTaps: { name: string, input_type: string, config_predefined: string[], agents: { total: 1 } }[];

  availableInputs: { [propName: string]: { version: string, config: any } };

  availableHandlers: { [propName: string]: { version: string, config: any } };

  // #if edit
  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  isEdit: boolean;

  // #load controls
  isLoading = { 'taps': false, 'backend': false, 'inputs': false, 'handlers': false };

  agentPolicyLoading = false;

  constructor(
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentPolicy = this.router.getCurrentNavigation().extras.state?.agentPolicy as AgentPolicy || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.agentPolicyID;
    this.agentPolicyLoading = this.isEdit;

    !!this.agentPolicyID && agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(resp => {
      this.agentPolicy = resp;
      this.agentPolicyLoading = false;
      this.readyForms();
      this.getBackendsList();
    });

    this.readyForms();
    if (!this.isEdit) this.getBackendsList();
  }

  readyForms() {
    const { name, description, backend } = this.agentPolicy = !!this.agentPolicy ? this.agentPolicy : {
      name: '',
      description: '',
      backend: 'pktvisor',
      policy: {
        input: {
          tap: '',
          input_type: '',
          config_predefined: {},
        },
        handlers: { modules: {} },
      },
    };

    // TODO uncomment this line - BE currently not ssaving or not returning handlers
    // this.handlers = Object.entries(policy.handlers.modules).map(([key, value]) => ({name: key, handler: value}));
    this.handlers = [];

    this.detailsFormGroup = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_:][a-zA-Z0-9_]*$')]],
      description: [description],
      backend: [{ value: backend }, Validators.required],
    });

    this.handlerSelectorFormGroup = this._formBuilder.group({
      'selected_handler': ['', [Validators.required]],
      'label': ['', [Validators.required]],
    });
    this.dynamicHandlerConfigFormGroup = this._formBuilder.group({});
  }

  getBackendsList() {
    this.isLoading['backend'] = true;
    this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
      this.availableBackends = !!backends['data'] && backends['data'] || [];

      if (this.isEdit) {
        this.detailsFormGroup.controls.backend.disable();
        this.onBackendSelected(this.agentPolicy.backend);
      }

      this.isLoading['backend'] = false;
    });
  }

  onBackendSelected(selectedBackend) {
    this.backend = this.availableBackends[selectedBackend] || { backend: 'pktvisor' };
    this.backend.config = {};

    // todo hardcoded for pktvisor
    this.getTaps();
    this.getInputs();
    this.getHandlers();

  }

  getTaps() {
    this.isLoading['taps'] = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'taps'])
      .subscribe(taps => {
        this.availableTaps = !!taps['data'] && taps['data'] || [];

        const { input } = this.agentPolicy.policy;
        const { tap, input_type } = input;

        this.tapFormGroup = this._formBuilder.group({
          'selected_tap': [tap, Validators.required],
          'input_type': [input_type, Validators.required],
        });

        this.isLoading['taps'] = false;
      });
  }

  getInputs() {
    this.isLoading['inputs'] = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'inputs'])
      .subscribe(inputs => {
        this.availableInputs = !!inputs['data'] && inputs['data'] || {};

        this.isLoading['inputs'] = false;

        const input_type = !!this.tap ? this.tap : null;

        if (input_type) {
          this.onInputSelected(input_type);
          this.tapFormGroup.controls.input_type.disable();
        } else {
          this.input = null;
          this.tapFormGroup.controls.input_type.enable();
          this.tapFormGroup.controls.input_type.reset('');
        }
      });
  }

  onTapSelected(selectedTap) {
    this.tap = this.availableTaps[selectedTap];

    if (!this.tap?.config_predefined) this.tap['config_predefined'] = [];

    const { input_type } = this.tap;

    if (input_type && !!this.availableInputs) {
      this.onInputSelected(input_type);
      this.tapFormGroup.controls.input_type.disable();
    } else {
      this.input = null;
      this.tapFormGroup.controls.input_type.enable();
      this.tapFormGroup.controls.input_type.reset('');
    }
  }

  onInputSelected(selectedInput) {
    this.input = this.availableInputs[selectedInput];

    this.tapFormGroup.controls.input_type.setValue(selectedInput);

    // input type config model
    const { config: inputConfig } = this.input;
    // if editing, some values might not be overrideable any longer, all should be prefilled in form
    const agentConfig = !!this.isEdit && this.agentPolicy.policy?.config || null;
    // tap config values, cannot be overridden if set -- scratch that
    // only fields in config_predefined should be available for editing;
    const preConfig = this.tap.config_predefined;
    // assemble config obj with a three way merge of sorts
    // TODO this is under revision
    const finalConfig = { ...agentConfig, ...preConfig.reduce((acc, curr) => {
      acc[curr] = '';
      return acc;
    }, {}) };

    // populate form controls
    const dynamicFormControls = Object.keys(finalConfig)
      .reduce((acc, key) => {
        const value = finalConfig?.[key] || '';
        // const disabled = !!preConfig?.[key];
        const disabled = false;
        acc[key] = [
          { value, disabled },
          !!inputConfig?.[key]?.required ? Validators.required : null,
        ];
        return acc;
      }, {});

    this.inputFormGroup = this._formBuilder.group(dynamicFormControls);
  }

  getHandlers() {
    this.isLoading['handlers'] = true;

    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
      .subscribe(handlers => {
        this.availableHandlers = !!handlers['data'] && handlers['data'] || {};

        this.isLoading['handlers'] = false;
      });
  }


  onHandlerSelected(selectedHandler) {
    const { config } = this.availableHandlers[selectedHandler];

    const dynamicControls = Object.keys(config).reduce((acc, key) => {
      const field = config[key];
      acc[field.name] = [
        '',
        field.required ? Validators.required : null,
      ];
      return acc;
    }, {});

    this.handlerSelectorFormGroup.controls.label.setValue('');

    this.dynamicHandlerConfigFormGroup = this._formBuilder.group(dynamicControls);

    this.liveHandler = this.availableHandlers[selectedHandler];
  }

  onHandlerAdded() {
    const handlerName = this.handlerSelectorFormGroup.controls.label.value;
    this.handlers.push({
      name: handlerName,
      handler: {
        type: this.handlerSelectorFormGroup.controls.selected_handler.value,
        config: Object.keys(this.dynamicHandlerConfigFormGroup.controls)
          .reduce((acc, control) => {
            acc[control] = this.dynamicHandlerConfigFormGroup.controls[control].value;
            return acc;
            }, {}),
      },
    });
  }

  onHandlerRemoved(selectedHandler) {
    this.handlers.splice(selectedHandler, 1);
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFormGroup.controls.name.value,
      description: this.detailsFormGroup.controls.description.value,
      backend: this.availableBackends[this.detailsFormGroup.controls.backend.value].backend,
      tags: {},
      policy: {
        kind: 'collection',
        input: {
          tap: this.availableTaps[this.tapFormGroup.controls.selected_tap.value].name,
          input_type: this.tapFormGroup.controls.input_type.value,
          config: Object.keys(this.inputFormGroup.controls)
            .map(key => ({ [key]: this.inputFormGroup.controls[key].value }))
            .reduce((acc, curr) => {
              for (const [key, value] of Object.entries(curr)) {
                if (!!value && value !== '') acc[key] = value;
              }
              return acc;
            }, {}),
        },
        handlers: {
          modules: this.handlers.reduce((prev, handler) => {
            prev[handler.name] = {
              type: handler.handler?.['type'] || '',
              config: handler.handler?.['config'] || {},
            };
            return prev;
          }, {}),
        },
      },
      window_config: {
        num_periods: 5,
        deep_sample_rate: 100,
      },
    };

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
