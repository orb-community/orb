import { Component, OnInit } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { tagInputValidator } from 'app/shared/directives/tag-input.validator';

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

  agentGroup: AgentGroup;

  isEdit: boolean;

  constructor(
    private agentsService: AgentsService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentGroup = this.router.getCurrentNavigation().extras.state?.agentGroup as AgentGroup || null;
  }

  ngOnInit() {
    this.firstFormGroup = this._formBuilder.group({
      name: ['', Validators.required],
      description: [''],
    });

    this.secondFormGroup = this._formBuilder.group({
      tags: [{}, tagInputValidator()],
      key: [''],
      value: [''],
    });

  }

  goBack() {
    this.router.navigate(['../../agents'], {relativeTo: this.route});
  }

  onAddTag() {
    const {tags, key, value} = this.secondFormGroup.controls;
    tags.value[key.value] = value.value;
  }

  onRemoveTag() {

  }

  onFormSubmit() {
    const payload = {
      name: this.firstFormGroup.controls.name.value,
      description: this.firstFormGroup.controls.description.value,
      // TODO tag input
      tags: {
        hardcoded: 'payload',
      },
      validate_only: false, // Apparently this guy is required..
    };

    // // TODO remove line bellow
    // console.log(payload);

    this.agentsService.addAgentGroup(payload).subscribe(resp => {
      this.notificationsService.success('Agent Group successfully created', '');
      this.goBack();
    });
  }

}
