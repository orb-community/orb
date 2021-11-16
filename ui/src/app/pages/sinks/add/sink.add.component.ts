import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';

@Component({
  selector: 'ngx-sink-add-component',
  templateUrl: './sink.add.component.html',
  styleUrls: ['./sink.add.component.scss'],
})
export class SinkAddComponent {
  strings = STRINGS;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  thirdFormGroup: FormGroup;

  customSinkSettings: {};

  selectedSinkSetting: any[];

  selectedTags: {[propName: string]: string};

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
    this.initializeForms();
    this.sinkID = this.route.snapshot.paramMap.get('id');
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean || !!this.sinkID;
    this.sinkLoading = this.isEdit;
    this.isLoading = true;

    Promise.all([this.getSink(), this.getSinkBackends()]).then(() => {
      // builds secondFormGroup
      const {backend} = this.sink;
      this.isLoading = false;
      if (backend !== '') this.onSinkTypeSelected(backend);
    }).catch(reason => console.warn(`Couldn't retrieve data. Reason: ${reason}`));
  }

  newSink() {
    return {
      name: '',
      description: '',
      backend: 'prometheus', // default sink
      tags: {},
    } as Sink;
  }

  initializeForms() {
    const { name, description, backend, tags } = this.sink = {
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

    this.selectedTags = {...tags};

    this.thirdFormGroup = this._formBuilder.group({
      key: [''],
      value: [''],
    });
  }

  getSink() {
    return new Promise(resolve => {
      if (this.sinkID) {
        this.sinksService.getSinkById(this.sinkID).subscribe(resp => {
          const {name, backend, description, tags} = this.sink = {...this.sink, ...resp};
          this.sinkLoading = false;
          this.selectedTags = tags;

          this.firstFormGroup.patchValue({name, description, backend});
          this.firstFormGroup.controls.backend.disable();
          this.firstFormGroup.controls.name.disable();
          resolve(resp);
        });
      } else {
        const sink = this.newSink();
        resolve(sink);
      }
    });
  }

  getSinkBackends() {
    return new Promise(resolve => {
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
        resolve(true);
      });
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
      tags: { ...this.selectedTags },
    };

    if (this.isEdit) {
      // updating existing sink
      this.sinksService.editSink({ ...payload, id: this.sinkID }).subscribe(() => {
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

  checkValidName() {
    const { tags } = this.sink;
    const { value } = this.thirdFormGroup?.controls?.label || {value: ''};
    return !(value === '' || Object.keys(tags || {}).find(name => value === name));
  }

  // addTag button should be [disabled] = `$sf.controls.key.value !== ''`
  onAddTag() {
    const { key, value } = this.thirdFormGroup.controls;

    this.selectedTags[key.value] = value.value;
    key.reset('');
    value.reset('');
  }

  onRemoveTag(tag: any) {
    delete this.selectedTags[tag];
  }
}
