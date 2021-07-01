import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { OpcuaNode, OpcuaTableRow } from 'app/common/interfaces/opcua.interface';
import { ThingsService } from 'app/common/services/things/things.service';
import { ChannelsService } from 'app/common/services/channels/channels.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

const defLimit: number = 20;

@Injectable()
export class OpcuaService {
  typeOpcua = 'opcua';
  typeOpcuaServer = 'OPC-UA-Server';
  typeOpcuaNode = 'OPC-UA-Node';

  constructor(
    private http: HttpClient,
    private thingsService: ThingsService,
    private channelsService: ChannelsService,
    private notificationsService: NotificationsService,
  ) { }

  getNode(id: string) {
    return this.thingsService.getThing(id);
  }

  getNodes(offset?: number, limit?: number, name?: string) {
    const filters = {
      offset: offset || 0,
      limit: limit || defLimit,
      type: this.typeOpcua,
      name: name,
    };

    return this.thingsService.getThings(filters);
  }

  addNodes(serverURI: string, nodes: any) {
    const filters = {
      offset: 0,
      limit: 1,
      type: this.typeOpcua,
      metadata: `{"server_uri": "${serverURI}"}`,
    };

    // Check if a channel exist for serverURI
    return this.channelsService.getChannels(filters).map(
      (resp: any) => {
        if (resp.total === 0) {
          const chanReq = {
            name: `${this.typeOpcuaServer}-${serverURI}`,
            metadata: {
              type: this.typeOpcua,
              opcua: {
                server_uri: serverURI,
              },
            },
          };

          this.channelsService.addChannel(chanReq).subscribe(
            respChan => {
              const chanID = respChan.headers.get('location').replace('/channels/', '');
              this.addAndConnect(chanID, nodes);
            },
          );
        } else {
          const chanID = resp.channels[0].id;
          this.addAndConnect(chanID, nodes);
        }
      },
    );
  }

  addAndConnect(chanID: string, nodes: any) {
    const nodesReq: OpcuaNode[] = [];
    nodes.forEach(node => {
      const nodeReq: OpcuaNode = {
        name: node.name,
        metadata: {
          type: this.typeOpcua,
          opcua: {
            node_id: node.nodeID,
            server_uri: node.serverURI,
          },
          channel_id: chanID,
        },
      };
      nodesReq.push(nodeReq);
    });

    this.thingsService.addThings(nodesReq).subscribe(
      (respThings: any) => {
        const channels = [chanID];
        const nodesIDs = respThings.body.things.map( thing => thing.id);
        this.channelsService.connectThings(channels, nodesIDs).subscribe(
          respCon => {
            this.notificationsService.success('OPC-UA Nodes successfully created', '');
          },
          err => {
            nodesIDs.forEach( id => {
              this.thingsService.deleteThing(id).subscribe();
            });
          },
        );
      },
    );
  }

  editNode(node: OpcuaTableRow) {
    const nodeReq: OpcuaNode = {
      id: node.id,
      name: node.name,
      metadata: {
        type: this.typeOpcua,
        opcua: {
          server_uri: node.serverURI,
          node_id: node.nodeID,
        },
        channel_id: node.metadata.channel_id,
      },
    };

    return this.thingsService.editThing(nodeReq).map(
      resp => {
        this.notificationsService.success('OPC-UA Node successfully edited', '');
      },
    );
  }

  deleteNode(node: any) {
    return this.thingsService.deleteThing(node.id).map(
      respThing => {
        const serverURI = node.metadata.opcua.server_uri;
        const filters = {
          offset: 0,
          limit: 1,
          type: this.typeOpcua,
          metadata: `{"server_uri": "${serverURI}"}`,
        };
        this.thingsService.getThings(filters).subscribe(
          (respChan: any) => {
            if (respChan.total === 0) {
              const channelID = node.metadata.channel_id;
              this.channelsService.deleteChannel(channelID).subscribe();
            }
          },
        );
        this.notificationsService.success('OPC-UA Node successfully deleted', '');
      },
    );
  }

  browseServerNodes(uri: string, ns: string, id: string) {
    const params = new HttpParams()
      .set('server', uri)
      .set('namespace', ns)
      .set('identifier', id);

    return this.http.get(environment.browseUrl, { params })
      .map(
        resp => {
          this.notificationsService.success('OPC-UA browsing finished', '');
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Browse',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }
}
