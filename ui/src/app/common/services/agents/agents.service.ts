import {HttpClient, HttpParams} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import 'rxjs/add/observable/empty';

import {environment} from '../../../../environments/environment';
import {NotificationsService} from '../../../common/services/notifications/notifications.service';
import {Agent} from 'app/common/interfaces/orb/agent.interface';
import {NgxDatabalePageInfo, OrbPagination} from 'app/common/interfaces/orb/pagination';
import {AgentGroup} from 'app/common/interfaces/orb/agent.group.interface';

// default filters
const defLimit: number = 20;
const defOrder: string = 'name';
const defDir = 'desc';

@Injectable()
export class AgentsService {
  picture = 'assets/images/mainflux-logo.png';

  paginationCache: any = {};
  cache: OrbPagination<AgentGroup>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
    this.clean();
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
  }

  addAgentGroup(agentItem: any) {
    return this.http.post(environment.agentsUrl,
      agentItem,
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

  getAgentGroupById(agentId: string): any {
    return this.http.get(`${environment.agentsUrl}/${agentId}`)
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

  getAgentGroups(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    const offset = pageInfo.offset || this.cache.offset;
    let params = new HttpParams()
      .set('offset', offset.toString())
      .set('limit', (pageInfo.limit || this.cache.limit).toString())
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

    if (this.paginationCache[pageInfo.offset]) {
      return Observable.of(this.cache);
    }

    return this.http.get(environment.agentsUrl, {params})
      .map(
        (resp: any) => {
          this.paginationCache[pageInfo.offset] = true;
          // This is the position to insert the new data
          const start = pageInfo.offset * resp.limit;
          const newData = [...this.cache.data];
          newData.splice(start, resp.limit, ...resp.agents);
          this.cache = {
            ...this.cache,
            total: resp.total,
            data: newData,
          };
          if (pageInfo.name) this.cache.name = pageInfo.name;
          if (pageInfo.tags) this.cache.tags = pageInfo.tags;
          return this.cache;
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
}
