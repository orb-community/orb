import {Component, OnInit} from '@angular/core';

import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {SinksService} from 'app/common/services/sinks/sinks.service';
import {ActivatedRoute, Router} from '@angular/router';
import {Sink} from 'app/common/interfaces/orb/sink.interface';
import {STRINGS} from 'assets/text/strings';
import {sinkTypesList} from 'app/pages/sinks/sinks.component';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';

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

  strings = STRINGS.agents;

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
      name: ['', Validators.required],
      description: [''],
    });

    this.secondFormGroup = this._formBuilder.group({});

  }

  goBack() {
    this.router.navigate(['../../sinks'], {relativeTo: this.route});
  }

  onFormSubmit() {
    // const payload = {
    //   name: this.firstFormGroup.controls.name.value,
    //   backend: this.firstFormGroup.controls.backend.value,
    //   description: this.firstFormGroup.controls.description.value,
    //   config: this.selectedSinkSetting.reduce((accumulator, current) => {
    //     accumulator[current.prop] = this.secondFormGroup.controls[current.prop].value;
    //     return accumulator;
    //   }, {}),
    //   tags: {
    //     cloud: 'aws',
    //   },
    //   validate_only: false, // Apparently this guy is required..
    // };
    // // TODO Check this out
    // // console.log(payload);
    // this.sinksService.addSink(payload).subscribe(resp => {
    //   this.notificationsService.success('Sink successfully created', '');
    //   this.goBack();
    // });
  }

}
