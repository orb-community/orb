import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { TapConfig } from 'app/common/interfaces/orb/policy/config/tap.config.interface';

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

  // !!!
  // all form groups from here on should be dynamic, not
  // even declared like this
  // hey, I'm walking here
  // this is input/taps which come
  tapFormGroup: FormGroup;

  // with conf objects, hence tapConfigForm
  tapConfigFormGroup: FormGroup;

  // pktvisors configure handlers when iface is set
  handlersFormGroup: FormGroup;

  // !!!

  // Refactor while coding :)
  backendConfigForms: { [propName: string]: FormGroup };

  availableBackends: { [propName: string]: any };

  backend: any;

  availableTaps: { [propName: string]: any };

  selectedTap: any;

  tapConfig: { [propName: string]: any };

  availableHandlers: { [propName: string]: any };

  selectedHandlers: any;

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
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.agentPolicyID;
    this.agentPolicyLoading = this.isEdit;

    !!this.agentPolicyID && agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(resp => {
      this.agentPolicy = resp;
      this.agentPolicyLoading = false;
    });

    const { name, description, backend } = this.agentPolicy;

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
      this.availableBackends = { backends };

      this.isEdit && this.detailsFormGroup.controls.backend.disable();

      // builds secondFormGroup
      this.onBackendSelected(this.detailsFormGroup.controls.backend.value);

      this.isLoading = false;
    });
  }

  onBackendSelected(selectedBackend) {
    // const conf = !!this.agentPolicy &&
    //   this.isEdit &&
    //   (selectedBackend === this.agentPolicy.backend) &&
    //   this.agentPolicy?.policy &&
    //   this.agentPolicy.policy as TapConfig || null;

    this.backend = this.availableBackends[selectedBackend];
    this.getBackendConfig();
  }

  getBackendConfig() {
    // at this point I'm going to start hardcoding this whole api discovery I'm mocking
    // and at the same time creating the dynamic forms.
    // lets finish the forms and deliver.
    // it is very likely frontend will just get one single json with all
    // possible configs to build the forms instead of querying and doing the discovery itself
    // but it would go like
    // last argument being array inside array really
    // this.agentsPoliciesService.discover([this.agentsBackendUrl, selectedBackend.name, Object.keys(selectedBackend)])
    this.getTapsList();
  }

  getTapsList() {
    // this.isLoading = true;
    // this.agentPoliciesService.getPktVisorInputs().subscribe((taps) => {
    //   this.availableBackends['pktvisor']['taps'] = taps;
    //   this.backend['inputs'] = taps;
    const { taps } = this.availableBackends['pktvisor'];
    this.backendConfigForms['taps'] = this._formBuilder
      .group(Object.keys(taps['config'])
        .reduce((accumulator, curr) => {
          accumulator[curr] = Object.keys(taps['config'])
            .map(entry => ({
              type: taps['config'][entry].type,
              label: taps['config'][entry].title,
              prop: taps['config'][entry].name,
              input: taps['config'][entry].input,
              required: taps['config'][entry].required,
            }));
          return accumulator;
        }, {}));

    this.getHandlersList();
  }

  onTapSelected(selectedTap) {
    // const dynamicFormControls = this.selectedTap.reduce((accumulator, curr) => {
    // accumulator[curr.prop] = [
    //   !!conf && (curr.prop in conf) && conf[curr.prop] ||
    //   '',
    //   curr.required ? Validators.required : null,
    // ];
    // return accumulator;
    // }, {});

    // this.tapFormGroup = this._formBuilder.group(dynamicFormControls);
  }

  getHandlersList() {
    this.isLoading = true;
    this.agentPoliciesService.getPktVisorHandlers().subscribe((handlers) => {
      this.availableHandlers = handlers;
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
