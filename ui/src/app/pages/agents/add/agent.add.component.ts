import { Component, OnInit } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { TagMatch } from 'app/common/interfaces/orb/tag.match.interface';

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

  strings = {...STRINGS.agents, stepper: STRINGS.stepper};

  agentGroup: AgentGroup;

  matchingAgents: [];

  tagMatch: TagMatch = {};

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
      tags: [[], Validators.minLength(1)],
      key: [''],
      value: [''],
    });

    this.tagMatch.total = this.tagMatch.online = 0;
  }

  goBack() {
    this.router.navigate(['../../agents'], {relativeTo: this.route});
  }

  // addTag button should be [disabled] = `$sf.controls.key.value !== ''`
  onAddTag() {
    const {tags, key, value} = this.secondFormGroup.controls;
    // sanitize minimally anyway
    if (key?.value && key.value !== '') {
      if (value?.value && value.value !== '') {
        // key and value fields
        tags.setValue([{[key.value]: value.value}].concat(tags.value));
        key.reset('');
        value.reset('');
        this.updateTagMatches();
      }
    } else {
      // TODO remove this else clause and error
      console.error('This shouldn\'t be happening');
    }
  }

  onRemoveTag(tag: any) {
    const {tags, tags: {value: tagsList}} = this.secondFormGroup.controls;
    const indexToRemove = tagsList.indexOf(tag);

    if (indexToRemove >= 0) {
      tags.setValue(tagsList.slice(0, indexToRemove).concat(tagsList.slice(indexToRemove + 1)));
      this.updateTagMatches();
    }
  }

  wrapPayload(validate: boolean) {
    const {name, description} = this.firstFormGroup.controls;
    const {tags: {value: tagsList}} = this.secondFormGroup.controls;
    const tagsObj = tagsList.reduce((prev, curr) => {
      for (const [key, value] of Object.entries(curr)) {
        prev[key] = value;
      }
      return prev;
    }, {});

    return {
      name: name.value,
      description: description.value,
      tags: {...tagsObj},
      validate_only: !!validate && validate, // Apparently this guy is required..
    };
  }

  // query agent group matches
  updateTagMatches() {
    // validate:true
    const payload = this.wrapPayload(true);
    // // remove line bellow
    // console.log(payload)

    // just validate and get matches summary
    this.agentsService.validateAgentGroup(payload).subscribe((resp: any) => {
      this.tagMatch = {
        total: resp.matchingAgents.total,
        online: resp.matchingAgents.online,
      };
      this.notificationsService.success(this.strings.match.updated, '');
    });
  }

  // saves current agent group
  onFormSubmit() {
    // validate:false
    const payload = this.wrapPayload(false);

    // // remove line bellow
    // console.log(payload)
    this.agentsService.addAgentGroup(payload).subscribe(resp => {
      this.notificationsService.success(this.strings.add.success, '');
      this.goBack();
    });
  }

}
