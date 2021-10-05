import { Component, OnInit } from '@angular/core';

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
export class AgentPolicyAddComponent implements OnInit {
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
  tap: { [propName: string]: any };

  // selected input object
  input: { [propName: string]: any };

  // holds selected handler conf.
  // handler template currently selected, to be edited by user and then added to the handlers list or discarded
  liveHandler: { [propName: string]: any };

  // holds all handlers added by user
  handlers: { [propName: string]: any }[] = [];

  // #services responses
  // hold info retrieved
  availableBackends: { [propName: string]: any }[];

  availableTaps: { [propName: string]: any }[];

  availableInputs: { [propName: string]: any };

  availableHandlers: { [propName: string]: any };

  // #if edit
  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  isEdit: boolean;

  // #load controls
  isLoading = false;

  agentPolicyLoading = false;

  constructor(
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentPolicy = this.router.getCurrentNavigation().extras.state?.agentPolicy as AgentPolicy || {
      name: '',
      description: '',
      tags: {},
      backend: 'pktvisor',
    };
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.route.snapshot.paramMap.get('agentPolicy') as AgentPolicy;

    this.isEdit = !!this.agentPolicyID;
    this.agentPolicyLoading = this.isEdit;

    !!this.agentPolicyID && agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(resp => {
      this.agentPolicy = resp;
      this.agentPolicyLoading = false;
    });

    this.getBackendsList();
  }

  ngOnInit() {
    const { name, description, backend } = this.agentPolicy || { name: '', description: '', backend: '' };

    this.detailsFormGroup = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_:][a-zA-Z0-9_]*$')]],
      description: [description],
      backend: [backend, Validators.required],
    });
    this.tapFormGroup = this._formBuilder.group({
      'selected_tap': ['', Validators.required],
      'input_type': ['', Validators.required],
    });
    this.handlerSelectorFormGroup = this._formBuilder.group({ 'selected_handler': [''] });
  }

  getBackendsList() {
    this.isLoading = true;
    this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
      this.availableBackends = !!backends['data'] && backends['data'] || [];

      if (this.isEdit) {
        this.detailsFormGroup.controls.backend.disable();
        this.onBackendSelected(this.agentPolicy.backend);
      }

      this.isLoading = false;
    });
  }

  onBackendSelected(selectedBackend) {
    this.backend = this.availableBackends[selectedBackend];
    this.backend.config = {};

    // todo hardcoded for pktvisor
    this.getTaps();
    this.getInputs();
    this.getHandlers();

  }

  getTaps() {
    this.isLoading = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'taps'])
      .subscribe(taps => {
        this.availableTaps = !!taps['data'] && taps['data'] || [];

        this.isLoading = false;
      });
  }

  getInputs() {
    this.isLoading = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'inputs'])
      .subscribe(inputs => {
        this.availableInputs = !!inputs['data'] && inputs['data'] || {};

        this.isLoading = false;
      });
  }

  onTapSelected(selectedTap) {
    this.tap = this.availableTaps[selectedTap];
    const { input_type } = this.tap;
    if (input_type) {
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
    const { config } = this.input;
    const dynamicFormControls = Object.keys(config || {})
      .reduce((acc, key) => {
        acc[key] = [
          // TODO predef conf below can come from editing existing policy or from tap ->pre-conf<-(!has!precedence!)
          '', // or predefined conf when editing.
          config[key].required ? Validators.required : null,
        ];
        return acc;
      }, {});

    this.inputFormGroup = this._formBuilder.group(dynamicFormControls);

    // reconfig dynamic forms based on backend selected
    // this.backendConfigForms = Object.keys(this.backend.config)
    //   .reduce((formGroups, groupName, groupIndex) => {
    //     formGroups[groupName] = this._formBuilder.group({ [groupName]: ['', Validators.required] });
    //     return formGroups;
    //   }, {});

  }

  getHandlers() {
    this.isLoading = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
      .subscribe(handlers => {
        this.availableHandlers = !!handlers['data'] && handlers['data'] || {};

        this.handlerSelectorFormGroup = this._formBuilder.group({
          'selected_handler': ['', []],
        });

        this.isLoading = false;
      });
  }


  onHandlerSelected(selectedHandler) {
    this.liveHandler = this.availableHandlers[selectedHandler];
    const {config} = this.liveHandler;
    const dynamicControls = Object.keys(config).reduce((acc, key) => {
      acc[key] = [
        '',
        config[key].required ? Validators.required : null,
      ];
      return acc;
    }, {});

    this.dynamicHandlerConfigFormGroup = this._formBuilder.group(dynamicControls);
  }

  onHandlerAdded() {
    const { label: { value } } = this.dynamicHandlerConfigFormGroup.controls;
    this.handlers[value] = this.dynamicHandlerConfigFormGroup.value;
  }

  onHandlerRemoved(selectedHandler) {
    delete this.handlers[selectedHandler];
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFormGroup.controls.name.value,
      description: this.detailsFormGroup.controls.description.value,
      backend: this.detailsFormGroup.controls.backend.value,
      // config: this.selectedTap.reduce((accumulator, current) => {
      //   accumulator[current.prop] = this.tapFormGroup.controls[current.prop].value;
      //   return accumulator;
      // }, {}),
      validate_only: false, // Apparently this guy is required..
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
