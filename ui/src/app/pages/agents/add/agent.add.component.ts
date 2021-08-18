import {Component, OnInit} from '@angular/core';

import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {SinksService} from 'app/common/services/sinks/sinks.service';
import {ActivatedRoute, Router} from '@angular/router';
import {Sink} from 'app/common/interfaces/orb/sink.interface';
import {STRINGS} from 'assets/text/strings';
import {sinkTypesList} from 'app/pages/sinks/sinks.component';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';


const SETTINGS_EXAMPLE = {
  prometheus: [
    {
      type: 'text',
      input: 'text',
      title: 'Remote Host',
      name: 'remote_host',
      required: true,
    },
    {
      type: 'text',
      input: 'text',
      title: 'Username',
      name: 'username',
      required: true,
    },
    {
      type: 'password',
      input: 'text',
      title: 'Password',
      name: 'password',
      required: true,
    },
  ],
};


@Component({
  selector: 'ngx-agent-add-component',
  templateUrl: './agent.add.component.html',
  styleUrls: ['./agent.add.component.scss'],
})
export class AgentAddComponent implements OnInit {
  // stepper vars
  firstFormGroup: FormGroup;
  secondFormGroup: FormGroup;
  isEditable = false;

  strings = STRINGS;

  customSinkSettings: {};
  selectedSinkSetting: any[];

  sinkForm = {
    name: '',
    description: '',
    backend: sinkTypesList.prometheus,
    config: {
      host_name: '',
      username: '',
      password: '',
    },
    tags: {},
  };
  sink: Sink;

  sinkTypesList = Object.values(sinkTypesList);

  isEdit: boolean;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.sink = this.router.getCurrentNavigation().extras.state?.sink as Sink || null;

  }

  ngOnInit() {
    this.firstFormGroup = this._formBuilder.group({
      name: ['my-prom-sink-', Validators.required],
      description: [''],
      backend: ['', Validators.required],
    });

    this.secondFormGroup = this._formBuilder.group({});

    /***TODO THIS IS JUST AN EXAMPLE OF HOW TO MAP WHAT COMES FROM THE BE TO SOMETHING THAT MAKES MORE
     * SENSE IN THE FRONTEND.
     */

    this.customSinkSettings = Object.keys(SETTINGS_EXAMPLE).reduce((accumulator, curr) => {
      accumulator[curr] = SETTINGS_EXAMPLE[curr].map(entry => ({
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
    this.sinksService.addSink(payload).subscribe(resp => {
      this.notificationsService.success('Sink successfully created', '');
      this.goBack();
    });
  }

  onSinkTypeSelected(selectedValue) {
    this.selectedSinkSetting = this.customSinkSettings[selectedValue];
    const dynamicFormControls = this.selectedSinkSetting.reduce((accumulator, curr) => {
      accumulator[curr.prop] = ['', curr.required ? Validators.required : null];
      return accumulator;
    }, {});
    this.secondFormGroup = this._formBuilder.group(dynamicFormControls);
  }
}
