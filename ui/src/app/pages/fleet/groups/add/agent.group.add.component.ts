import { AfterViewInit, Component, TemplateRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { TagMatch } from 'app/common/interfaces/orb/tag.match.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';


@Component({
  selector: 'ngx-agent-group-add-component',
  templateUrl: './agent.group.add.component.html',
  styleUrls: ['./agent.group.add.component.scss'],
})
export class AgentGroupAddComponent implements AfterViewInit {
  // page vars
  strings = { ...STRINGS.agentGroups, stepper: STRINGS.stepper };

  isEdit: boolean;

  columnMode = ColumnMode;

  columns: TableColumn[];

  // templates
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
    },
    {
      id: '1',
      label: 'Tags',
      prop: 'tags',
      selected: false,
    },
  ];

  expanded: boolean;

  // stepper vars
  firstFormGroup: FormGroup;

  secondFormGroup: FormGroup;

  // agent vars
  agentGroup: AgentGroup;

  selectedTags: { [propName: string]: string };

  matchingAgents: Agent[];

  tagMatch: TagMatch = {};

  isLoading = false;

  agentGroupID;

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentsService: AgentsService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.isLoading = true;

    this.selectedTags = {};
    this.tagMatch.total = this.tagMatch.online = 0;
    this.expanded = false;
    this.agentsService.clean();
    this.agentGroupsService.clean();

    this.agentGroupID = this.route.snapshot.paramMap.get('id');
    this.isEdit = !!this.agentGroupID;

    this.getAgentGroup()
      .then((agentGroup) => {
        this.agentGroup = agentGroup;
        this.selectedTags = agentGroup.tags;
        this.initializeForms();
        this.isLoading = false;
      })
      .then(() => this.updateMatches())
      .catch(reason => console.warn(`Couldn't retrieve data. Reason: ${ reason }`));
  }

  initializeForms() {
    const { name, description } = this.agentGroup;

    this.firstFormGroup = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      description: [description],
    });

    this.secondFormGroup = this._formBuilder.group({
      key: [''],
      value: [''],
    });
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Agent Name',
        resizeable: false,
        flexGrow: 1,
        minWidth: 90,
      },
      {
        prop: 'orb_tags',
        name: 'Tags',
        resizeable: false,
        minWidth: 100,
        flexGrow: 2,
        cellTemplate: this.agentTagsTemplateCell,
      },
      {
        prop: 'state',
        name: 'Status',
        minWidth: 90,
        flexGrow: 1,
        cellTemplate: this.agentStateTemplateRef,
      },
      {
        name: 'Last Activity',
        prop: 'ts_last_hb',
        minWidth: 130,
        resizeable: false,
        sortable: false,
        flexGrow: 1,
      },
    ];
  }

  newAgentGroup() {
    return {
      name: '',
      description: '',
      tags: {},
    } as AgentGroup;
  }

  getAgentGroup() {
    return new Promise<AgentGroup>(resolve => {
      if (this.agentGroupID) {
        this.agentGroupsService.getAgentGroupById(this.agentGroupID).subscribe(resp => {
          resolve(resp);
        });
      } else {
        resolve(this.newAgentGroup());
      }
    });
  }

  goBack() {
    this.router.navigateByUrl('/pages/fleet/groups');
  }

  // addTag button should be [disabled] = `$sf.controls.key.value !== ''`
  onAddTag() {
    const { key, value } = this.secondFormGroup.controls;

    this.selectedTags[key.value] = value.value;

    // key and value fields
    key.reset('');
    value.reset('');

    this.updateMatches();
  }

  onRemoveTag(tag: any) {
    delete this.selectedTags[tag];

    this.updateMatches();
  }

  // query agent group matches
  updateMatches() {
    const tagMatches = new Promise<TagMatch>(resolve => {
      const name = this.firstFormGroup.controls.name.value;
      if (name !== '' && Object.keys(this.selectedTags).length !== 0) {
        const payload = this.wrapPayload(true);
        // just validate and get matches summary
        this.agentGroupsService.validateAgentGroup(payload).subscribe((resp: any) => {
          resolve({
            total: resp.body.matching_agents.total,
            online: resp.body.matching_agents.online,
          });
        });
      } else {
        resolve({ total: 0, online: 0 });
      }
    });

    const matchingAgents = new Promise<Agent[]>(resolve => {
      if (Object.keys(this.selectedTags).length !== 0) {
        this.agentsService.getMatchingAgents(this.selectedTags).subscribe(
          resp => {
            resolve(resp.agents);
          });
      } else {
        resolve([]);
      }
    });

    Promise.all([tagMatches, matchingAgents]).then(responses => {
      const summary = responses[0] as TagMatch;
      const matches = responses[1] as Agent[];

      this.tagMatch = summary;
      this.matchingAgents = matches;
    }).catch(reason => console.warn(`Couldn't retrieve data. Reason: ${ reason }`));
  }

  toggleExpandMatches() {
    this.expanded = !this.expanded;
    !!this.expanded && this.updateMatches();
  }

  wrapPayload(validate: boolean) {
    const { name, description } = this.firstFormGroup.controls;
    return {
      name: name.value,
      description: description.value,
      tags: { ...this.selectedTags },
      validate_only: !!validate && validate, // Apparently this guy is required..
    };
  }

  // saves current agent group
  onFormSubmit() {
    // validate:false
    const payload = this.wrapPayload(false);

    // // remove line bellow
    // console.log(payload)
    if (this.isEdit) {
      this.agentGroupsService.editAgentGroup({ ...payload, id: this.agentGroupID }).subscribe(() => {
        this.notificationsService.success('Agent Group successfully updated', '');
        this.goBack();
      });
    } else {
      this.agentGroupsService.addAgentGroup(payload).subscribe(() => {
        this.notificationsService.success('Agent Group successfully created', '');
        this.goBack();
      });
    }
  }

}
