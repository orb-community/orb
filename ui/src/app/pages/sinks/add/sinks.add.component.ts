import {Component} from '@angular/core';

import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {SinksService} from 'app/common/services/sinks/sinks.service';
import {ActivatedRoute, Router} from '@angular/router';
import {Sink} from 'app/common/interfaces/orb/sink.interface';
import {STRINGS} from 'assets/text/strings';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {SinkConfig} from 'app/common/interfaces/orb/sink.config/sink.config.interface';

@Component({
  selector: 'ngx-sinks-add-component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent {
  strings = STRINGS;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  customSinkSettings: {};

  selectedSinkSetting: any[];

  sink: Sink;

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
    this.sink = this.router.getCurrentNavigation().extras.state?.sink as Sink || null;
    this.isEdit = this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.getSinkBackends();
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
      const {name, description, backend} = !!this.sink ? this.sink : {
        name: '',
        description: '',
        backend: 'prometheus', // default sink
      } as SinkConfig<string>;
      this.firstFormGroup = this._formBuilder.group({
        name: [name, Validators.required],
        description: [description],
        backend: [backend, Validators.required],
      });

      this.isEdit && this.firstFormGroup.controls.backend.disable();

      // builds secondFormGroup
      this.onSinkTypeSelected(backend);
      this.isLoading = false;
    });
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
