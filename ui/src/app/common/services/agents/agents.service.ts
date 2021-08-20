import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { PageFilters } from 'app/common/interfaces/mainflux.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';

// default filters
const defLimit: number = 20;
const defOrder: string = 'id';
const defDir: string = 'desc';

@Injectable()
export class AgentsService {
  picture = 'assets/images/mainflux-logo.png';

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
  }

  addAgentGroup(agentGroupItem: AgentGroup) {
    return this.http.post(environment.agentsUrl,
      agentGroupItem,
      {observe: 'response'})
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Agent',
            `Error: ${err.status} - ${err.statusText} - ${err.error.error}`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentGroupById(id: string): any {
    return this.http.get(`${environment.agentsUrl}/${id}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Agent',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentGroups(filters: PageFilters) {
    filters.offset = filters.offset || 0;
    filters.limit = filters.limit || defLimit;
    filters.order = filters.order || defOrder;
    filters.dir = filters.dir || defDir;

    let params = new HttpParams()
      .set('offset', filters.offset.toString())
      .set('limit', filters.limit.toString())
      .set('order', filters.order)
      .set('dir', 'asc');

    if (filters.name) {
      params = params.append('name', filters.name);
    }

    return this.http.get(environment.agentsUrl, {params})
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Agents',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  editAgentGroup(agentItem: Agent): any {
    return this.http.put(`${environment.agentsUrl}/${agentItem.id}`, agentItem)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Agent',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  deleteAgentGroup(agentId: string) {
    return this.http.delete(`${environment.agentsUrl}/${agentId}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to delete Agent',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }
}
