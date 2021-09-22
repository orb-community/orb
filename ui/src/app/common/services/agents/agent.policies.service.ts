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

  backendsCache: OrbPagination<{[propName: string]: any}>;

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
    //       this.notificationsService.error('Failed to get Available Backends',
    //         `Error: ${ err.status } - ${ err.statusText }`);
    //       return Observable.throwError(err);
    //     },
    //   );
    // TODO remove mock and uncomment http request
    return new Observable(subscriber => {
      subscriber.next({
        'pktvisor': {
          'taps': {
            'pcap': {
              'version': '1.0',
              'info': {
                'available_iface': ['eth0', 'eth1'],
              },
              'config': {
                'iface': {
                  'required': true,
                  'type': 'string',
                  'input': 'string',
                  'title': 'Interface',
                  'iface': 'string',
                  'description': 'The ethernet interface to capture on',
                },
                'bpf': {
                  'required': false,
                  'type': 'string',
                  'input': 'string',
                  'title': 'Filter Expression',
                  'bpf': 'string',
                  'description': 'tcpdump compatible filter expression for limiting the traffic examined (with BPF). Example: "port 53"',
                },
                'host_spec': {
                  'required': false,
                  'type': 'string',
                  'input': 'string',
                  'title': 'Host Specification',
                  'host_spec': 'string',
                  'description': 'Subnets (comma separated) to consider this HOST, in CIDR form. Example: "10.0.1.0/24,10.0.2.1/32,2001:db8::/64"',
                },
                'pcap_source': {
                  'required': false,
                  'type': 'string',
                  'input': 'string',
                  'title': 'pcap Engine',
                  'pcap_source': 'string',
                  'description': 'pcap backend engine to use. Defaults to best for platform.',
                },
              },
            },
            'dnstap': {
              'version': '1.0',
              'config': {
                'type': {
                  'type': 'text',
                  'input': 'select',
                  'title': 'Type',
                  'name': 'Type',
                  'options': [
                    'AUTH_QUERY',
                    'AUTH_RESPONSE',
                    'RESOLVER_QUERY',
                    'RESOLVER_RESPONSE',
                    'TOOL_QUERY',
                    'TOOL_RESPONSE',
                  ],
                  'required': true,
                  'description': 'AUTH_QUERY, AUTH_RESPONSE, RESOLVER_QUERY,\n' +
                    'RESOLVER_RESPONSE, ..., TOOL_QUERY, TOOL_RESPONSE',
                },
                'socket_family': {
                  'type': 'text',
                  'input': 'select',
                  'title': 'Socket Family',
                  'name': 'Socket Family',
                  'options': ['INET', 'INET6'],
                  'required': true,
                  'description': 'INET, INET6',
                },
                'socket_protocol': {
                  'type': 'text',
                  'input': 'select',
                  'title': 'Socket Protocol',
                  'name': 'Socket Protocol',
                  'options': ['UDP', 'TCP'],
                  'required': true,
                  'description': 'UDP, TCP',
                },
                'query_address': {
                  'type': 'text',
                  'input': 'text',
                  'title': 'Query Address',
                  'name': 'Query Address',
                  'required': false,
                  'description': '',
                },
                'query_port': {
                  'type': 'text',
                  'input': 'text',
                  'title': 'Query Port',
                  'name': 'Query Port',
                  'required': false,
                  'description': '',
                },
                'response_address': {
                  'type': 'text',
                  'input': 'text',
                  'title': 'Response Address',
                  'name': 'Response Address',
                  'required': false,
                  'description': '',
                },
              },
            },
          },
          'handlers': {
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
                'metric_groups': {
                  'input': 'checkgroup',
                  'title': 'Metric Groups',
                  'name': 'metric_groups',
                  'cardinality': {
                    'input': 'checkbox',
                    'type': 'checkbox',
                    'name': 'cardinality',
                    'title': 'Cardinality',
                    'description': 'Metrics counting the unique number of items in the stream',
                    'metrics': [
                      'cardinality.qname',
                    ],
                  },
                  'dns_transactions': {
                    'input': 'checkbox',
                    'type': 'checkbox',
                    'name': 'dns_transactions',
                    'title': 'DNS Transactions (Query/Reply pairs)',
                    'description': 'Metrics based on tracking queries and their associated replies',
                    'metrics': [
                      'xact.counts.timed_out',
                      'xact.counts.total',
                      'xact.in.top_slow',
                    ],
                  },
                  'top_dns_wire': {
                    'input': 'checkbox',
                    'type': 'checkbox',
                    'name': 'top_dns_wire',
                    'title': 'Top N Metrics (Various)',
                    'description': 'Top N metrics across various details from the DNS wire packets',
                    'metrics': [],
                  },
                  'top_qnames': {
                    'input': 'checkbox',
                    'type': 'checkbox',
                    'name': 'top_qnames',
                    'title': 'Top N QNames (All)',
                    'description': 'Top QNames across all DNS queries in stream',
                    'metrics': [],
                  },
                  'top_qnames_by_rcode': {
                    'input': 'checkbox',
                    'type': 'checkbox',
                    'name': 'top_qnames_by_rcode',
                    'title': 'Top N QNames (Failing RCodes) ',
                    'description': 'Top QNames across failing result codes',
                    'metrics': [],
                  },
                },
              },
              'metric_groups': {
                'title': 'Metric Groups',
                'name': 'metric_groups',
                'cardinality': {
                  'input': 'checkbox',
                  'type': 'checkbox',
                  'name': 'cardinality',
                  'title': 'Cardinality',
                  'description': 'Metrics counting the unique number of items in the stream',
                  'metrics': [
                    'cardinality.qname',
                  ],
                },
                'dns_transactions': {
                  'input': 'checkbox',
                  'type': 'checkbox',
                  'name': 'dns_transactions',
                  'title': 'DNS Transactions (Query/Reply pairs)',
                  'description': 'Metrics based on tracking queries and their associated replies',
                  'metrics': [
                    'xact.counts.timed_out',
                    'xact.counts.total',
                    'xact.in.top_slow',
                  ],
                },
                'top_dns_wire': {
                  'input': 'checkbox',
                  'type': 'checkbox',
                  'name': 'top_dns_wire',
                  'title': 'Top N Metrics (Various)',
                  'description': 'Top N metrics across various details from the DNS wire packets',
                  'metrics': [],
                },
                'top_qnames': {
                  'input': 'checkbox',
                  'type': 'checkbox',
                  'name': 'top_qnames',
                  'title': 'Top N QNames (All)',
                  'description': 'Top QNames across all DNS queries in stream',
                  'metrics': [],
                },
                'top_qnames_by_rcode': {
                  'input': 'checkbox',
                  'type': 'checkbox',
                  'name': 'top_qnames_by_rcode',
                  'title': 'Top N QNames (Failing RCodes) ',
                  'description': 'Top QNames across failing result codes',
                  'metrics': [],
                },
              },
            },
            'net': {
              'version': '1.0',
              'config': {},
              'metric_groups': {
                'ip_cardinality': {
                  'title': 'IP Address Cardinality',
                  'description': 'Unique IP addresses seen in the stream',
                  'metrics': [
                    'cardinality.dst_ips_out',
                    'cardinality.src_ips_in',
                  ],
                },
                'top_geo': {
                  'title': 'Top Geo',
                  'description': 'Top Geo IP and ASN in the stream',
                  'metrics': [
                    'top_ASN',
                    'top_geoLoc',
                  ],
                },
                'top_ips': {
                  'title': 'Top IPs',
                  'description': 'Top IP addresses in the stream',
                  'metrics': [
                    'top_ipv4',
                    'top_ipv6',
                  ],
                },
              },
            },
          },
          'input': {
            'pcap': {
              'version': '1.0',
              'info': {
                'available_iface': ['eth0', 'eth1'],
              },
              'config': {
                'iface': {
                  'required': true,
                  'type': 'string',
                  'input': 'string',
                  'title': 'Interface',
                  'iface': 'string',
                  'description': 'The ethernet interface to capture on',
                },
                'bpf': {
                  'required': false,
                  'type': 'string',
                  'input': 'string',
                  'title': 'Filter Expression',
                  'bpf': 'string',
                  'description': 'tcpdump compatible filter expression for limiting the traffic examined (with BPF). Example: "port 53"',
                },
                'host_spec': {
                  'required': false,
                  'type': 'string',
                  'input': 'string',
                  'title': 'Host Specification',
                  'host_spec': 'string',
                  'description': 'Subnets (comma separated) to consider this HOST, in CIDR form. Example: "10.0.1.0/24,10.0.2.1/32,2001:db8::/64"',
                },
                'pcap_source': {
                  'required': false,
                  'type': 'string',
                  'input': 'string',
                  'title': 'pcap Engine',
                  'pcap_source': 'string',
                  'description': 'pcap backend engine to use. Defaults to best for platform.',
                },
              },
            },
            'dnstap': {
              'version': '1.0',
              'config': {
                'type': {
                  'type': 'text',
                  'input': 'select',
                  'title': 'Type',
                  'name': 'Type',
                  'options': [
                    'AUTH_QUERY',
                    'AUTH_RESPONSE',
                    'RESOLVER_QUERY',
                    'RESOLVER_RESPONSE',
                    'TOOL_QUERY',
                    'TOOL_RESPONSE',
                  ],
                  'required': true,
                  'description': 'AUTH_QUERY, AUTH_RESPONSE, RESOLVER_QUERY,\n' +
                    'RESOLVER_RESPONSE, ..., TOOL_QUERY, TOOL_RESPONSE',
                },
                'socket_family': {
                  'type': 'text',
                  'input': 'select',
                  'title': 'Socket Family',
                  'name': 'Socket Family',
                  'options': ['INET', 'INET6'],
                  'required': true,
                  'description': 'INET, INET6',
                },
                'socket_protocol': {
                  'type': 'text',
                  'input': 'select',
                  'title': 'Socket Protocol',
                  'name': 'Socket Protocol',
                  'options': ['UDP', 'TCP'],
                  'required': true,
                  'description': 'UDP, TCP',
                },
                'query_address': {
                  'type': 'text',
                  'input': 'text',
                  'title': 'Query Address',
                  'name': 'Query Address',
                  'required': false,
                  'description': '',
                },
                'query_port': {
                  'type': 'text',
                  'input': 'text',
                  'title': 'Query Port',
                  'name': 'Query Port',
                  'required': false,
                  'description': '',
                },
                'response_address': {
                  'type': 'text',
                  'input': 'text',
                  'title': 'Response Address',
                  'name': 'Response Address',
                  'required': false,
                  'description': '',
                },
              },
            },
          },
        },
      });

    });
  }

  // getPktVisorInputs() {
    // return this.http.get(environment.pktvisorTapsUrl)
    //   .map((resp: any) => {
    //     return resp.input;
    //   }).catch(
    //     err => {
    //       this.notificationsService.error('Failed to get PKTVisor Inputs',
    //         `Error: ${ err.status } - ${ err.statusText }`);
    //       return Observable.throwError(err);
    //     },
    //   );
  //
  //   return new Observable(subscriber => {
  //     subscriber.next({
  //       'pcap': {
  //         'version': '1.0',
  //         'info': {
  //           'available_iface': ['eth0', 'eth1'],
  //         },
  //         'config': {
  //           'iface': {
  //             'required': true,
  //             'type': 'string',
  //             'input': 'string',
  //             'title': 'Interface',
  //             'iface': 'string',
  //             'description': 'The ethernet interface to capture on',
  //           },
  //           'bpf': {
  //             'required': false,
  //             'type': 'string',
  //             'input': 'string',
  //             'title': 'Filter Expression',
  //             'bpf': 'string',
  //             'description': 'tcpdump compatible filter expression for limiting the traffic examined (with BPF). Example: "port 53"',
  //           },
  //           'host_spec': {
  //             'required': false,
  //             'type': 'string',
  //             'input': 'string',
  //             'title': 'Host Specification',
  //             'host_spec': 'string',
  //             'description': 'Subnets (comma separated) to consider this HOST, in CIDR form. Example: "10.0.1.0/24,10.0.2.1/32,2001:db8::/64"',
  //           },
  //           'pcap_source': {
  //             'required': false,
  //             'type': 'string',
  //             'input': 'string',
  //             'title': 'pcap Engine',
  //             'pcap_source': 'string',
  //             'description': 'pcap backend engine to use. Defaults to best for platform.',
  //           },
  //         },
  //       },
  //       'dnstap': {
  //         'version': '1.0',
  //         'config': {
  //           'type': {
  //             'type': 'text',
  //             'input': 'select',
  //             'title': 'Type',
  //             'name': 'Type',
  //             'options': [
  //               'AUTH_QUERY',
  //               'AUTH_RESPONSE',
  //               'RESOLVER_QUERY',
  //               'RESOLVER_RESPONSE',
  //               'TOOL_QUERY',
  //               'TOOL_RESPONSE',
  //             ],
  //             'required': true,
  //             'description': 'AUTH_QUERY, AUTH_RESPONSE, RESOLVER_QUERY,\n' +
  //               'RESOLVER_RESPONSE, ..., TOOL_QUERY, TOOL_RESPONSE',
  //           },
  //           'socket_family': {
  //             'type': 'text',
  //             'input': 'select',
  //             'title': 'Socket Family',
  //             'name': 'Socket Family',
  //             'options': ['INET', 'INET6'],
  //             'required': true,
  //             'description': 'INET, INET6',
  //           },
  //           'socket_protocol': {
  //             'type': 'text',
  //             'input': 'select',
  //             'title': 'Socket Protocol',
  //             'name': 'Socket Protocol',
  //             'options': ['UDP', 'TCP'],
  //             'required': true,
  //             'description': 'UDP, TCP',
  //           },
  //           'query_address': {
  //             'type': 'text',
  //             'input': 'text',
  //             'title': 'Query Address',
  //             'name': 'Query Address',
  //             'required': false,
  //             'description': '',
  //           },
  //           'query_port': {
  //             'type': 'text',
  //             'input': 'text',
  //             'title': 'Query Port',
  //             'name': 'Query Port',
  //             'required': false,
  //             'description': '',
  //           },
  //           'response_address': {
  //             'type': 'text',
  //             'input': 'text',
  //             'title': 'Response Address',
  //             'name': 'Response Address',
  //             'required': false,
  //             'description': '',
  //           },
  //         },
  //       },
  //     });
  //   });
  // }

  // getPktVisorHandlers() {
    // return this.http.get(environment.pktvisorTapsUrl)
    //   .map((resp: any) => {
    //     return resp.handler;
    //   }).catch(
    //     err => {
    //       this.notificationsService.error('Failed to get Taps Configurations',
    //         `Error: ${ err.status } - ${ err.statusText }`);
    //       return Observable.throwError(err);
    //     },
    //   );
  //
  //   return new Observable(subscriber => {
  //     subscriber.next({
  //       'dns': {
  //         'version': '1.0',
  //         'config': {
  //           'filter_exclude_noerror': {
  //             'title': 'Filter: Exclude NOERROR',
  //             'name': 'filter_exclude_noerror',
  //             'type': 'checkbox',
  //             'input': 'checkbox',
  //             'description': 'Filter out all NOERROR responses',
  //           },
  //           'filter_only_rcode': {
  //             'title': 'Filter: Include Only RCode',
  //             'name': 'filter_only_rcode',
  //             'type': 'number',
  //             'input': 'number',
  //             'description': 'Filter out any queries which are not the given RCODE',
  //             'min': 0,
  //             'max': 65536,
  //             'step': 1,
  //           },
  //           'filter_only_qname_suffix': {
  //             'title': 'Filter: Include Only QName With Suffix',
  //             'name': 'filter_only_qname_suffix',
  //             'type': 'text',
  //             'input': 'text',
  //             'pattern': '(\w+);',
  //             'description': 'Filter out any queries whose QName does not end in a suffix on the list',
  //           },
  //           'metric_groups': {
  //             'input': 'checkgroup',
  //             'title': 'Metric Groups',
  //             'name': 'metric_groups',
  //             'cardinality': {
  //               'input': 'checkbox',
  //               'type': 'checkbox',
  //               'name': 'cardinality',
  //               'title': 'Cardinality',
  //               'description': 'Metrics counting the unique number of items in the stream',
  //               'metrics': [
  //                 'cardinality.qname',
  //               ],
  //             },
  //             'dns_transactions': {
  //               'input': 'checkbox',
  //               'type': 'checkbox',
  //               'name': 'dns_transactions',
  //               'title': 'DNS Transactions (Query/Reply pairs)',
  //               'description': 'Metrics based on tracking queries and their associated replies',
  //               'metrics': [
  //                 'xact.counts.timed_out',
  //                 'xact.counts.total',
  //                 'xact.in.top_slow',
  //               ],
  //             },
  //             'top_dns_wire': {
  //               'input': 'checkbox',
  //               'type': 'checkbox',
  //               'name': 'top_dns_wire',
  //               'title': 'Top N Metrics (Various)',
  //               'description': 'Top N metrics across various details from the DNS wire packets',
  //               'metrics': [],
  //             },
  //             'top_qnames': {
  //               'input': 'checkbox',
  //               'type': 'checkbox',
  //               'name': 'top_qnames',
  //               'title': 'Top N QNames (All)',
  //               'description': 'Top QNames across all DNS queries in stream',
  //               'metrics': [],
  //             },
  //             'top_qnames_by_rcode': {
  //               'input': 'checkbox',
  //               'type': 'checkbox',
  //               'name': 'top_qnames_by_rcode',
  //               'title': 'Top N QNames (Failing RCodes) ',
  //               'description': 'Top QNames across failing result codes',
  //               'metrics': [],
  //             },
  //           },
  //         },
  //         'metric_groups': {
  //           'title': 'Metric Groups',
  //           'name': 'metric_groups',
  //           'cardinality': {
  //             'input': 'checkbox',
  //             'type': 'checkbox',
  //             'name': 'cardinality',
  //             'title': 'Cardinality',
  //             'description': 'Metrics counting the unique number of items in the stream',
  //             'metrics': [
  //               'cardinality.qname',
  //             ],
  //           },
  //           'dns_transactions': {
  //             'input': 'checkbox',
  //             'type': 'checkbox',
  //             'name': 'dns_transactions',
  //             'title': 'DNS Transactions (Query/Reply pairs)',
  //             'description': 'Metrics based on tracking queries and their associated replies',
  //             'metrics': [
  //               'xact.counts.timed_out',
  //               'xact.counts.total',
  //               'xact.in.top_slow',
  //             ],
  //           },
  //           'top_dns_wire': {
  //             'input': 'checkbox',
  //             'type': 'checkbox',
  //             'name': 'top_dns_wire',
  //             'title': 'Top N Metrics (Various)',
  //             'description': 'Top N metrics across various details from the DNS wire packets',
  //             'metrics': [],
  //           },
  //           'top_qnames': {
  //             'input': 'checkbox',
  //             'type': 'checkbox',
  //             'name': 'top_qnames',
  //             'title': 'Top N QNames (All)',
  //             'description': 'Top QNames across all DNS queries in stream',
  //             'metrics': [],
  //           },
  //           'top_qnames_by_rcode': {
  //             'input': 'checkbox',
  //             'type': 'checkbox',
  //             'name': 'top_qnames_by_rcode',
  //             'title': 'Top N QNames (Failing RCodes) ',
  //             'description': 'Top QNames across failing result codes',
  //             'metrics': [],
  //           },
  //         },
  //       },
  //       'net': {
  //         'version': '1.0',
  //         'config': {},
  //         'metric_groups': {
  //           'ip_cardinality': {
  //             'title': 'IP Address Cardinality',
  //             'description': 'Unique IP addresses seen in the stream',
  //             'metrics': [
  //               'cardinality.dst_ips_out',
  //               'cardinality.src_ips_in',
  //             ],
  //           },
  //           'top_geo': {
  //             'title': 'Top Geo',
  //             'description': 'Top Geo IP and ASN in the stream',
  //             'metrics': [
  //               'top_ASN',
  //               'top_geoLoc',
  //             ],
  //           },
  //           'top_ips': {
  //             'title': 'Top IPs',
  //             'description': 'Top IP addresses in the stream',
  //             'metrics': [
  //               'top_ipv4',
  //               'top_ipv6',
  //             ],
  //           },
  //         },
  //       },
  //     });
  //   });
  // }
}
