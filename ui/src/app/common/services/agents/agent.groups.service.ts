import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import {
  NgxDatabalePageInfo,
  OrbPagination,
} from 'app/common/interfaces/orb/pagination.interface';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { catchError, expand, reduce } from 'rxjs/operators';

// default filters
const defLimit: number = 100;
const defOrder: string = 'name';
const defDir = 'asc';

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
    return this.http
      .post(
        environment.agentGroupsUrl,
        {
          ...agentGroupItem,
          validate_only: false,
        },
        { observe: 'response' },
      )
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to create Agent Group',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  validateAgentGroup(agentGroupItem: AgentGroup) {
    return this.http
      .post(
        environment.validateAgentGroupsUrl,
        {
          ...agentGroupItem,
          validate_only: true,
        },
        { observe: 'response' },
      )
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Validate Agent Group',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  getAgentGroupById(id: string): Observable<AgentGroup> {
    return this.http.get(`${environment.agentGroupsUrl}/${id}`).pipe(
      catchError((err) => {
        this.notificationsService.error(
          'Failed to fetch Agent Group',
          `Error: ${err.status} - ${err.statusText}`,
        );
        err['id'] = id;
        return of(err);
      }),
    );
  }

  getAllAgentGroups() {
    this.clean();
    const pageInfo = AgentGroupsService.getDefaultPagination();

    return this.getAgentGroups(pageInfo).pipe(
      expand((data) => {
        return data.next ? this.getAgentGroups(data.next) : Observable.empty();
      }),
      reduce<OrbPagination<AgentGroup>>((acc, value) => {
        acc.data = value.data;
        acc.offset = 0;
        acc.total = acc.data.length;
        return acc;
      }, this.cache),
    );
  }

  getAgentGroups(pageInfo: NgxDatabalePageInfo, isFilter = false) {
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
        doClean =
          !this.paginationCache?.name ||
          this.paginationCache?.name !== pageInfo.name;
      }
      // was filtered, no longer
    } else if (this.paginationCache?.isFilter === true) {
      doClean = true;
    }

    if (
      pageInfo.order !== this.cache.order ||
      pageInfo.dir !== this.cache.dir
    ) {
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
      .set('offset', offset.toString())
      .set('limit', limit.toString())
      .set('order', order)
      .set('dir', dir);

    return this.http
      .get(environment.agentGroupsUrl, { params })
      .map((resp: any) => {
        this.paginationCache[pageInfo?.offset || 0] = true;

        // This is the position to insert the new data
        const start = pageInfo?.offset || 0;

        const newData = [...this.cache.data];

        newData.splice(start, resp.limit, ...resp.agentGroups);

        this.cache = {
          ...this.cache,
          next: resp.offset + resp.limit < resp.total && {
            limit: resp.limit,
            offset: (
              parseInt(resp.offset, 10) + parseInt(resp.limit, 10)
            ).toString(),
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
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Agent Groups',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  editAgentGroup(agentGroup: AgentGroup): any {
    return this.http
      .put(`${environment.agentGroupsUrl}/${agentGroup.id}`, agentGroup)
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to edit Agent Group',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  deleteAgentGroup(agentGroupId: string) {
    return this.http
      .delete(`${environment.agentGroupsUrl}/${agentGroupId}`)
      .map((resp) => {
        this.cache.data.splice(
          this.cache.data.map((ag) => ag.id).indexOf(agentGroupId),
          1,
        );
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Delete Agent Group',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }
}
