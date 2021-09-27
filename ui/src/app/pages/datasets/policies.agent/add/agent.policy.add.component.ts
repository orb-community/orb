import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';

@Component({
  selector: 'ngx-agent-policy-add-component',
  templateUrl: './agent.policy.add.component.html',
  styleUrls: ['./agent.policy.add.component.scss'],
})
export class AgentPolicyAddComponent {
  strings = STRINGS;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  thirdFormGroup: FormGroup;

  customSinkSettings: {};

  selectedSinkSetting: any[];

  sink: Sink;

  sinkID: string;

  sinkTypesList = [];

  isEdit: boolean;
  isLoading = false;
  sinkLoading = false;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.sink = this.router.getCurrentNavigation().extras.state?.sink as Sink || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.sinkID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.sinkID;
    this.sinkLoading = this.isEdit;

    !!this.sinkID && sinksService.getSinkById(this.sinkID).subscribe(resp => {
      this.sink = resp;
      this.sinkLoading = false;
      this.getSinkBackends();
    });
    !this.sinkLoading && this.getSinkBackends();
  }

  getSinkBackends() {
    this.isLoading = true;
    this.sinksService.getSinkBackends().subscribe(backends => {
      this.sinkTypesList = backends.map(entry => entry.backend);
      this.customSinkSettings = this.sinkTypesList.reduce((accumulator, curr) => {
        const index = backends.findIndex(entry => entry.backend === curr);
        accumulator[curr] = backends[index].config.map(entry => ({
          type: entry.type,
          label: entry.title,
          prop: entry.name,
          input: entry.input,
          required: entry.required,
        }));
        return accumulator;
      }, {});
      const {name, description, backend, tags} = !!this.sink ? this.sink : {
        name: '',
        description: '',
        backend: 'prometheus', // default sink
        tags: {},
      } as Sink;
      this.firstFormGroup = this._formBuilder.group({
        name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
        description: [description],
        backend: [backend, Validators.required],
      });

      this.isEdit && this.firstFormGroup.controls.backend.disable();

      // builds secondFormGroup
      this.onSinkTypeSelected(backend);

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
    this.router.navigateByUrl('/pages/sinks');
  }

  onFormSubmit() {
    const payload = {
      name: this.firstFormGroup.controls.name.value,
      backend: this.firstFormGroup.controls.backend.value,
      description: this.firstFormGroup.controls.description.value,
      config: this.selectedSinkSetting.reduce((accumulator, current) => {
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
      this.sinksService.editSink({...payload, id: this.sinkID}).subscribe(() => {
        this.notificationsService.success('Sink successfully updated', '');
        this.goBack();
      });
    } else {
      this.sinksService.addSink(payload).subscribe(() => {
        this.notificationsService.success('Sink successfully created', '');
        this.goBack();
      });
    }

  }

  onSinkTypeSelected(selectedValue) {
    // SinkConfig<string> being the generic of all other `sinkTypes`.
    const conf = !!this.sink &&
      this.isEdit &&
      (selectedValue === this.sink.backend) &&
      this.sink?.config &&
      this.sink.config as SinkConfig<string> || null;

    this.selectedSinkSetting = this.customSinkSettings[selectedValue];

    const dynamicFormControls = this.selectedSinkSetting.reduce((accumulator, curr) => {
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
