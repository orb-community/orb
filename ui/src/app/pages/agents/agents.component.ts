import { Component, OnInit } from '@angular/core';
import { NbDialogService } from '@nebular/theme';

import {
  DropdownFilterItem,
  PageFilters,
  TableConfig,
  TablePage,
  User,
} from 'app/common/interfaces/mainflux.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { STRINGS } from 'assets/text/strings';
import { AgentsMockService } from 'app/common/services/agents/agents.mock.service';
import { AgentDeleteComponent } from 'app/pages/agents/delete/agent.delete.component';
import { AgentDetailsComponent } from 'app/pages/agents/details/agent.details.component';

const defFreq: number = 100;

/**
 * Available sink statuses
 */
export enum sinkStatus {
  active = 'active',
  error = 'error',
}

export enum sinkTypesList {
  prometheus = 'prometheus',
  // aws = 'aws',
  // s3 = 's3',
  // azure = 'azure',
}

@Component({
  selector: 'ngx-agents-component',
  templateUrl: './agents.component.html',
  styleUrls: ['./agents.component.scss'],
})
export class AgentsComponent implements OnInit {
  strings = STRINGS.sink;

  tableConfig: TableConfig = {
    colNames: ['Name', 'Description', 'Type', 'Status', 'Tags', 'orb-sink-add'],
    keys: ['name', 'description', 'type', 'status', 'tags', 'orb-action-hover'],
  };

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
    private dialogService: NbDialogService,
    private agentsService: AgentsMockService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.tableFilters = this.tableConfig.colNames.map((name, index) => ({
      id: index.toString(),
      name,
      order: 'asc',
      selected: false,
    })).filter((filter) => (!filter.name.startsWith('orb-')));
  }

  ngOnInit() {
    // Fetch all sinks
    this.getAgents();
  }

  getAgents(name?: string): void {
    this.pageFilters.name = name;
    this.agentsService.getAgents(this.pageFilters).subscribe(
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
          this.agentsService.deleteAgent(row.id).subscribe(
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

  filterByInactive = (sink) => sink.status === 'inactive';

}
