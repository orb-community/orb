import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { OpcuaService } from 'app/common/services/opcua/opcua.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { MessagesService } from 'app/common/services/messages/messages.service';
import { OpcuaStore } from 'app/common/store/opcua.store';
import { FsService } from 'app/common/services/fs/fs.service';
import { OpcuaTableRow } from 'app/common/interfaces/opcua.interface';
import { OpcuaAddComponent } from 'app/pages/services/opcua/add/opcua.add.component';

const defSearchBarMs: number = 100;

@Component({
  selector: 'ngx-opcua-component',
  templateUrl: './opcua.component.html',
  styleUrls: ['./opcua.component.scss'],
})
export class OpcuaComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', '', 'Name', 'Server URI', 'Node ID', 'Messages', 'Last Seen'],
    keys: ['edit', 'delete', 'details', 'name', 'serverURI', 'nodeID', 'messages', 'seen'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  browseServerURI = '';
  // Standard root OPC-UA server NodeID (ns=0;i=84)
  browseNamespace = '';
  browseIdentifier = '';

  browsedNodes = [];
  browseSearch = [];
  checkedNodes = [];

  offset = 0;
  limit = 20;
  total = 0;

  searchTime = 0;
  columnChar = '|';

  constructor(
    private router: Router,
    private opcuaService: OpcuaService,
    private messagesService: MessagesService,
    private notificationsService: NotificationsService,
    private opcuaStore: OpcuaStore,
    private dialogService: NbDialogService,
    private fsService: FsService,
  ) { }

  ngOnInit() {
    const browseStore = this.opcuaStore.getBrowsedNodes();
    this.browseServerURI = browseStore.uri;
    this.browsedNodes = browseStore.nodes;
    this.getOpcuaNodes();
  }

  getOpcuaNodes(name?: string): void {
    this.opcuaService.getNodes(this.offset, this.limit, name).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.things,
        };

        this.page.rows.forEach((node: OpcuaTableRow) => {
          node.serverURI = node.metadata.opcua.server_uri;
          node.nodeID = node.metadata.opcua.node_id;

          const chanID: string = node.metadata ? node.metadata.channel_id : '';
          this.messagesService.getMessages(chanID, node.key, {publisher: node.id}).subscribe(
            (msgResp: any) => {
              if (msgResp.messages) {
                node.seen = msgResp.messages[0].time;
                node.messages = msgResp.total;
              }
            },
          );
        });
      },
    );
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      this.pageFilters.offset = this.page.offset - this.page.limit;
    }
    if (dir === 'next') {
      this.pageFilters.offset = this.page.offset + this.page.limit;
    }
    this.getOpcuaNodes();
  }

  onChangeLimit(lim: number) {
    this.pageFilters.limit = lim;
    this.getOpcuaNodes();
  }

  openAddModal() {
    this.dialogService.open(OpcuaAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          setTimeout(
            () => {
              this.getOpcuaNodes();
            }, 3000,
          );
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(OpcuaAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          setTimeout(
            () => {
              this.getOpcuaNodes();
            }, 3000,
          );
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'OPC-UA Node' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.opcuaService.deleteNode(row).subscribe(
            resp => {
              this.notificationsService.success('OPC-UA Node successfully deleted', '');
              setTimeout(
                () => {
                  this.getOpcuaNodes();
                }, 3000,
              );
            },
          );
        }
      },
    );
  }

  onOpenDetails(row: any) {
    if (row.id) {
      this.router.navigate([`${this.router.routerState.snapshot.url}/details/${row.id}`]);
    }
  }

  browseOpcuaNodes() {
    this.opcuaService.browseServerNodes(this.browseServerURI, this.browseNamespace, this.browseIdentifier).subscribe(
      (resp: any) => {
        this.browsedNodes = resp.nodes;
        const browsedNodes = {
          uri: this.browseServerURI,
          nodes: this.browsedNodes,
        };
        this.opcuaStore.setBrowsedNodes(browsedNodes);
      },
      err => {
      },
    );
  }

  onCheckboxChanged(event: boolean, node: any) {
    if (event === true) {
      this.checkedNodes.push(node.NodeID);
    } else {
      this.checkedNodes = this.checkedNodes.filter(n => n !== node.NodeID);
    }
  }

  subscribeOpcuaNodes() {
    const nodesReq = [];
    this.checkedNodes.forEach( (checkedNode, i) => {
      const nodeReq = {
        name: checkedNode,
        serverURI: this.browseServerURI,
        nodeID: checkedNode,
      };

      // Check if subscription already exist
      if (!this.isSubscribed(this.browseServerURI, checkedNode)) {
        nodesReq.push(nodeReq);
      }
    });

    this.opcuaService.addNodes(this.browseServerURI,  nodesReq).subscribe(
      resp => {
        setTimeout(
          () => {
            this.getOpcuaNodes();
          }, 3000,
        );
      },
    );
  }

  isSubscribed(serverURI: string, nodeID: string) {
    const subs = this.page.rows.filter((n: OpcuaTableRow) => n.serverURI === serverURI && n.nodeID === nodeID);
    if (subs.length !== 0) {
      this.notificationsService.warn(`Subscribtion to server ${serverURI} and nodeID ${nodeID} already exist`, '');
      return true;
    }

    return false;
  }

  searchNode(input) {
    const t = new Date().getTime();
    if ((t - this.searchTime) > defSearchBarMs) {
      this.getOpcuaNodes(input);
      this.searchTime = t;
    }
  }

  searchBrowse(input) {
    const browseStore = this.opcuaStore.getBrowsedNodes();
    this.browsedNodes = browseStore.nodes.filter( node =>
      (node.NodeID.includes(input) ||
      node.Description.includes(input) ||
      node.DataType.includes(input) ||
      node.BrowseName.includes(input)));
  }

  onClickSave() {
    this.fsService.exportToCsv('opcua_nodes.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
    if (files && files.length > 0) {
      const file: File = files.item(0);
      const reader: FileReader = new FileReader();
      reader.readAsText(file);
      reader.onload = () => {
        const csv: string = reader.result as string;
        const lines = csv.split('\n');

        // Split all file lines using a separator
        lines.forEach( (line, i) => {
          const cols = line.split(this.columnChar);
          const name = cols[0];
          const nodes = [];
          if (name !== '' && name !== '<empty string>' && cols.length > 2) {
            const serv = cols[1];
            for (let j = 2; j < cols.length; j++) {
              const node = {
                name: cols[0],
                serverURI: serv,
                nodeID: cols[j],
              };
              nodes.push(node);
            }

            this.opcuaService.addNodes(serv, nodes).subscribe(
              resp => {
                setTimeout( () => {
                  this.getOpcuaNodes();
                }, 3000 * i);
              },
            );
          } else {
            this.notificationsService.warn('Incomplete line found in file', '');
          }
        });
      };
    }
  }
}
