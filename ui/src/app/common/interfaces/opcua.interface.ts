export interface OpcuaMetadata {
    type?: string;
    opcua: {
      server_uri?: string,
      node_id?: string,
    };
    channel_id?: string;
}

export interface OpcuaNode {
  name?: string;
  id?: string;
  key?: string;
  metadata?: OpcuaMetadata;
}

export interface OpcuaTableRow {
  name: string;
  id?: string;
  key?: string;
  serverURI: string;
  nodeID: string;
  messages?: number;
  metadata?: OpcuaMetadata;
  seen?: any;
}
