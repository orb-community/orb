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
    const offset = !!pageInfo ? pageInfo.offset : this.cache.offset;
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
    //   }).catch(err => {
    //       this.notificationsService.error('Failed to get Available Backends',
    //         `Error: ${ err.status } - ${ err.statusText }`);
    //       return Observable.throwError(err);
    //     });
    // TODO uncomment mock above
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
    // TODO remove this if and uncomment code above - this allows only for taps
    if (final === 'pktvisor/taps') {
      return this.http.get(`${environment.agentsBackendUrl}/${final}`)
      .map((response: any) => {
        return response.backend;
      }).catch(
        err => {
          this.notificationsService.error('Failed to get Available Backends',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
    }

    let resp;
    switch (final) {
      // case 'pktvisor/taps':
      //   resp = {
      //     data: [
      //       {
      //         'name': 'ethernet',
      //         'input_type': 'pcap',
      //         'config_predefined': [
      //           'iface',
      //         ],
      //         'agents': {
      //           'total': 1,
      //         },
      //       },
      //     ],
      //   };
      //   break;
      case 'pktvisor/inputs':
        resp = {
          data: {
            'pcap': {
              '1.0': {
                'filter': {
                  'bpf': {
                    'type': 'string',
                    'input': 'text',
                    'label': 'Filter Expression',
                    'description': 'tcpdump compatible filter expression for limiting the traffic examined (with BPF). See https://www.tcpdump.org/manpages/tcpdump.1.html',
                    'props': {
                      'example': 'udp port 53 and host 127.0.0.1',
                    },
                  },
                },
                'config': {
                  'iface': {
                    'type': 'string',
                    'input': 'text',
                    'label': 'Network Interface',
                    'description': 'The network interface to capture traffic from',
                    'props': {
                      'required': true,
                      'example': 'eth0',
                    },
                  },
                  'host_spec': {
                    'type': 'string',
                    'input': 'text',
                    'label': 'Host Specification',
                    'description': 'Subnets (comma separated) which should be considered belonging to this host, in CIDR form. Used for ingress/egress determination, defaults to host attached to the network interface.',
                    'props': {
                      'advanced': true,
                      'example': '10.0.1.0/24,10.0.2.1/32,2001:db8::/64',
                    },
                  },
                  'pcap_source': {
                    'type': 'string',
                    'input': 'select',
                    'label': 'Packet Capture Engine',
                    'description': 'Packet capture engine to use. Defaults to best for platform.',
                    'props': {
                      'advanced': true,
                      'example': 'libpcap',
                      'options': {
                        'libpcap': 'libpcap',
                        'af_packet (linux only)': 'af_packet',
                      },
                    },
                  },
                },
              },
            },
          },
        };
        break;
      case 'pktvisor/handlers':
        resp = {
          data: {
            'dns': {
              '1.0': {
                'filter': {
                  'exclude_noerror': {
                    'label': 'Exclude NOERROR',
                    'type': 'bool',
                    'input': 'checkbox',
                    'description': 'Filter out all NOERROR responses',
                  },
                  'only_rcode': {
                    'label': 'Include Only RCODE',
                    'type': 'number',
                    'input': 'select',
                    'description': 'Filter out any queries which are not the given RCODE',
                    'props': {
                      'allow_custom_options': true,
                      'options': {
                        'NOERROR': 0,
                        'SERVFAIL': 2,
                        'NXDOMAIN': 3,
                        'REFUSED': 5,
                      },
                    },
                  },
                  'only_qname_suffix': {
                    'label': 'Include Only QName With Suffix',
                    'type': 'string[]',
                    'input': 'text',
                    'description': 'Filter out any queries whose QName does not end in a suffix on the list',
                    'props': {
                      'example': '.foo.com,.example.com',
                    },
                  },
                },
                'config': {},
                'metrics': {},
                'metric_groups': {
                  'cardinality': {
                    'label': 'Cardinality',
                    'description': 'Metrics counting the unique number of items in the stream',
                    'metrics': [],
                  },
                  'dns_transactions': {
                    'label': 'DNS Transactions (Query/Reply pairs)',
                    'description': 'Metrics based on tracking queries and their associated replies',
                    'metrics': [],
                  },
                  'top_dns_wire': {
                    'label': 'Top N Metrics (Various)',
                    'description': 'Top N metrics across various details from the DNS wire packets',
                    'metrics': [],
                  },
                  'top_qnames': {
                    'label': 'Top N QNames (All)',
                    'description': 'Top QNames across all DNS queries in stream',
                    'metrics': [],
                  },
                  'top_qnames_by_rcode': {
                    'label': 'Top N QNames (Failing RCodes) ',
                    'description': 'Top QNames across failing result codes',
                    'metrics': [],
                  },
                },
              },
            },
            'net': {
              '1.0': {
                'filter': {},
                'config': {},
                'metrics': {},
                'metric_groups': {
                  'ip_cardinality': {
                    'label': 'IP Address Cardinality',
                    'description': 'Unique IP addresses seen in the stream',
                    'metrics': [],
                  },
                  'top_geo': {
                    'label': 'Top Geo',
                    'description': 'Top Geo IP and ASN in the stream',
                    'metrics': [],
                  },
                  'top_ips': {
                    'label': 'Top IPs',
                    'description': 'Top IP addresses in the stream',
                    'metrics': [],
                  },
                },
              },
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
