import { Component, OnInit } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';
import { SINK_BACKEND_SETTINGS, SINK_BACKEND_TYPES } from 'app/common/services/sinks/sink.settings';

@Component({
  selector: 'ngx-sinks-add-component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent implements OnInit {
  strings = STRINGS;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  customSinkSettings: {};

  selectedSinkSetting: any[];

  sink: Sink;

  sinkTypesList = Object.values(SINK_BACKEND_TYPES);

  isEdit: boolean;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.sink = this.router.getCurrentNavigation().extras.state?.sink as Sink || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
  }

  ngOnInit() {
    const {name, description, backend} = !!this.sink ? this.sink : {
      name: 'my-prom-sink-',
      description: '',
      backend: '',
    } as SinkConfig<string>;
    this.firstFormGroup = this._formBuilder.group({
      name: [name, Validators.required],
      description: [description],
      backend: [backend, Validators.required],
    });

    this.secondFormGroup = this._formBuilder.group({});

    /**
     * TODO map interface to settings obj and fields OR get it from service-Backend
     * THIS IS JUST AN EXAMPLE OF HOW TO MAP WHAT COMES FROM THE BE TO SOMETHING THAT MAKES MORE
     * SENSE IN THE FRONTEND.
     */
    this.customSinkSettings = Object.keys(SINK_BACKEND_SETTINGS).reduce((accumulator, curr) => {
      accumulator[curr] = SINK_BACKEND_SETTINGS[curr].map(entry => ({
        type: entry.type,
        label: entry.title,
        prop: entry.name,
        input: entry.input,
        required: entry.required,
      }));
      return accumulator;
    }, {});
  }

  goBack() {
    this.router.navigate(['../../sinks'], {relativeTo: this.route});
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
      tags: {
        cloud: 'aws',
      },
      validate_only: false, // Apparently this guy is required..
    };
    // TODO Check this out
    // console.log(payload);
    if (this.isEdit) {
      // updating existing sink
      this.sinksService.editSink(payload).subscribe(resp => {
        this.notificationsService.success('Sink successfully created', '');
        this.goBack();
      });
    } else {
      this.sinksService.addSink(payload).subscribe(resp => {
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
        !!conf && (curr.prop in conf) && curr.prop || '',
        curr.required ? Validators.required : null,
      ];
      return accumulator;
    }, {});

    this.secondFormGroup = this._formBuilder.group(dynamicFormControls);
  }
}
