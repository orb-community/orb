import { Component, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogService } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { AgentKeyComponent } from '../key/agent.key.component';


@Component({
  selector: 'ngx-agent-add-component',
  templateUrl: './agent.add.component.html',
  styleUrls: ['./agent.add.component.scss'],
})
export class AgentAddComponent {
  // page vars
  strings = {...STRINGS.agents, stepper: STRINGS.stepper};

  isEdit: boolean;

  // templates
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  // agent vars
  agent: Agent;

  selectedTags: {[propName: string]: string} = {};

  isLoading = false;

  agentID;

  constructor(
    private agentsService: AgentsService,
    private dialogService: NbDialogService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.firstFormGroup = this._formBuilder.group({
      name: ['', [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
    });

    this.secondFormGroup = this._formBuilder.group({
      key: [''],
      value: [''],
    });

    this.agentsService.clean();
    this.agent = this.router.getCurrentNavigation().extras.state?.agent as Agent || null;
    this.agentID = this.route.snapshot.paramMap.get('id');

    this.isEdit = !!this.agentID && this.router.getCurrentNavigation().extras.state?.edit as boolean;

    this.isLoading = this.isEdit;

    !!this.agentID && this.agentsService.getAgentById(this.agentID).subscribe(resp => {
      this.agent = resp;
      this.selectedTags = this.agent.orb_tags;
      this.isLoading = false;
      this.updateForm();
    });

  }

  updateForm() {
    const {name, orb_tags} = !!this.agent ? this.agent : {
      name: '',
      orb_tags: {},
    } as Agent;

    this.firstFormGroup.controls.name.patchValue(name);

    this.secondFormGroup.setValue({
      orb_tags: !!orb_tags ? Object.entries(orb_tags).map(([key, value]) => ({[key]: value})) : [],
      key: '',
      value: '',
    });

    this.agentsService.clean();
  }

  resetFormValues() {
    const {name, orb_tags} = !!this.agent ? this.agent : {
      name: '',
      orb_tags: {},
    } as Agent;

    this.firstFormGroup.setValue({name: name});

    this.secondFormGroup.controls.orb_tags.setValue(
      Object.keys(orb_tags).map(key => ({[key]: orb_tags[key]})),
    );

    this.agentsService.clean();
  }

  goBack() {
    this.router.navigateByUrl('/pages/fleet/agents');
  }

  validateSelectedTags() {
    return Object.keys(this.selectedTags).length !== 0;
  }

  // TODO - this method can be refactored and has multiple occurrences in ui components
  checkValidName() {
    const { value } = this.secondFormGroup.controls.key;
    if (value === '') return false;
    return !Object.keys(this.selectedTags).find(name => name === value) !== undefined;
  }

  onAddTag() {
    const {orb_tags, key, value} = this.secondFormGroup.controls;
    // sanitize minimally anyway
    if (key?.value && key.value !== '') {
      if (value?.value && value.value !== '') {
        // key and value fields
        key.reset('');
        value.reset('');
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
      tags.patchValue(tagsList.slice(0, indexToRemove).concat(tagsList.slice(indexToRemove + 1)));
    }
  }

  wrapPayload(validate: boolean) {
    const {name} = this.firstFormGroup.controls;
    const {tags: {value: tagsList}} = this.secondFormGroup.controls;
    const tagsObj = tagsList.reduce((prev, curr) => {
      for (const [key, value] of Object.entries(curr)) {
        prev[key] = value;
      }
      return prev;
    }, {});
    return {
      name: name.value,
      orb_tags: {...tagsObj},
      validate_only: !!validate && validate, // Apparently this guy is required..
    };
  }

  openKeyModal(row: any) {
    this.dialogService.open(AgentKeyComponent, {
      context: { agent: row },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      this.goBack();
    });
    this.notificationsService.success('Agent successfully created', '');
  }

  // saves current agent group
  onFormSubmit() {
    const payload = this.wrapPayload(false);

    if (this.isEdit) {
      this.agentsService.editAgent({...payload, id: this.agentID}).subscribe(() => {
        this.notificationsService.success('Agent successfully updated', '');
        this.goBack();
      });
    } else {
      this.agentsService.addAgent(payload).subscribe((resp) => {
        this.openKeyModal(resp.body);
      });
    }
  }

}
