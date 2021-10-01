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

  /**
   * Forms
   * //NOTE: refactor to be all dynamic
   */
    // agent policy general information
  detailsFormGroup: FormGroup;

  // Refactor while coding :)
  backendConfigForms: { [propName: string]: FormGroup };

  availableBackends: { [propName: string]: any }[];

  availableTaps: { [propName: string]: any }[];

  availableInputs: { [propName: string]: any };

  availableHandlers: { [propName: string]: any }[];

  backend: { [propName: string]: any };

  tap: { [propName: string]: any };

  input: { [propName: string]: any };

  handlers: { [propName: string]: any }[] = [];

  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  isEdit: boolean;

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

    const { name, description, backend } = this.agentPolicy || { name: '', description: '', backend: '' };

    this.backendConfigForms = {};

    this.detailsFormGroup = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_:][a-zA-Z0-9_]*$')]],
      description: [description, Validators.required],
      backend: [backend, Validators.required],
    });

    this.getBackendsList();
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
    // reconfig dynamic forms based on backend selected
    // this.backendConfigForms = Object.keys(this.backend.config)
    //   .reduce((formGroups, groupName, groupIndex) => {
    //     formGroups[groupName] = this._formBuilder.group({ [groupName]: ['', Validators.required] });
    //     return formGroups;
    //   }, {});

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

        this.backendConfigForms['taps'] = this._formBuilder.group({
          'tap': ['', [Validators.required]], // tap name
        });

        this.isLoading = false;
      });
  }

  getInputs() {
    this.isLoading = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'inputs'])
      .subscribe(inputs => {
        this.availableInputs = !!inputs['data'] && inputs['data'] || {};
        this.backend.config['inputs'] = this.availableInputs;

        this.isLoading = false;
      });
  }

  onTapSelected(selectedTap) {
    this.tap = this.availableTaps[selectedTap];
    const { taps } = this.backendConfigForms;
    Object.keys(this.tap.config).forEach(key => {
      taps.addControl(key, this._formBuilder.control('', [Validators.required]));
    });

    this.backendConfigForms['taps'].addControl(
      'input_type',
      this._formBuilder.control([this.tap.input_type, [Validators.required]]),
    );

    this.onInputSelected(this.tap.input_type);
    this.backendConfigForms['taps'].controls.input_type.disable();
  }

  onInputSelected(selectedInput) {
    this.input = this.availableInputs[selectedInput];
    const inputs = this._formBuilder.group({});
    Object.keys(this.input.config).forEach(key => {
      inputs.addControl(key, this._formBuilder.control('', [Validators.required]));
    });
    this.backendConfigForms['inputs'] = inputs;
  }

  getHandlers() {
    this.isLoading = true;
    this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
      .subscribe(handlers => {
        this.backend.config['handlers'] = !!handlers['data'] && handlers['data'] || {};

        this.backendConfigForms['handlers'] = this._formBuilder.group({'selected_handler': ['', []]});

        this.isLoading = false;
      });
  }


  onHandlerSelected(selectedHandler) {

  }

  onHandlerAdded() {

  }

  onHandlerRemoved(selectedHandler) {

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
