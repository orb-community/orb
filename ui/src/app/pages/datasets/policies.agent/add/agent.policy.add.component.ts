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
  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  thirdFormGroup: FormGroup;

  availableTaps: [];

  selectedTap: any;

  availableInputs: [];

  selectedInput: any;

  availableHandlers: [];

  selectedHandlers: [];

  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  tapList = [];

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
    this.agentPolicy = this.router.getCurrentNavigation().extras.state?.agentPolicy as AgentPolicy || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.agentPolicyID;
    this.agentPolicyLoading = this.isEdit;

    !!this.agentPolicyID && agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(resp => {
      this.agentPolicy = resp;
      this.agentPolicyLoading = false;
      this.getTapsList();
    });
    !this.agentPolicyLoading && this.getTapsList();
  }

  getTapsList() {
    this.isLoading = true;
    this.agentPoliciesService.getPktVisorTaps().subscribe(taps => {
      this.tapList = taps.map(entry => entry.backend);
      this.availableTaps = this.tapList.reduce((accumulator, curr) => {
        const index = taps.findIndex(entry => entry.backend === curr);
        accumulator[curr] = taps[index].config.map(entry => ({
          type: entry.type,
          label: entry.title,
          prop: entry.name,
          input: entry.input,
          required: entry.required,
        }));
        return accumulator;
      }, {});
      const {name, description, backend, tags} = !!this.agentPolicy ? this.agentPolicy : {
        name: '',
        description: '',
        backend: 'dns', // default sink
        tags: {},
      } as AgentPolicy;
      this.firstFormGroup = this._formBuilder.group({
        name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_:][a-zA-Z0-9_]*$')]],
        description: [description],
        backend: [backend, Validators.required],
      });

      this.isEdit && this.firstFormGroup.controls.backend.disable();

      // builds secondFormGroup
      this.onTapTypeSelected(backend);

      this.thirdFormGroup = this._formBuilder.group({
        tags: [Object.keys(tags || {}).map(key => ({[key]: tags[key]})),
          Validators.minLength(1)],
        key: [''],
        value: [''],
      });

      this.isLoading = false;
    });
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  onFormSubmit() {
    const payload = {
      name: this.firstFormGroup.controls.name.value,
      backend: this.firstFormGroup.controls.backend.value,
      description: this.firstFormGroup.controls.description.value,
      config: this.selectedTap.reduce((accumulator, current) => {
        accumulator[current.prop] = this.secondFormGroup.controls[current.prop].value;
        return accumulator;
      }, {}),
      tags: this.thirdFormGroup.controls.tags.value.reduce((prev, curr) => {
        for (const [key, value] of Object.entries(curr)) {
          prev[key] = value;
        }
        return prev;
      }, {}),
      validate_only: false, // Apparently this guy is required..
    };
    // TODO Check this out
    // console.log(payload);
    if (this.isEdit) {
      // updating existing sink
      this.agentPoliciesService.editAgentPolicy({...payload, id: this.agentPolicyID}).subscribe(() => {
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

  onTapTypeSelected(selectedValue) {
    // SinkConfig<string> being the generic of all other `sinkTypes`.
    const conf = !!this.agentPolicy &&
      this.isEdit &&
      (selectedValue === this.agentPolicy.backend) &&
      this.agentPolicy?.policy &&
      this.agentPolicy.policy as TapConfig || null;

    this.selectedTap = this.availableTaps[selectedValue];

    const dynamicFormControls = this.selectedTap.reduce((accumulator, curr) => {
      accumulator[curr.prop] = [
        !!conf && (curr.prop in conf) && conf[curr.prop] ||
        '',
        curr.required ? Validators.required : null,
      ];
      return accumulator;
    }, {});

    this.secondFormGroup = this._formBuilder.group(dynamicFormControls);
  }

  // addTag button should be [disabled] = `$sf.controls.key.value !== ''`
  onAddTag() {
    const {tags, key, value} = this.thirdFormGroup.controls;
    // sanitize minimally anyway
    if (key?.value && key.value !== '') {
      if (value?.value && value.value !== '') {
        // key and value fields
        tags.reset([{[key.value]: value.value}].concat(tags.value));
        key.reset('');
        value.reset('');
      }
    } else {
      // TODO remove this else clause and error
      console.error('This shouldn\'t be happening');
    }
  }

  onRemoveTag(tag: any) {
    const {tags, tags: {value: tagsList}} = this.thirdFormGroup.controls;
    const indexToRemove = tagsList.indexOf(tag);

    if (indexToRemove >= 0) {
      tags.setValue(tagsList.slice(0, indexToRemove).concat(tagsList.slice(indexToRemove + 1)));
    }
  }
}
