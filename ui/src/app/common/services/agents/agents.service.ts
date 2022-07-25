import {HttpClient, HttpParams} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {EMPTY, Observable} from 'rxjs';
import 'rxjs/add/observable/empty';

import {Agent, AgentPolicyAggStates} from 'app/common/interfaces/orb/agent.interface';
import {OrbPagination} from 'app/common/interfaces/orb/pagination.interface';
import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {environment} from 'environments/environment';
import {expand, map, scan, takeWhile} from 'rxjs/operators';
import {AgentPolicyState, AgentPolicyStates} from 'app/common/interfaces/orb/agent.policy.interface';

export enum AvailableOS {
  DOCKER = 'docker',
}

@Injectable()
export class AgentsService {
  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {}

  addAgent(agentItem: Agent) {
    return this.http
      .post<Agent>(
        environment.agentsUrl,
        { ...agentItem, validate_only: false },
        { observe: 'response' },
      )
      .map((resp) => {
        let { body: agent } = resp;
        agent = {
          ...agent,
          combined_tags: { ...agent?.orb_tags, ...agent?.agent_tags },
        };
        return agent;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to create Agent',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  resetAgent(id: string) {
    return this.http
      .post(
        `${environment.agentsUrl}/${id}/rpc/reset`,
        {},
        { observe: 'response' },
      )
      .catch((err) => {
        this.notificationsService.error(
          'Failed to reset Agent',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  validateAgent(agentItem: Agent) {
    return this.http
      .post(
        environment.validateAgentsUrl,
        { ...agentItem, validate_only: true },
        { observe: 'response' },
      )
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Validate Agent',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  getAgentById(id: string): Observable<Agent> {
    return this.http
      .get<Agent>(`${environment.agentsUrl}/${id}`)
      .map((agent) => {
        return {
          ...agent,
          combined_tags: { ...agent?.orb_tags, ...agent?.agent_tags },
        };
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to fetch Agent',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  editAgent(agent: Agent): any {
    return this.http
      .put<Agent>(`${environment.agentsUrl}/${agent.id}`, agent)
      .map((resp) => {
        return {
          ...resp,
          combined_tags: { ...resp?.orb_tags, ...resp?.agent_tags },
        };
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to edit Agent',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  deleteAgent(agentId: string) {
    return this.http
      .delete(`${environment.agentsUrl}/${agentId}`)
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Delete Agent',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getAllAgents(tags?: any) {
    const page = {
      order: 'name',
      dir: 'asc',
      limit: 100,
      data: [],
      offset: 0,
      tags,
    } as OrbPagination<Agent>;
    return this.getAgents(page).pipe(
      expand((data) => {
        return data.next ? this.getAgents(data.next) : EMPTY;
      }),
      takeWhile((data) => data.next !== undefined),
      map((_page) => _page.data),
      scan((acc, v) => [...acc, ...v]),
    );
  }

  getAgents(page: OrbPagination<Agent>) {
    let params = new HttpParams()
      .set('order', page.order)
      .set('dir', page.dir)
      .set('offset', page.offset.toString())
      .set('limit', page.limit.toString());

    if (page.tags) {
      params = params.set(
        'tags',
        JSON.stringify(page.tags).replace('[', '').replace(']', ''),
      );
    }

    return this.http
      .get(`${environment.agentsUrl}`, { params })
      .pipe(
        map((resp: any) => {
          const {
            order,
            direction: dir,
            offset,
            limit,
            total,
            agents,
            tags,
          } = resp;
          const next = offset + limit < total && {
            limit,
            order,
            dir,
            tags,
            offset: (parseInt(offset, 10) + parseInt(limit, 10)).toString(),
          };
          const data = this.mapUIAggregates(agents);
          return {
            order,
            dir,
            offset,
            limit,
            total,
            data,
            next,
          } as OrbPagination<Agent>;
        }),
      )
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Agents',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  mapUIAggregates(agents) {
    return agents.map((agent) => {
      // combined tags helper
      agent.combined_tags = { ...agent?.orb_tags, ...agent?.agent_tags };
      // map agg policy state
      const {agg_info, agg_state} = this.policyAggState(agent);
      agent.policy_agg_info = agg_info;
      agent.policy_agg_state = agg_state;
      return agent;
    });
  }

  policyAggState(agent) {
    const { policy_state } = agent;
    let agg_info = 'No Policies Applied';
    let agg_state = AgentPolicyAggStates.none;

    const policies = !!policy_state && Object.values(policy_state) as AgentPolicyState[] || [];
    if (policies.length > 0) {
      let err = 0;
      policies.reduce((prev, curr) => {
        if (curr.state !== AgentPolicyStates.running) {
          err++;
        }
        return curr;
      });
      if (err > 0) {
        agg_info = `${err} out of ${policies.length} policies are not running`;
        agg_state = AgentPolicyAggStates.failure;
      } else {
        agg_info = `All policies are running`;
        agg_state = AgentPolicyAggStates.healthy;
      }
    }

    return { agg_info, agg_state};
  }
}
