/**
 * Sink Data Model Interface
 *
 * [Sinks Architecture]{@link https://github.com/orb-community/orb/wiki/Architecture:-Sinks}
 */


import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';
import { OtlpConfig } from './sink/config/otlp.config.interface';

/**
 * @enum SinkStates
 */
export enum SinkStates {
  active = 'active',
  error = 'error',
  idle = 'idle',
  unknown = 'unknown',
}

/**
 * @enum SinkBackends
 */
export enum SinkBackends {
  prometheus = 'prometheus',
  otlp = 'otlphttp'
}

/**
 * @interface Sink
 */
export interface Sink extends OrbEntity {
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
   *  State: {string} = 'active'|'error'
   */
  state?: string;

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
  config?: SinkTypes;

  /**
   *  Default = json, can be Yaml
   */
  format?: string;
  /**
   *  Only used for Yaml payload
   */
  config_data?: string;
}

export type SinkTypes = OtlpConfig;

/**
 * Prometheus Sink Type
 * @type OtlpSink
 */
export type OtlpSink =
  | Sink
  | {
      config?: OtlpConfig;
    };

/**
 * for future
 * Sink<T> = {..., config?: <T>, ...}
 * or
 * SpecificSink = Sink | {//{[overrides: string]: any};
 * or
 * SpecificSink extends PrometheusConfig;
 */
