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
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { catchError, expand, reduce } from 'rxjs/operators';

// default filters
const defLimit: number = 100;
const defOrder: string = 'name';
const defDir = 'asc';

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

  addAgentPolicy(agentPolicyItem: AgentPolicy): Observable<AgentPolicy> {
    return this.http
      .post(
        environment.agentPoliciesUrl,
        { ...agentPolicyItem },
        { observe: 'response' },
      )
      .map((resp) => {
        return resp.body as AgentPolicy;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to create Agent Policy',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return of(err);
      });
  }

  duplicateAgentPolicy(id: string): any {
    return this.http
      .post(
        `${environment.agentPoliciesUrl}/${id}/duplicate`,
        {},
        { observe: 'response' },
      )
      .map((resp) => {
        const { body } = resp;
        return body;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to duplicate Agent Policy',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return of(err);
      });
  }

  getAgentPolicyById(id: string): Observable<AgentPolicy> {
    return this.http.get(`${environment.agentPoliciesUrl}/${id}`).pipe(
      catchError((err) => {
        this.notificationsService.error(
          'Failed to fetch Agent Policy',
          `Error: ${err.status} - ${err.statusText}`,
        );
        err['id'] = id;
        return of(err);
      }),
    );
  }

  editAgentPolicy(agentPolicy: AgentPolicy): any {
    return this.http
      .put(`${environment.agentPoliciesUrl}/${agentPolicy.id}`, agentPolicy)
      .catch((err) => {
        this.notificationsService.error(
          'Failed to edit Agent Policy',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return of(err);
      });
  }

  deleteAgentPolicy(agentPoliciesId: string) {
    return this.http
      .delete(`${environment.agentPoliciesUrl}/${agentPoliciesId}`)
      .map((resp) => {
        this.cache.data.splice(
          this.cache.data.map((ap) => ap.id).indexOf(agentPoliciesId),
          1,
        );
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Delete Agent Policies',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getAllAgentPolicies() {
    this.clean();
    const pageInfo = AgentPoliciesService.getDefaultPagination();

    return this.getAgentsPolicies(pageInfo).pipe(
      expand((data) => {
        return data.next
          ? this.getAgentsPolicies(data.next)
          : Observable.empty();
      }),
      reduce<OrbPagination<AgentPolicy>>((acc, value) => {
        acc.data = value.data;
        acc.offset = 0;
        acc.total = acc.data.length;
        return acc;
      }, this.cache),
    );
  }

  getAgentsPolicies(pageInfo: NgxDatabalePageInfo, isFilter = false) {
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
      .get(environment.agentPoliciesUrl, { params })
      .map((resp: any) => {
        this.paginationCache[pageInfo?.offset / pageInfo?.limit || 0] = true;

        // This is the position to insert the new data
        const start = pageInfo?.offset;

        const newData = [...this.cache.data];

        newData.splice(start, resp.limit, ...resp.data);

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
          'Failed to get Agent Policies',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getAvailableBackends() {
    return this.http
      .get(environment.agentsBackendUrl)
      .map((resp: any) => {
        return resp.backends;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Available Backends',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getBackendConfig(route: string[]) {
    const final = route.join('/');

    return this.http
      .get(`${environment.agentsBackendUrl}/${final}`)
      .map((response: any) => {
        return response;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Available Backends',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }
}
