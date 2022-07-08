import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnChanges,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { TagMatch } from 'app/common/interfaces/orb/tag.match.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import {
  ColumnMode,
  DatatableComponent,
  TableColumn,
} from '@swimlane/ngx-datatable';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-agent-group-add-component',
  templateUrl: './agent.group.add.component.html',
  styleUrls: ['./agent.group.add.component.scss'],
})
export class AgentGroupAddComponent
  implements OnInit, OnChanges, AfterViewInit {
  // page vars
  strings = { ...STRINGS.agentGroups, stepper: STRINGS.stepper };

  isEdit: boolean;

  columnMode = ColumnMode;

  columns: TableColumn[];

  // table
  @ViewChild(DatatableComponent) table: DatatableComponent;

  // templates
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;

  @ViewChild('agentStateTemplateCell') agentStateTemplateRef: TemplateRef<any>;

  @ViewChild('agentLastHBTemplateCell') agentLastHBTemplateRef: TemplateRef<
    any
  >;

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
    private cdr: ChangeDetectorRef,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.isLoading = true;

    this.selectedTags = {};
    this.tagMatch.total = this.tagMatch.online = 0;
    this.expanded = false;

    this.agentGroupID = this.route.snapshot.paramMap.get('id');
    this.isEdit = !!this.agentGroupID;
  }

  ngOnInit() {
    this.getAgentGroup()
      .then((agentGroup) => {
        this.agentGroup = agentGroup;
        this.selectedTags = agentGroup.tags;
        this.initializeForms();
        this.isLoading = false;
      })
      .then(() => this.updateMatches())
      .catch((reason) =>
        console.warn(`Couldn't retrieve data. Reason: ${reason}`),
      );
  }

  ngOnChanges() {
    this.table.rows = this.matchingAgents;
    this.table.recalculate();
  }

  initializeForms() {
    const { name: name, description } = this.agentGroup;

    this.firstFormGroup = this._formBuilder.group({
      name: [
        name,
        [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')],
      ],
      description: [description],
    });
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Agent Name',
        flexGrow: 2,
        canAutoResize: true,
        resizeable: false,
        minWidth: 90,
        width: 120,
        maxWidth: 200,
      },
      {
        prop: 'combined_tags',
        name: 'Tags',
        flexGrow: 6,
        resizeable: false,
        canAutoResize: true,
        minWidth: 300,
        width: 450,
        maxWidth: 1000,
        cellTemplate: this.agentTagsTemplateCell,
        comparator: (a, b) =>
          Object.entries(a)
            .map(([key, value]) => `${key}:${value}`)
            .join(',')
            .localeCompare(
              Object.entries(b)
                .map(([key, value]) => `${key}:${value}`)
                .join(','),
            ),
      },
      {
        prop: 'state',
        name: 'Status',
        flexGrow: 1,
        canAutoResize: true,
        resizeable: false,
        minWidth: 90,
        maxWidth: 150,
        cellTemplate: this.agentStateTemplateRef,
      },
      {
        name: 'Last Activity',
        prop: 'ts_last_hb',
        cellTemplate: this.agentLastHBTemplateRef,
        flexGrow: 2,
        resizeable: false,
        canAutoResize: true,
        minWidth: 180,
        width: 250,
        maxWidth: 400,
        sortable: false,
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
    return new Promise<AgentGroup>((resolve) => {
      if (this.agentGroupID) {
        this.agentGroupsService
          .getAgentGroupById(this.agentGroupID)
          .subscribe((resp) => {
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

  // query agent group matches
  updateMatches() {
    const tagMatches = new Promise<TagMatch>((resolve) => {
      const name = this.firstFormGroup.controls.name.value;
      if (name !== '' && Object.keys(this.selectedTags).length !== 0) {
        const payload = this.wrapPayload(true);
        // just validate and get matches summary
        this.agentGroupsService
          .validateAgentGroup(payload)
          .subscribe((resp: any) => {
            resolve({
              total: resp.body.matching_agents.total,
              online: resp.body.matching_agents.online,
            });
          });
      } else {
        resolve({ total: 0, online: 0 });
      }
    });

    const matchingAgents = new Promise<Agent[]>((resolve) => {
      if (Object.keys(this.selectedTags).length !== 0) {
        this.agentsService
          .getMatchingAgents(this.selectedTags)
          .subscribe((agents) => {
            resolve(agents);
          });
      } else {
        resolve([]);
      }
    });

    Promise.all([tagMatches, matchingAgents])
      .then((responses) => {
        const summary = responses[0] as TagMatch;
        const matches = responses[1] as Agent[];

        this.tagMatch = summary;
        this.matchingAgents = matches;
        this.cdr.markForCheck();
      })
      .catch((reason) =>
        console.warn(`Couldn't retrieve data. Reason: ${reason}`),
      );
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

  selectedTagsValid() {
    return Object.keys(this.selectedTags).length > 0;
  }

  // saves current agent group
  onFormSubmit() {
    // validate:false
    const payload = this.wrapPayload(false);

    // // remove line bellow
    // console.log(payload)
    if (this.isEdit) {
      this.agentGroupsService
        .editAgentGroup({ ...payload, id: this.agentGroupID })
        .subscribe(() => {
          this.notificationsService.success(
            'Agent Group successfully updated',
            '',
          );
          this.goBack();
        });
    } else {
      this.agentGroupsService.addAgentGroup(payload).subscribe(() => {
        this.notificationsService.success(
          'Agent Group successfully created',
          '',
        );
        this.goBack();
      });
    }
  }
}
