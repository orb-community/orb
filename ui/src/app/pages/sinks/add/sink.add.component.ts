import { Component, OnDestroy, OnInit } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { forkJoin, of, Subscription } from 'rxjs';
import { SinkFeature } from 'app/common/interfaces/orb/sink/sink.feature.interface';

@Component({
  selector: 'ngx-sink-add-component',
  templateUrl: './sink.add.component.html',
  styleUrls: ['./sink.add.component.scss'],
})
export class SinkAddComponent implements OnInit, OnDestroy {
  strings = STRINGS;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  thirdFormGroup: FormGroup;

  customSinkSettings: {};

  backends: { [propName: string]: SinkFeature };

  backendConfig: any[];

  selectedTags: { [propName: string]: string };

  sink: Sink;

  sinkID: string;

  sinkTypesList = [];

  isEdit: boolean;

  isLoading = false;

  subscription: Subscription;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.isLoading = true;
    this.sinkID = this.route.snapshot.paramMap.get('id');
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean || !!this.sinkID;
  }

  ngOnInit() {
    const sink$ = this.isEdit ?
      // retrieve sink by id
      this.sinksService.getSinkById(this.sinkID)
      :
      // use a blank sink
      of({
        name: '',
        description: '',
        backend: 'prometheus', // default sink
        tags: {},
      } as Sink);

    this.subscription = forkJoin(
      {
        sink: sink$,
        backends: this.sinksService.getSinkBackends(),
      })
      .subscribe(values => {
        const { sink: { backend: backend } } = { sink: this.sink, backends: this.backends } = values;

        this.isLoading = false;

        this.initializeForms();

        if (backend !== '') this.onSinkTypeSelected(backend);
      });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  initializeForms() {
    const { name, description, backend, tags } = this.sink;

    this.firstFormGroup = this._formBuilder.group({
      name: [name, [this.sinksService.sinkNameValidator, Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      description: [description],
      backend: [backend, Validators.required],
    });

    this.selectedTags = { ...tags };

    this.thirdFormGroup = this._formBuilder.group({
      key: [''],
      value: [''],
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
      config: this.backendConfig.reduce((accumulator, current) => {
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

  onSinkTypeSelected(selectedBackend) {
    // SinkConfig<string> being the generic of all other `sinkTypes`.
    const conf = !!this.sink &&
      (selectedBackend === this.sink.backend) &&
      this.sink?.config || null;

    this.backendConfig = this.backends[selectedBackend].config;

    const dynamicFormControls = this.backendConfig.reduce((accumulator, curr) => {
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
    const { value } = this.thirdFormGroup?.controls?.key;
    const hasTagForKey = Object.keys(this.selectedTags).find(key => key === value);
    return value && value !== '' && !hasTagForKey;
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
