import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { catchError, delay, expand, reduce } from 'rxjs/operators';

// default filters
const defLimit: number = 20;
const defOrder: string = 'name';
const defDir = 'desc';

@Injectable()
export class AgentGroupsService {
  paginationCache: any = {};

  cache: OrbPagination<AgentGroup>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
    this.clean();
  }

  public static getDefaultPagination(): OrbPagination<AgentGroup> {
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

  addAgentGroup(agentGroupItem: AgentGroup) {
    return this.http.post(environment.agentGroupsUrl,
        { ...agentGroupItem, validate_only: false },
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Agent Group',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  validateAgentGroup(agentGroupItem: AgentGroup) {
    return this.http.post(environment.validateAgentGroupsUrl,
        { ...agentGroupItem, validate_only: true },
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Validate Agent Group',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  getAgentGroupById(id: string): Observable<AgentGroup> {
    return this.http.get(`${ environment.agentGroupsUrl }/${ id }`)
      .pipe(
        catchError(err => {
          this.notificationsService.error('Failed to fetch Agent Group',
            `Error: ${ err.status } - ${ err.statusText }`);
          return of(err);
        }),
      );
  }

  getAllAgentGroups() {
    const pageInfo = AgentGroupsService.getDefaultPagination();
    pageInfo.limit = 100;

    return this.getAgentGroups(pageInfo)
      .pipe(
        expand(data => {
          return data.next ? this.getAgentGroups(data.next) : Observable.empty();
        }),
        delay(250),
        reduce<OrbPagination<AgentGroup>>((acc, value) => {
          acc.data.splice(value.offset, value.limit, ...value.data);
          acc.offset = 0;
          return acc;
        }, this.cache),
      );
  }

  getAgentGroups(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    const offset = pageInfo?.offset || this.cache.offset;
    const limit = pageInfo?.limit || this.cache.limit;
    let params = new HttpParams()
      .set('offset', (offset * limit).toString())
      .set('limit', limit.toString())
      .set('order', this.cache.order)
      .set('dir', this.cache.dir);

    if (isFilter) {
      if (pageInfo?.name) {
        params = params.append('name', pageInfo.name);
      }
      if (pageInfo?.tags) {
        params.append('tags', JSON.stringify(pageInfo.tags));
      }
      this.paginationCache[offset] = false;
    }

    if (this.paginationCache[pageInfo?.offset]) {
      return of(this.cache);
    }

    return this.http.get(environment.agentGroupsUrl, { params })
      .map(
        (resp: any) => {
          this.paginationCache[pageInfo?.offset || 0] = true;
          // This is the position to insert the new data
          const start = pageInfo?.offset || 0;
          const newData = [...this.cache.data];
          newData.splice(start, resp.limit, ...resp.agentGroups);
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
          this.notificationsService.error('Failed to get Agent Groups',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  editAgentGroup(agentGroup: AgentGroup): any {
    return this.http.put(`${ environment.agentGroupsUrl }/${ agentGroup.id }`, agentGroup)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Agent Group',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  deleteAgentGroup(agentGroupId: string) {
    return this.http.delete(`${ environment.agentGroupsUrl }/${ agentGroupId }`)
      .map(
        resp => {
          this.cache.data.splice(this.cache.data.map(ag => ag.id).indexOf(agentGroupId), 1);
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Delete Agent Group',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }
}
