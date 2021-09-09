import { AfterViewInit, Component, OnInit, TemplateRef, ViewChild } from '@angular/core';
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


@Component({
  selector: 'ngx-agent-group-add-component',
  templateUrl: './agent.group.add.component.html',
  styleUrls: ['./agent.group.add.component.scss'],
})
export class AgentGroupAddComponent implements OnInit, AfterViewInit {
  // page vars
  strings = {...STRINGS.agentGroups, stepper: STRINGS.stepper};

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

  matchingAgents: Agent[];

  tagMatch: TagMatch = {};

  isLoading = false;

  agentGroupID;

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentsService: AgentsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentsService.clean();
    this.agentGroup = this.router.getCurrentNavigation().extras.state?.agentGroup as AgentGroup || null;
    this.agentGroupID = this.route.snapshot.paramMap.get('id');
    !!this.agentGroupID && this.agentGroupsService.getAgentGroupById(this.agentGroupID).subscribe(resp => {
      this.agentGroup = resp.agentGroup;
      this.isLoading = false;
    });
    this.isEdit = !!this.agentGroupID && this.router.getCurrentNavigation().extras.state?.edit as boolean;
    this.isLoading = this.isEdit;
  }

  ngOnInit() {
    const {name, description, tags} = !!this.agentGroup ? this.agentGroup : {
      name: '',
      description: '',
      tags: {},
    } as AgentGroup;
    this.firstFormGroup = this._formBuilder.group({
      name: [name, Validators.required],
      description: [description],
    });

    this.secondFormGroup = this._formBuilder.group({
      tags: [Object.keys(tags).map(key => ({[key]: tags[key]})) || [],
        Validators.minLength(1)],
      key: [''],
      value: [''],
    });

    this.tagMatch.total = this.tagMatch.online = 0;
    this.expanded = false;

    this.agentGroupsService.clean();
  }

  resetFormValues() {
    const {name, description, tags} = !!this.agentGroup ? this.agentGroup : {
      name: '',
      description: '',
      tags: {},
    } as AgentGroup;

    this.firstFormGroup.setValue({name: name, description: description});

    this.secondFormGroup.controls.tags.setValue(
      Object.keys(tags).map(key => ({[key]: tags[key]})));

    this.updateTagMatches();

    this.updateMatchingAgents();

    this.agentGroupsService.clean();
  }

  goBack() {
    if (this.isEdit) {
      this.router.navigate(['../../../agent-groups'], {relativeTo: this.route});
    } else {
      this.router.navigate(['../../agent-groups'], {relativeTo: this.route});
    }
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
        this.updateMatchingAgents();
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
      this.updateMatchingAgents();
    }
  }

  // query agent group matches
  updateTagMatches() {
    const payload = this.wrapPayload(true);
    // just validate and get matches summary
    this.agentGroupsService.validateAgentGroup(payload).subscribe((resp: any) => {
      this.tagMatch = {
        total: resp.body.matching_agents.total,
        online: resp.body.matching_agents.online,
      };
    });
  }

  updateMatchingAgents() {
    const tags = this.secondFormGroup.controls.tags.value;
    this.agentsService.getMatchingAgents(tags).subscribe(
      resp => {
        this.matchingAgents = resp;
      },
    );
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
    if (this.isEdit) {
      this.agentGroupsService.editAgentGroup({...payload, id: this.agentGroupID}).subscribe(resp => this.goBack());
    } else {
      this.agentGroupsService.addAgentGroup(payload).subscribe(() => this.goBack());
    }
  }

}
