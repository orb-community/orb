import {HttpClient, HttpParams} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {EMPTY, Observable, of} from 'rxjs';
import 'rxjs/add/observable/empty';

import {Agent, AgentPolicyAggStates,} from 'app/common/interfaces/orb/agent.interface';
import {AgentPolicyState, AgentPolicyStates,} from 'app/common/interfaces/orb/agent.policy.interface';
import {OrbPagination} from 'app/common/interfaces/orb/pagination.interface';
import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {environment} from 'environments/environment';
import {expand, map, scan, takeWhile} from 'rxjs/operators';

export enum AvailableOS {
  DOCKER = 'docker',
}

@Injectable()
export class AgentsService {
  constructor(
      private http: HttpClient,
      private notificationsService: NotificationsService,
  ) {
  }

  getPktVisorMetrics(agentID: string) {
    const metrics = {
      'default': {
        'default-default-dhcp': {
          'dhcp': {
            'period': {
              'length': 296,
              'start_ts': 1660152794
            },
            'rates': {
              'total': {
                'live': 0,
                'p50': 0,
                'p90': 0,
                'p95': 0,
                'p99': 0
              }
            },
            'wire_packets': {
              'ack': 0,
              'deep_samples': 0,
              'discover': 0,
              'filtered': 0,
              'offer': 0,
              'request': 0,
              'total': 0
            }
          }
        },
        'default-default-dns': {
          'dns': {
            'cardinality': {
              'qname': 1
            },
            'period': {
              'length': 296,
              'start_ts': 1660152794
            },
            'rates': {
              'total': {
                'live': 0,
                'p50': 0,
                'p90': 0,
                'p95': 0,
                'p99': 1
              }
            },
            'top_nodata': [
              {
                'estimate': 1,
                'name': 'connectivity-check.ubuntu.com'
              }
            ],
            'top_nxdomain': [],
            'top_qname2': [
              {
                'estimate': 2,
                'name': '.ubuntu.com'
              }
            ],
            'top_qname3': [
              {
                'estimate': 2,
                'name': 'connectivity-check.ubuntu.com'
              }
            ],
            'top_qname_by_resp_bytes': [
              {
                'estimate': 58,
                'name': 'connectivity-check.ubuntu.com'
              }
            ],
            'top_qtype': [
              {
                'estimate': 2,
                'name': 'AAAA'
              }
            ],
            'top_rcode': [
              {
                'estimate': 11,
                'name': 'NOERROR'
              }
            ],
            'top_refused': [],
            'top_srvfail': [],
            'top_udp_ports': [
              {
                'estimate': 10,
                'name': '5353'
              },
              {
                'estimate': 2,
                'name': '55622'
              }
            ],
            'wire_packets': {
              'deep_samples': 12,
              'filtered': 0,
              'ipv4': 10,
              'ipv6': 2,
              'nodata': 1,
              'noerror': 11,
              'nxdomain': 0,
              'queries': 1,
              'refused': 0,
              'replies': 11,
              'srvfail': 0,
              'tcp': 0,
              'total': 12,
              'udp': 12
            },
            'xact': {
              'counts': {
                'timed_out': 0,
                'total': 1
              },
              'in': {
                'top_slow': [],
                'total': 0
              },
              'out': {
                'quantiles_us': {
                  'p50': 1954,
                  'p90': 1954,
                  'p95': 1954,
                  'p99': 1954
                },
                'top_slow': [],
                'total': 1
              },
              'ratio': {
                'quantiles': {
                  'p50': 1.0,
                  'p90': 1.0,
                  'p95': 1.0,
                  'p99': 1.0
                }
              }
            }
          }
        },
        'default-default-net': {
          'packets': {
            'cardinality': {
              'dst_ips_out': 45,
              'src_ips_in': 36
            },
            'deep_samples': 140237,
            'filtered': 0,
            'in': 78827,
            'ipv4': 138244,
            'ipv6': 1970,
            'other_l4': 90,
            'out': 61387,
            'payload_size': {
              'p50': 952,
              'p90': 1000,
              'p95': 1000,
              'p99': 1000
            },
            'period': {
              'length': 276,
              'start_ts': 1660152814
            },
            'protocol': {
              'tcp': {
                'syn': 22
              }
            },
            'rates': {
              'bytes_in': {
                'live': 202634,
                'p50': 196752,
                'p90': 208529,
                'p95': 211604,
                'p99': 247499
              },
              'bytes_out': {
                'live': 123204,
                'p50': 122321,
                'p90': 141873,
                'p95': 144025,
                'p99': 153737
              },
              'pps_in': {
                'live': 277,
                'p50': 285,
                'p90': 302,
                'p95': 310,
                'p99': 353
              },
              'pps_out': {
                'live': 198,
                'p50': 207,
                'p90': 276,
                'p95': 281,
                'p99': 298
              },
              'pps_total': {
                'live': 475,
                'p50': 497,
                'p90': 565,
                'p95': 575,
                'p99': 601
              }
            },
            'tcp': 4625,
            'top_ASN': [],
            'top_geoLoc': [],
            'top_ipv4': [
              {
                'estimate': 116940,
                'name': '206.247.14.178'
              },
              {
                'estimate': 19623,
                'name': '35.215.237.233'
              },
              {
                'estimate': 620,
                'name': '66.22.202.35'
              },
              {
                'estimate': 292,
                'name': '162.159.128.235'
              },
              {
                'estimate': 125,
                'name': '164.163.6.3'
              },
              {
                'estimate': 104,
                'name': '239.255.255.250'
              },
              {
                'estimate': 96,
                'name': '54.232.34.114'
              },
              {
                'estimate': 76,
                'name': '187.16.226.151'
              },
              {
                'estimate': 64,
                'name': '162.159.136.234'
              },
              {
                'estimate': 48,
                'name': '52.109.108.52'
              }
            ],
            'top_ipv6': [
              {
                'estimate': 749,
                'name': '2620:1ec:a92::171'
              },
              {
                'estimate': 394,
                'name': '2603:1056:c03:2424::2'
              },
              {
                'estimate': 357,
                'name': '2600:1419:1e00:593::4b36'
              },
              {
                'estimate': 216,
                'name': '2602:fd3f:3:ff02::2d'
              },
              {
                'estimate': 48,
                'name': '2603:1056:c03:2401::2'
              },
              {
                'estimate': 39,
                'name': '2603:1030:b00::4ee'
              },
              {
                'estimate': 24,
                'name': '2603:1056:1400:1::'
              },
              {
                'estimate': 21,
                'name': '2a01:111:f100:3001::8987:1046'
              },
              {
                'estimate': 20,
                'name': '2800:3f0:4001:82e::2004'
              },
              {
                'estimate': 20,
                'name': 'fe80::dac6:78ff:fe39:e360'
              }
            ],
            'total': 140237,
            'udp': 135522
          }
        },
        'default-default-pcap_stats': {
          'pcap': {
            'if_drops': 0,
            'os_drops': 0,
            'period': {
              'length': 276,
              'start_ts': 1660152814
            },
            'tcp_reassembly_errors': 0
          }
        }
      },
      'default-6410a3cbc0d30617-resources': {
        'default-6410a3cbc0d30617-resources': {
          'input_resources': {
            'cpu_usage': {
              'p50': 2.0491803278688523,
              'p90': 2.658486707566462,
              'p95': 2.8112449799196786,
              'p99': 3.006012024048096
            },
            'deep_samples': 61,
            'event_rate': {
              'live': 1,
              'p50': 0,
              'p90': 1,
              'p95': 1,
              'p99': 2
            },
            'handler_count': 5,
            'memory_bytes': {
              'p50': 22462464,
              'p90': 22462464,
              'p95': 22462464,
              'p99': 22462464
            },
            'period': {
              'length': 276,
              'start_ts': 1660152814
            },
            'policy_count': 2,
            'total': 61
          }
        }
      }
    };

    return of(metrics);
  }

  addAgent(agentItem: Agent) {
    return this.http
        .post<Agent>(
            environment.agentsUrl,
            {...agentItem, validate_only: false},
            {observe: 'response'},
        )
        .map((resp) => {
          let {body: agent} = resp;
          agent = {
            ...agent,
            combined_tags: {...agent?.orb_tags, ...agent?.agent_tags},
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
            {observe: 'response'},
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
            {...agentItem, validate_only: true},
            {observe: 'response'},
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
            combined_tags: {...agent?.orb_tags, ...agent?.agent_tags},
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
            combined_tags: {...resp?.orb_tags, ...resp?.agent_tags},
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
        .get(`${environment.agentsUrl}`, {params})
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
      agent.combined_tags = {...agent?.orb_tags, ...agent?.agent_tags};
      // map agg policy state
      const {agg_info, agg_state} = this.policyAggState(agent);
      agent.policy_agg_info = agg_info;
      agent.policy_agg_state = agg_state;
      return agent;
    });
  }

  policyAggState(agent) {
    const {policy_state} = agent;
    let agg_info = 'No Policies Applied';
    let agg_state = AgentPolicyAggStates.none;

    const policies =
        (!!policy_state && (Object.values(policy_state) as AgentPolicyState[])) ||
        [];
    if (policies.length > 0) {
      let err = 0;
      policies.forEach((policy) => {
        if (policy.state !== AgentPolicyStates.running) {
          err = err + 1;
        }
      });
      if (err > 0) {
        if (err === policies.length) {
          agg_info = 'All Policies not running';
        } else {
          agg_info = `${err} out of ${policies.length} policies are not running`;
        }
        agg_state = AgentPolicyAggStates.failure;
      } else {
        agg_info = `All policies are running`;
        agg_state = AgentPolicyAggStates.healthy;
      }
    }

    return {agg_info, agg_state};
  }
}
