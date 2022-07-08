import { Component, TemplateRef, ViewChild } from '@angular/core';
import { NbDialogService } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { AgentKeyComponent } from '../key/agent.key.component';
import { Tags } from 'app/common/interfaces/orb/tag';

@Component({
  selector: 'ngx-agent-add-component',
  templateUrl: './agent.add.component.html',
  styleUrls: ['./agent.add.component.scss'],
})
export class AgentAddComponent {
  // page vars
  strings = { ...STRINGS.agents, stepper: STRINGS.stepper };

  isEdit: boolean;

  // templates
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  // stepper vars
  firstFormGroup: FormGroup;

  selectedTags: Tags;

  // agent vars
  agent: Agent;

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
    this.isLoading = true;

    this.agentID = this.route.snapshot.paramMap.get('id');
    this.isEdit = !!this.agentID;

    this.getAgent()
      .then((agent) => {
        this.agent = agent;
        this.initializeForms();
        this.isLoading = false;
      })
      .catch((reason) =>
        console.warn(`Couldn't fetch data. Reason: ${reason}`),
      );
  }

  newAgent() {
    return {
      name: '',
      orb_tags: {},
    } as Agent;
  }

  getAgent() {
    return new Promise<Agent>((resolve) => {
      if (this.agentID) {
        !!this.agentID &&
          this.agentsService.getAgentById(this.agentID).subscribe((resp) => {
            resolve(resp);
          });
      } else {
        resolve(this.newAgent());
      }
    });
  }

  initializeForms() {
    const { name: name, orb_tags } = this.agent;

    this.selectedTags = { ...orb_tags };

    this.firstFormGroup = this._formBuilder.group({
      name: [
        name,
        [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')],
      ],
    });
  }

  goBack() {
    this.router.navigateByUrl('/pages/fleet/agents');
  }

  wrapPayload(validate: boolean) {
    const { name } = this.firstFormGroup.controls;
    return {
      name: name.value,
      orb_tags: { ...this.selectedTags },
      validate_only: !!validate && validate, // Apparently this guy is required..
    };
  }

  openKeyModal(row: any) {
    this.dialogService
      .open(AgentKeyComponent, {
        context: { agent: row },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((resp) => {
        this.goBack();
      });
    this.notificationsService.success('Agent successfully created', '');
  }

  // saves current agent group
  onFormSubmit() {
    const payload = this.wrapPayload(false);

    if (this.isEdit) {
      this.agentsService
        .editAgent({ ...payload, id: this.agentID })
        .subscribe(() => {
          this.notificationsService.success('Agent successfully updated', '');
          this.goBack();
        });
    } else {
      this.agentsService.addAgent(payload).subscribe((resp) => {
        this.openKeyModal(resp);
      });
    }
  }
}
