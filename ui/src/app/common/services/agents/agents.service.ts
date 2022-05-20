import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { delay, expand, reduce } from 'rxjs/operators';

// default filters
const defLimit: number = 100;
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
                                          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`);
          return Observable.throwError(err);
        },
      );
  }

  resetAgent(id: string) {
    return this.http.post(`${environment.agentsUrl}/${id}/rpc/reset`, {}, { observe: 'response' })
      .catch(err => {
        this.notificationsService.error('Failed to reset Agent',
                                        `Error: ${err.status} - ${err.statusText} - ${err.error.error}`);
        return Observable.throwError(err);
      });
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
                                          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentById(id: string): Observable<Agent> {
    return this.http.get<Agent>(`${environment.agentsUrl}/${id}`)
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

  editAgent(agent: Agent): any {
    return this.http.put(`${environment.agentsUrl}/${agent.id}`, agent)
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

  deleteAgent(agentId: string) {
    return this.http.delete(`${environment.agentsUrl}/${agentId}`)
      .map(
        resp => {
          this.cache.data.splice(this.cache.data.map(a => a.id).indexOf(agentId), 1);
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Delete Agent',
                                          `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  getMatchingAgents(tagsInfo: any): Observable<Agent[]> {
    const params = new HttpParams()
      .set('offset', AgentsService.getDefaultPagination().offset.toString())
      .set('limit', AgentsService.getDefaultPagination().limit.toString())
      .set('order', AgentsService.getDefaultPagination().order.toString())
      .set('dir', AgentsService.getDefaultPagination().dir.toString())
      .set('tags', JSON.stringify(tagsInfo).replace('[', '').replace(']', ''));

    return this.http.get(environment.agentsUrl, { params })
      .map(
        (resp: any) => {
          return resp.agents;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Matching Agents',
                                          `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  getAllAgents() {
    this.clean();
    const pageInfo = AgentsService.getDefaultPagination();


    return this.getAgents(pageInfo)
      .pipe(
        expand(data => {
          return data.next ? this.getAgents(data.next) : Observable.empty();
        }),
        delay(250),
        reduce<OrbPagination<Agent>>((acc, value) => {
          acc.data = value.data;
          acc.offset = 0;
          acc.total = acc.data.length;
          return acc;
        }, this.cache),
      );
  }

  getAgents(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    let limit = pageInfo?.limit || this.cache.limit;
    let order = pageInfo?.order || this.cache.order;
    let dir = pageInfo?.dir || this.cache.dir;
    let offset = pageInfo?.offset || 0;
    let doClean = false;
    let params = new HttpParams();

    if (isFilter) {
      if (pageInfo?.name) {
        params = params.set('name', pageInfo.name);
        // is filter different than last filter?
        doClean = !this.paginationCache?.name || this.paginationCache?.name !== pageInfo.name;
      }
      // was filtered, no longer
    } else if (this.paginationCache?.isFilter === true) {
      doClean = true;
    }

    if (pageInfo.order !== this.cache.order || pageInfo.dir !== this.cache.dir) {
      doClean = true;
    }

    if (doClean) {
      this.clean();
      offset = 0;
      limit = this.cache.limit = pageInfo.limit;
      dir = pageInfo.dir;
      order = pageInfo.order;
    }

    if (this.paginationCache[offset]) {
      return of(this.cache);
    }
    params = params
      .set('offset', (offset).toString())
      .set('limit', limit.toString())
      .set('order', order)
      .set('dir', dir);

    return this.http.get(environment.agentsUrl, { params })
      .map(
        (resp: any) => {
          this.paginationCache[pageInfo?.offset / pageInfo?.limit || 0] = true;

          // This is the position to insert the new data
          const start = pageInfo?.offset;

          const newData = [...this.cache.data];

          newData.splice(start, resp.limit,
                         ...resp.agents.map(agent => {
                           agent.combined_tags = { ...agent?.orb_tags, ...agent?.agent_tags };
                           return agent;
                         }));

          this.cache = {
            ...this.cache,
            next: resp.offset + resp.limit < resp.total && {
              limit: resp.limit,
              offset: (parseInt(resp.offset, 10) + parseInt(resp.limit, 10)).toString(),
              order: 'name',
              dir: 'desc',
            },
            limit: resp.limit,
            offset: resp.offset,
            dir: resp.direction,
            order: resp.order,
            total: resp.total,
            data: newData,
            name: pageInfo?.name,
          };

          return this.cache;
        })
      .catch(
        err => {
          this.notificationsService.error('Failed to get Agents',
                                          `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }
}
