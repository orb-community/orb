export interface GatewayMetadata {
  ctrl_channel_id?: string;
  data_channel_id?: string;
  export_channel_id?: string;
  export_thing_id?: string;
  external_key?: string;
  external_id?: string;
  cfg_id?: string;
  type?: string;
}

export interface Gateway {
  id?: string;
  key?: string;
  name?: string;
  externalID?: string;
  metadata?: GatewayMetadata;
  seen?: string;
  messages?: string;
}
