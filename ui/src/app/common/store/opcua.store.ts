import { Injectable } from '@angular/core';

@Injectable()
export class OpcuaStore {
  browsedNodes = {
    uri: '',
    nodes: [],
  };

  constructor(
  ) { }

  setBrowsedNodes(browse: any) {
    this.browsedNodes.uri = browse.uri;
    this.browsedNodes.nodes = browse.nodes;
  }

  getBrowsedNodes() {
    return this.browsedNodes;
  }
}
