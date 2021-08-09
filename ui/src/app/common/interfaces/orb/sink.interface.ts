/**
 * Sink Data Model Interface
 *
 * https://github.com/ns1labs/orb/wiki/Architecture:-Sinks
 */

import { PrometheusConfig } from 'app/common/interfaces/orb/sink.config/prometheus.config.interface';
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';

export interface Sink {
  /** id: UUIDv4 (read only) */
  id?: string;
  /** Name: string [a-zA-Z_:][a-zA-Z0-9_]* */
  name?: string;
  /** Description: string */
  description?: string;
  /**
   * ORB Tags: orb_tags string<JSON>
   * simple key/values - no recursive objects
   */
  tags?: { [propName: string]: string };
  /** Status: string ['active'|'error'] */
  status?: string;
  /** Error Message: string contains error message if status is 'error' (read only) */
  error?: string;
  /**
   * Backend Type: string (set once)
   * Match a backend from /features/sinks.
   * Cannot change once created (read only)
   */
  backend?: string;
  /** config: object containing sink specific info */
  config?: SinkConfig<string>;
  /** ts_created: UUIDv4 (read only) */
  ts_created?: string;
}

/**
 * Prometheus Sink Type
 */
export type PromSink = Sink | {
  config?: PrometheusConfig;
};

/**
 * for future
 * Sink<T> = {..., config?: <T>, ...}
 * or
 * SpecificSink = Sink | {//{[overrides: string]: any};
 * or
 * SpecificSink extends PrometheusConfig;
 */
