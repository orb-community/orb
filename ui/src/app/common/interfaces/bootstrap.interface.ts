export interface Config {
  thing_id: string;
  thing_key: string;
  channels: Array<string>;
  external_id: string;
  external_key: string;
  content: string;
  state: number;
}

export interface ConfigContent {
  log_level: string;
  http_port: string;
  mqtt_url: string;
  edgex_url: string;
  nats_url: string;
  export_config: ExportConfig;
}

export interface ExportConfig {
  file?: string;
  exp: ExpConf;
  mqtt: MqttConfig;
  routes: Array<Route>;
}

export interface ExpConf {
  log_level?: string;
  nats?: string;
  port?: string;
  cache_db?: string;
  cache_pass?: string;
  cache_url?: string;
}

export interface Route {
  mqtt_topic?: string;
  nats_topic?: string;
  subtopic?: string;
  type?: string;
}

export interface MqttConfig {
  host?: string;
  ca_path?: string;
  cert_path?: string;
  priv_key_path?: string;
  channel?: string;
  qos?: number;
  mtls?: boolean;
  password?: string;
  username?: string;
  skip_tls_ver?: boolean;
  retain?: boolean;
}

export interface ConfigUpdate {
  content: string;
  name: string;
}
