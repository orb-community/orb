import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

// default filters
const defLimit: number = 20;
const defOrder: string = 'name';
const defDir = 'desc';

@Injectable()
export class AgentPoliciesService {
  paginationCache: any = {};

  cache: OrbPagination<AgentPolicy>;

  backendsCache: OrbPagination<{ [propName: string]: any }>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
    this.clean();
  }

  public static getDefaultPagination(): OrbPagination<AgentPolicy> {
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

  addAgentPolicy(agentPolicyItem: AgentPolicy) {
    return this.http.post(environment.agentPoliciesUrl,
        { ...agentPolicyItem },
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Agent Policy',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentPolicyById(id: string): any {
    return this.http.get(`${ environment.agentPoliciesUrl }/${ id }`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Agent Policy',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  editAgentPolicy(agentPolicy: AgentPolicy): any {
    return this.http.put(`${ environment.agentPoliciesUrl }/${ agentPolicy.id }`, agentPolicy)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Agent Policy',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  deleteAgentPolicy(agentPoliciesId: string) {
    return this.http.delete(`${ environment.agentPoliciesUrl }/${ agentPoliciesId }`)
      .map(
        resp => {
          this.cache.data.splice(this.cache.data.map(ap => ap.id).indexOf(agentPoliciesId), 1);
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Delete Agent Policies',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentsPolicies(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    const offset = pageInfo?.offset || this.cache.offset;
    const limit = pageInfo?.limit || this.cache.limit;
    let params = new HttpParams()
      .set('offset', (offset * limit).toString())
      .set('limit', limit.toString())
      .set('order', this.cache.order)
      .set('dir', this.cache.dir);

    if (isFilter) {
      if (pageInfo.name) {
        params = params.append('name', pageInfo.name);
      }
      if (pageInfo.tags) {
        params.append('tags', JSON.stringify(pageInfo.tags));
      }
      this.paginationCache[offset] = false;
    }

    if (this.paginationCache[pageInfo?.offset]) {
      return of(this.cache);
    }

    return this.http.get(environment.agentPoliciesUrl, { params })
      .map(
        (resp: any) => {
          this.paginationCache[pageInfo?.offset] = true;
          // This is the position to insert the new data
          const start = resp.offset;
          const newData = [...this.cache.data];
          // TODO figure out what field name for object data in response...
          newData.splice(start, resp.limit, ...resp.data);
          this.cache = {
            ...this.cache,
            offset: Math.floor(resp.offset / resp.limit),
            total: resp.total,
            data: newData,
          };
          if (pageInfo?.name) this.cache.name = pageInfo.name;
          if (pageInfo?.tags) this.cache.tags = pageInfo.tags;
          return this.cache;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Agent Policies',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getAvailableBackends() {
    return this.http.get(environment.agentsBackendUrl)
      .map((resp: any) => {
        return resp.backends;
      }).catch(err => {
        this.notificationsService.error('Failed to get Available Backends',
          `Error: ${ err.status } - ${ err.statusText }`);
        return Observable.throwError(err);
      });
  }

  getBackendConfig(route: string[]) {
    const final = route.join('/');

    return this.http.get(`${ environment.agentsBackendUrl }/${ final }`)
      .map((response: any) => {
        return response;
      }).catch(
        err => {
          this.notificationsService.error('Failed to get Available Backends',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );

  }
}
