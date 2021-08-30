import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination';
import { Agent } from 'app/common/interfaces/orb/agent.interface';

// default filters
const defLimit: number = 20;
const defOrder: string = 'name';
const defDir = 'desc';

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

   getMatchingAgents(tagsInfo: any) {
    const params = new HttpParams()
      .set('offset', AgentsService.getDefaultPagination().offset.toString())
      .set('limit', AgentsService.getDefaultPagination().limit.toString())
      .set('order', AgentsService.getDefaultPagination().order.toString())
      .set('dir', AgentsService.getDefaultPagination().dir.toString())
      .set('tags', JSON.stringify(tagsInfo).replace('[', '').replace(']', ''));

    return this.http.get(environment.agentsUrl, {params})
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


  getAgents(pageInfo: NgxDatabalePageInfo, isFilter = false) {
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
}
