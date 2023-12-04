import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NbDialogService } from '@nebular/theme';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { Tags } from 'app/common/interfaces/orb/tag';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { AgentResetComponent } from 'app/pages/fleet/agents/reset/agent.reset.component';

@Component({
  selector: 'ngx-agent-information',
  templateUrl: './agent-information.component.html',
  styleUrls: ['./agent-information.component.scss'],
})
export class AgentInformationComponent implements OnInit, OnChanges {
  @Input() agent: Agent;

  isResetting: boolean;

  agentStates = AgentStates;

  editMode: boolean;

  formGroup: FormGroup;

  selectedTags: Tags;

  isRequesting: boolean;

  @Output()
  refreshRequests = new EventEmitter<boolean>();

  constructor(
    protected agentsService: AgentsService,
    protected notificationService: NotificationsService,
    private fb: FormBuilder,
    private orb: OrbService,
    private dialogService: NbDialogService,
  ) {
    this.isResetting = false;
    this.isRequesting = false;
    this.editMode = false;
    this.updateForm();
  }

  ngOnInit(): void {
    this.selectedTags = this.agent?.orb_tags || {};
  }
  updateForm() {
    if (this.editMode) {
      const { name, orb_tags } = this.agent;
      this.formGroup = this.fb.group({
        name: [
          name,
          [
            Validators.required,
            Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$'),
            Validators.maxLength(64),
            Validators.minLength(2),
          ],
        ],
      });
      this.selectedTags = {...orb_tags} || {};
    } else {
      this.formGroup = this.fb.group({
        name: null,
      });
    }
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.editMode) {
      this.toggleEdit(changes.editMode.currentValue);
    }
    if (changes?.policy) {
      this.selectedTags = this.agent?.orb_tags || {};
    }
  }

  resetAgent() {
    if (!this.isResetting) {
      this.isResetting = true;
      this.agentsService.resetAgent(this.agent.id).subscribe(() => {
        this.isResetting = false;
        this.notifyResetSuccess();
      });
    }
  }

  getAgentVersion() {
    const agentVersion = this.agent?.agent_metadata?.orb_agent?.version;

    return agentVersion ? agentVersion : '-';
  }

  notifyResetSuccess() {
    this.notificationService.success('Agent Reset Requested', '');
  }
  toggleEdit(value) {

    this.editMode = value;
    if (this.editMode) {
      this.orb.pausePolling();
    } else {
      this.orb.startPolling();
    }
    this.updateForm();
  }
  canSave() {
    if (this.formGroup.status === 'VALID') {
      return true;
    }
    return false;
  }
  save () {
    this.isRequesting = true;
    const name = this.formGroup.controls.name.value;
    const payload = {
      name: name,
      orb_tags: { ...this.selectedTags },
    };
    this.agentsService.editAgent({ ...payload, id: this.agent.id }).subscribe(() => {
      this.notificationService.success('Agent successfully updated', '');
      this.orb.refreshNow();
      this.toggleEdit(false);
      this.isRequesting = false;
    },
    (error) => {
      this.isRequesting = false;
    });
  }

  hasChanges() {
    const name = this.formGroup.controls.name.value;

    const selectedTags = JSON.stringify(this.selectedTags);
    const orb_tags = JSON.stringify(this.agent.orb_tags);

    if (this.agent.name !== name || selectedTags !== orb_tags) {
      return true;
    }
    return false;
  }
  onOpenResetAgents() {
    const agent = this.agent;
    this.dialogService
      .open(AgentResetComponent, {
        context: { agent },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.resetAgent();
          this.orb.refreshNow();
        }
      });
  }
}
