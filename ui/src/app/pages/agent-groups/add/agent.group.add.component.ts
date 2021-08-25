import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { TagMatch } from 'app/common/interfaces/orb/tag.match.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination';
import { ColumnMode, TableColumn } from '@swimlane/ngx-datatable';


@Component({
  selector: 'ngx-agent-group-add-component',
  templateUrl: './agent.group.add.component.html',
  styleUrls: ['./agent.group.add.component.scss'],
})
export class AgentGroupAddComponent implements OnInit {
  // page vars
  strings = {...STRINGS.agents, stepper: STRINGS.stepper};

  isEdit: boolean;

  // // expandable table vars
  // tableConfig: TableConfig = {
  //   colNames: ['Agent Name', 'Tags', 'Status', 'Last Activity'],
  //   keys: ['name', 'agent_tags', 'state', 'ts_lst_hb'],
  // };
  columnMode = ColumnMode;
  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<Agent>;

  searchPlaceholder = 'Search by name';
  filterSelectedIndex = '0';

  // templates
  @ViewChild('agentsTemplateCell') agentsTemplateCell: TemplateRef<any>;
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;
  @ViewChild('addAgentTemplateRef') addAgentTemplateRef: TemplateRef<any>;
  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

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

  addForm: FormGroup;
  // agent vars
  agentGroup: AgentGroup;

  matchingAgents: Agent[];

  tagMatch: TagMatch = {};

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentsService: AgentsService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentsService.clean();
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

    this.addForm = this._formBuilder.group({
      firstFormGroup: this.firstFormGroup,
      secondFormGroup: this.secondFormGroup,
    });

    this.tagMatch.total = this.tagMatch.online = 0;
    this.expanded = false;

    this.agentGroupsService.clean();
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
        tags.reset([{[key.value]: value.value}].concat(tags.value));
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

  // query agent group matches
  updateTagMatches() {
    // validate:true
    const payload = this.wrapPayload(true);
    // // remove line bellow
    // console.log(payload)

    // just validate and get matches summary
    this.agentGroupsService.validateAgentGroup(payload).subscribe((resp: any) => {
      this.tagMatch = {
        total: resp.matchingAgents.total,
        online: resp.matchingAgents.online,
      };

      this.notificationsService.success(this.strings.match.updated, '');
    });
  }

  updateMatchingAgents(pageInfo: NgxDatabalePageInfo = null) {
    const isFilter = pageInfo === null;
    if (isFilter) {
      pageInfo = {
        offset: this.paginationControls.offset,
        limit: this.paginationControls.limit,
      };
      if (this.paginationControls.name?.length > 0) pageInfo.name = this.paginationControls.name;
      if (this.paginationControls.tags?.length > 0) pageInfo.tags = this.paginationControls.tags;
    }

    this.loading = true;
    this.agentsService.getAgents(pageInfo, isFilter).subscribe(
      (resp: OrbPagination<Agent>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo.offset;
        this.loading = false;
      },
    );
    // update list of agents
    // this.agentsService.getAgents();
    // this.matchingAgents = new Array(10)
    //   .fill(null)
    //   .map((_, i) => (
    //     {
    //       name: `Lorem Ipsum ${i}`,
    //       agent_tags: {cloud: `aws-${i}`},
    //       state: ['new', 'online', 'offline', 'stale'][i % 4],
    //       ts_lst_hb: `${+new Date()}`,
    //     }
    //   ));
    // update matching agent table
    // this.page = {
    //   offset: 0,
    //   limit: 10,
    //   total: 10,
    //   rows: this.matchingAgents,
    // };
  }

  toggleExpandMatches() {
    this.expanded = !this.expanded;
    !!this.expanded && this.updateMatchingAgents();
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

  // saves current agent group
  onFormSubmit() {
    // validate:false
    const payload = this.wrapPayload(false);

    // // remove line bellow
    // console.log(payload)
    this.agentGroupsService.addAgentGroup(payload).subscribe(() => {
      this.notificationsService.success(this.strings.add.success, '');
      this.goBack();
    });
  }

}
