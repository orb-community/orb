/**
 * Prometheus Sink Config Interface
 * for more details:
 * /src/cmd/prom-sink/main.go
 * https://github.com/ns1labs/orb/wiki/Architecture:-Sinks
 */
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';

export interface PrometheusConfig extends SinkConfig<string> {
  /** Remote Host Name: string */
  remote_host?: string;
  /** Username: string */
  username?: string;
  /** Password: string */
  password?: string;
}
