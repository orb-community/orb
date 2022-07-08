import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { SinkConfig } from 'app/common/interfaces/orb/sink/sink.config.interface';
import { SinkFeature } from 'app/common/interfaces/orb/sink/sink.feature.interface';

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

  customSinkSettings: {};

  selectedSinkSetting: any[];

  selectedTags: { [propName: string]: string };

  sink: Sink;

  sinkID: string;

  sinkTypesList = [];

  isEdit: boolean;

  isLoading = false;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.isLoading = true;
    this.sinkID = this.route.snapshot.paramMap.get('id');
    this.isEdit =
      (this.router.getCurrentNavigation().extras.state?.edit as boolean) ||
      !!this.sinkID;

    Promise.all([this.getSink(), this.getSinkBackends()])
      .then((responses) => {
        const { backend } = (this.sink = responses[0]);
        const backends = responses[1];

        this.sinkTypesList = backends.map((entry) => entry.backend);
        this.customSinkSettings = this.sinkTypesList.reduce(
          (accumulator, curr) => {
            const index = backends.findIndex((entry) => entry.backend === curr);
            accumulator[curr] = backends[index].config.map((entry) => ({
              type: entry.type,
              label: entry.title,
              prop: entry.name,
              input: entry.input,
              required: entry.required,
            }));
            return accumulator;
          },
          {},
        );

        this.initializeForms();

        this.isLoading = false;
        if (backend !== '') this.onSinkTypeSelected(backend);
      })
      .catch((reason) =>
        console.warn(`Couldn't retrieve data. Reason: ${reason}`),
      );
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
    const { name: name, description, backend, tags } = this.sink;

    this.firstFormGroup = this._formBuilder.group({
      name: [
        name,
        [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')],
      ],
      description: [description],
      backend: [backend, Validators.required],
    });

    this.selectedTags = { ...tags };
  }

  getSink() {
    return new Promise<Sink>((resolve) => {
      if (this.sinkID) {
        this.sinksService.getSinkById(this.sinkID).subscribe((resp) => {
          resolve(resp);
        });
      } else {
        resolve(this.newSink());
      }
    });
  }

  getSinkBackends() {
    return new Promise<SinkFeature[]>((resolve) => {
      this.sinksService.getSinkBackends().subscribe((backends) => {
        resolve(backends);
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
        accumulator[current.prop] = this.secondFormGroup.controls[
          current.prop
        ].value;
        return accumulator;
      }, {}),
      tags: { ...this.selectedTags },
    };

    if (this.isEdit) {
      // updating existing sink
      this.sinksService
        .editSink({ ...payload, id: this.sinkID })
        .subscribe(() => {
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
    const conf =
      (!!this.sink &&
        this.isEdit &&
        selectedValue === this.sink.backend &&
        this.sink?.config &&
        (this.sink.config as SinkConfig<string>)) ||
      null;

    this.selectedSinkSetting = this.customSinkSettings[selectedValue];

    const dynamicFormControls = this.selectedSinkSetting.reduce(
      (accumulator, curr) => {
        accumulator[curr.prop] = [
          (!!conf && curr.prop in conf && conf[curr.prop]) || '',
          curr.required ? Validators.required : null,
        ];
        return accumulator;
      },
      {},
    );

    this.secondFormGroup = this._formBuilder.group(dynamicFormControls);
  }
}
