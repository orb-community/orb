/**
 * Sink Data Model Interface
 *
 * [Sinks Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Sinks}
 */

import { PrometheusConfig } from 'app/common/interfaces/orb/sink.config/prometheus.config.interface';
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';

/**
 * @interface Sink
 */
export interface Sink {
  /**
   *  ID {string} UUIDv4 (read only)
   */
  id?: string;

  /**
   * Name {string} [a-zA-Z_:][a-zA-Z0-9_]*
   */
  name?: string;

  /**
   * Description {string}
   */
  description?: string;

  /**
   * A timestamp of creation {string}
   */
  ts_created?: string;

  /**
   * Tags {{[propName: string]: string}}
   */
  tags?: any;

  /**
   *  Status: {string} = 'active'|'error'
   */
  status?: string;

  /**
   * Error Message: {string}
   * Contains error message if status is 'error' (read only)
   */
  error?: string;

  /**
   * Backend Type: {string}
   * Match a backend from /features/sinks.
   * Cannot change once created (read only)
   */
  backend?: string;

  /**
   * Sink Config {{[propName: string]: string}}
   * config: object containing sink specific info
   */
  config?: SinkConfig<string>|any;
}

/**
 * Prometheus Sink Type
 * @type PromSink
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
