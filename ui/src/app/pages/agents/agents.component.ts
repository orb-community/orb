import {AfterViewInit, ChangeDetectorRef, Component, TemplateRef, ViewChild} from '@angular/core';
import {NbDialogService} from '@nebular/theme';

import {DropdownFilterItem, PageFilters, TablePage, User} from 'app/common/interfaces/mainflux.interface';
import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {ActivatedRoute, Router} from '@angular/router';
import {STRINGS} from 'assets/text/strings';
import {AgentDeleteComponent} from 'app/pages/agents/delete/agent.delete.component';
import {AgentDetailsComponent} from 'app/pages/agents/details/agent.details.component';
import {ColumnMode, TableColumn} from '@swimlane/ngx-datatable';
import {AgentsService} from 'app/common/services/agents/agents.service';
import {debounceTime, distinctUntilChanged} from "rxjs/operators";

const defFreq: number = 100;

@Component({
  selector: 'ngx-agents-component',
  templateUrl: './agents.component.html',
  styleUrls: ['./agents.component.scss'],
})
export class AgentsComponent implements AfterViewInit {
  strings = STRINGS.agents;

  columnMode = ColumnMode;
  columns: TableColumn[];

  // templates

  @ViewChild('agentsTemplateCell') agentsTemplateCell: TemplateRef<any>;
  @ViewChild('agentTagsTemplateCell') agentTagsTemplateCell: TemplateRef<any>;
  @ViewChild('addAgentTemplateRef') addAgentTemplateRef: TemplateRef<any>;
  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  //input
  @ViewChild('input') searchInput;

  page: TablePage = {
    limit: 10,
  };

  pageFilters: PageFilters = {
    offset: 0,
    order: 'id',
    dir: 'desc',
    name: '',
  };

  tableFilters: DropdownFilterItem[];

  searchFreq = 0;

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private agentsService: AgentsService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        width: 120,
        maxWidth: 243,
      },
      {
        name: 'Description',
        resizeable: false,
        width: 200,
        maxWidth: 350,
      },
      {
        prop: 'agents',
        name: 'Agents',
        resizeable: false,
        width: 80,
        maxWidth: 100,
        cellTemplate: this.agentsTemplateCell,
      },
      {
        name: 'Tags',
        width: 200,
        canAutoResize: true,
        cellTemplate: this.agentTagsTemplateCell,
      },
      {
        name: '',
        prop: 'actions',
        width: 120,
        resizeable: false,
        sortable: false,
        cellTemplate: this.actionsTemplateCell,
      },
    ];
    this.tableFilters = this.columns.map((entry, index) => ({
      id: index.toString(),
      name: entry.name,
      order: 'asc',
      selected: false,
    })).filter((filter) => (!filter.name?.startsWith('orb-')));
    debugger;
    this.searchInput.update
      .pipe(debounceTime(500))
      .pipe(distinctUntilChanged())
      .subscribe(model => (value) => {
        debugger;
        this.getAgents(value);
      });

    this.getAgents();
    this.cdr.detectChanges();
  }

  getAgents(name?: string): void {
    this.pageFilters.name = name;
    this.agentsService.getAgentGroups(this.pageFilters).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.agents,
        };
      },
    );
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      this.pageFilters.offset = this.page.offset - this.page.limit;
    }
    if (dir === 'next') {
      this.pageFilters.offset = this.page.offset + this.page.limit;
    }
    this.getAgents();
  }

  onChangeLimit(limit: number) {
    this.pageFilters.limit = limit;
    this.getAgents();
  }

  onOpenAdd() {
    this.router.navigate(['../agents/add'], {
      relativeTo: this.route,
    });
  }

  onOpenEdit(row: any) {
    this.router.navigate(['../agents/edit'], {
      relativeTo: this.route,
      queryParams: {id: row.id},
      state: {agent: row},
    });
  }

  openDeleteModal(row: any) {
    const {name, id} = row;
    this.dialogService.open(AgentDeleteComponent, {
      context: {agent: {name, id}},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.agentsService.deleteAgentGroup(row.id).subscribe(
            () => {
              this.page.rows = this.page.rows.filter((u: User) => u.id !== row.id);
              this.notificationsService.success('Sink Item successfully deleted', '');
            },
          );
        }
      },
    );
  }

  openDetailsModal(row: any) {
    const {name, description, backend, config, ts_created, id} = row;

    this.dialogService.open(AgentDetailsComponent, {
      context: {agent: {id, name, description, backend, config, ts_created}},
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getAgents();
        }
      },
    );
  }

  searchAgentByName(input) {
    const t = new Date().getTime();
    if ((t - this.searchFreq) > defFreq) {
      this.getAgents(input);
      this.searchFreq = t;
    }
  }

  filterByActive = (agent) => agent.status === 'active';

  mockCreate() {
    for (let i = 0; i < 10; i++) {
      this.agentsService.addAgentGroup({
        name: `sample-at-${Math.floor(Math.random() * 10000)}`,
        description: 'Lorem ipsum ipsils',
        tags: {
          node_type: 'dns',
          region: 'EU',
        },
      }).subscribe(evt => {
        console.log('added');
        this.getAgents();
      });
    }
  }
}
