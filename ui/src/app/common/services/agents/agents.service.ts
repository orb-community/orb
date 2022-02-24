import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';

// default filters
const defLimit: number = 20;
const defOrder: string = 'name';
const defDir = 'desc';

export enum AvailableOS {
  DOCKER = 'docker',
}

@Injectable()
export class AgentsService {
  paginationCache: any = {};

  cache: OrbPagination<Agent>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
    this.clean();
  }

  public static getDefaultPagination(): OrbPagination<Agent> {
    return {
      limit: defLimit,
      order: defOrder,
      dir: defDir,
      offset: 0,
      total: 0,
      data: null,
    };
  }

  clean() {
    this.cache = {
      limit: defLimit,
      offset: 0,
      order: defOrder,
      total: 0,
      dir: defDir,
      data: [],
    };
    this.paginationCache = {};
  }

  addAgent(agentItem: Agent) {
    return this.http.post(environment.agentsUrl,
        { ...agentItem, validate_only: false },
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Agent',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  validateAgent(agentItem: Agent) {
    return this.http.post(environment.validateAgentsUrl,
        { ...agentItem, validate_only: true },
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Validate Agent',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentById(id: string): Observable<Agent> {
    return this.http.get<Agent>(`${ environment.agentsUrl }/${ id }`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Agent',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  editAgent(agent: Agent): any {
    return this.http.put(`${ environment.agentsUrl }/${ agent.id }`, agent)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Agent',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  deleteAgent(agentId: string) {
    return this.http.delete(`${ environment.agentsUrl }/${ agentId }`)
      .map(
        resp => {
          this.cache.data.splice(this.cache.data.map(a => a.id).indexOf(agentId), 1);
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Delete Agent',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getMatchingAgents(tagsInfo: any) {
    const params = new HttpParams()
      .set('offset', AgentsService.getDefaultPagination().offset.toString())
      .set('limit', AgentsService.getDefaultPagination().limit.toString())
      .set('order', AgentsService.getDefaultPagination().order.toString())
      .set('dir', AgentsService.getDefaultPagination().dir.toString())
      .set('tags', JSON.stringify(tagsInfo).replace('[', '').replace(']', ''));

    return this.http.get(environment.agentsUrl, { params })
      .map(
        (resp: any) => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Matching Agents',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getAgents(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    const { limit, offset, name, tags } = pageInfo || this.cache;
    let params = new HttpParams()
      .set('offset', (offset * limit).toString())
      .set('limit', limit.toString())
      .set('order', this.cache.order)
      .set('dir', this.cache.dir);

    if (isFilter) {
      if (name) {
        params = params.append('name', name);
      }
      if (tags) {
        params.append('tags', JSON.stringify(tags));
      }
      this.paginationCache[offset] = false;
    }

    if (this.paginationCache[offset]) {
      return of(this.cache);
    }

    return this.http.get(environment.agentsUrl, { params })
      .map(
        (resp: any) => {
          this.paginationCache[offset || 0] = true;
          // This is the position to insert the new data
          const start = resp.offset;
          const newData = [...this.cache.data];
          newData.splice(start, resp.limit, ...resp.agents);
          this.cache = {
            ...this.cache,
            offset: Math.floor(resp.offset / resp.limit),
            total: resp.total,
            data: newData,
          };
          if (name) this.cache.name = name;
          if (tags) this.cache.tags = tags;
          return this.cache;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Agents',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }
}
