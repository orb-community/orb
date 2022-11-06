/**
 * Prometheus Sink Config Interface
 *
 * [Prometheus Sink]{@link https://github.com/etaques/orb/blob/develop/cmd/prom-sink/main.go}
 * [Sinks Architecture]{@link https://github.com/etaques/orb/wiki/Architecture:-Sinks}
 */
import { SinkConfig } from 'app/common/interfaces/orb/sink/sink.config.interface';

/**
 * @interface PrometheusConfig
 */
export interface PrometheusConfig extends SinkConfig<string> {
  name: 'Prometheus';
  /**
   *  Remote Host URL {string}
   */
  remote_host?: string;

  /**
   *  Username|Email(?) {string}
   */
  username?: string;

  /**
   *  Password {string}
   */
  password?: string;
}
