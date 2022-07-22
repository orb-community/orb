import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import 'rxjs/add/observable/empty';

import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { environment } from 'environments/environment';
import { catchError, expand, map, scan, takeWhile } from 'rxjs/operators';

@Injectable()
export class AgentPoliciesService {
  backendsCache: OrbPagination<{ [propName: string]: any }>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {}

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
        return throwError(err);
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
        return throwError(err);
      });
  }

  deleteAgentPolicy(agentPoliciesId: string) {
    return this.http
      .delete(`${environment.agentPoliciesUrl}/${agentPoliciesId}`)
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Delete Agent Policies',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getAllAgentPolicies() {
    const page = {
      order: 'name',
      dir: 'asc',
      limit: 100,
      data: [],
      offset: 0,
    } as OrbPagination<AgentPolicy>;

    return this.getAgentsPolicies(page).pipe(
      expand((data) => {
        return data.next
          ? this.getAgentsPolicies(data.next)
          : Observable.empty();
      }),
      takeWhile((data) => data.next !== undefined),
      map((_page) => _page.data),
      scan((acc, v) => [...acc, ...v]),
    );
  }

  getAgentsPolicies(page: OrbPagination<AgentPolicy>) {
    const params = new HttpParams()
      .set('order', page.order)
      .set('dir', page.dir)
      .set('offset', page.offset.toString())
      .set('limit', page.limit.toString());

    return this.http
      .get(environment.agentPoliciesUrl, { params })
      .pipe(
        map((resp: any) => {
          const { order, direction, offset, limit, total, data, tags } = resp;
          const next = offset + limit < total && {
            limit,
            order,
            dir: direction,
            tags,
            offset: (parseInt(offset, 10) + parseInt(limit, 10)).toString(),
          };
          return {
            order,
            dir: direction,
            offset,
            limit,
            total,
            data,
            next,
          } as OrbPagination<AgentPolicy>;
        }),
      )
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
