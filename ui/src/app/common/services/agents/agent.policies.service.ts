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
        { ...agentPolicyItem, validate_only: false },
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
    const offset = pageInfo.offset || this.cache.offset;
    const limit = pageInfo.limit || this.cache.limit;
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

    if (this.paginationCache[pageInfo.offset]) {
      return of(this.cache);
    }

    return this.http.get(environment.agentPoliciesUrl, { params })
      .map(
        (resp: any) => {
          this.paginationCache[pageInfo.offset] = true;
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
          if (pageInfo.name) this.cache.name = pageInfo.name;
          if (pageInfo.tags) this.cache.tags = pageInfo.tags;
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
    // return this.http.get(environment.agentsBackendUrl)
    //   .map((resp: any) => {
    //     return resp.backend;
    //   }).catch(
    //     err => {
    //       this.notificationsService.error('Failed to get Availab//         `Error: ${ err.status } - ${
    // err.statusText }`); return Observable.throwError(err); }, ); TODO remove mock and uncomment http request
    return new Observable(subscriber => {
      // TODO continue this format in future
      const resp = {
        data: [
          {
            'backend': 'pktvisor',
            'description': 'pktvisor observability agent from pktvisor.dev',
            // todo I could use some meta like this
            // 'config': ['taps', 'input', 'handlers'],
          },
        ],
      };
      subscriber.next(resp);
    });
  }

  // todo from this point on I have to assume pktvisor hardcoded steps
  // tap -> which will have a predefined input
  // fill input config form, will be dynamic to some extent
  // from there on, select handlers
  // ${backend}/${config[i]}/ // pktvisor/[taps,inputs,handlers]
  getBackendConfig(route: string[]) {
    const final = route.join('/');
    // return this.http.get(`${environment.agentsBackendUrl}/${final})
    //   .map((resp: any) => {
    //     return resp.backend;
    //   }).catch(
    //     err => {
    //       this.notificationsService.error('Failed to get Available Backends',
    //         `Error: ${ err.status } - ${ err.statusText }`);
    //       return Observable.throwError(err);
    //     },
    //   );
    // TODO remove mock and uncomment http request
    if (final === 'pktvisor/taps') {
      return this.http.get(`${ environment.agentsBackendUrl }/${ final }`)
        .map((resp: any) => {
          return resp.backend;
        }).catch(
          err => {
            this.notificationsService.error('Failed to get Available Backends',
              `Error: ${ err.status } - ${ err.statusText }`);
            return Observable.throwError(err);
          },
        );
    } else {
      let resp;
      switch (final) {
        case 'pktvisor/taps':
          resp = {
            data: [
              {
                'name': 'ethernet',
                'input_type': 'pcap',
                'config_predefined': [
                  'iface',
                ],
                'agents': {
                  'total': 1,
                },
              },
            ],
          };
          break;
        case 'pktvisor/inputs':
          resp = {
            data: {
              'pcap': {
                'version': '1.0',
                'config': {
                  'iface': {
                    'required': true,
                    'type': 'string',
                    'input': 'text',
                    'title': 'Interface',
                    'name': 'iface',
                    'description': 'The ethernet interface to capture on',
                  },
                  'bpf': {
                    'required': false,
                    'type': 'string',
                    'input': 'text',
                    'title': 'Filter Expression',
                    'name': 'bpf',
                    'description': 'tcpdump compatible filter expression for limiting the traffic examined (with BPF). Example: "port 53"',
                  },
                  'host_spec': {
                    'required': false,
                    'type': 'string',
                    'input': 'text',
                    'title': 'Host Specification',
                    'name': 'host_spec',
                    'description': 'Subnets (comma separated) to consider this HOST, in CIDR form. Example: "10.0.1.0/24,10.0.2.1/32,2001:db8::/64"',
                  },
                  'pcap_source': {
                    'required': false,
                    'type': 'string',
                    'input': 'text',
                    'title': 'pcap Engine',
                    'name': 'pcap_source',
                    'description': 'pcap backend engine to use. Defaults to best for platform.',
                  },
                },
              },
              'dnstap': {
                'version': '1.0',
                'config': {
                  'type': {
                    'type': 'string',
                    'input': 'select',
                    'title': 'Type',
                    'name': 'type',
                    'options': [
                      'AUTH_QUERY',
                      'AUTH_RESPONSE',
                      'RESOLVER_QUERY',
                      'RESOLVER_RESPONSE',
                      'TOOL_QUERY',
                      'TOOL_RESPONSE',
                    ],
                    'required': true,
                    'description': 'AUTH_QUERY, AUTH_RESPONSE, RESOLVER_QUERY, RESOLVER_RESPONSE, ..., TOOL_QUERY, TOOL_RESPONSE',
                  },
                  'socket_family': {
                    'type': 'string',
                    'input': 'select',
                    'title': 'Socket Family',
                    'name': 'socket_family',
                    'options': ['INET', 'INET6'],
                    'required': true,
                    'description': 'INET, INET6',
                  },
                  'socket_protocol': {
                    'type': 'string',
                    'input': 'select',
                    'title': 'Socket Protocol',
                    'name': 'socket_protocol',
                    'options': ['UDP', 'TCP'],
                    'required': true,
                    'description': 'UDP, TCP',
                  },
                  'query_address':
                    {
                      'type': 'string',
                      'input': 'text',
                      'title': 'Query Address',
                      'name': 'query_address',
                      'required': false,
                      'description': '',
                    },
                  'query_port': {
                    'type': 'string',
                    'input': 'text',
                    'title': 'Query Port',
                    'name': 'query_port',
                    'required': false,
                    'description': '',
                  },
                  'response_address': {
                    'type': 'string',
                    'input': 'text',
                    'title': 'Response Address',
                    'name': 'response_address',
                    'required': false,
                    'description': '',
                  },
                },
              },
            },
          };
          break;
        case 'pktvisor/handlers':
          resp = {
            'data': {
              'dns': {
                'version': '1.0',
                'config': {
                  'filter_exclude_noerror': {
                    'title': 'Filter: Exclude NOERROR',
                    'name': 'filter_exclude_noerror',
                    'type': 'checkbox',
                    'input': 'checkbox',
                    'description': 'Filter out all NOERROR responses',
                  },
                  'filter_only_rcode': {
                    'title': 'Filter: Include Only RCode',
                    'name': 'filter_only_rcode',
                    'type': 'number',
                    'input': 'number',
                    'description': 'Filter out any queries which are not the given RCODE',
                    'min': 0,
                    'max': 65536,
                    'step': 1,
                  },
                  'filter_only_qname_suffix': {
                    'title': 'Filter: Include Only QName With Suffix',
                    'name': 'filter_only_qname_suffix',
                    'type': 'text',
                    'input': 'text',
                    'pattern': '(\w+);',
                    'description': 'Filter out any queries whose QName does not end in a suffix on the list',
                  },
                },
              },
              'pcap': {
                'version': '1.0',
                'config': {},
              },
              'net': {
                'version': '1.0',
                'config': {},
              },
            },
          };
          break;
        default:
          resp = 'error';
      }
      return new Observable(subscriber => {
        subscriber.next(resp);
      });
    }
  }
}
