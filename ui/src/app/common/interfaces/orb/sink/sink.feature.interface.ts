/**
 * Base Sink Feature Interface
 *
 * [Sinks Architecture]{@link https://github.com/etaques/orb/wiki/Architecture:-Sinks}
 */

import { DynamicFormConfig } from 'app/common/interfaces/orb/dynamic.form.interface';

/**
 * @interface SinkFeature
 */
export interface SinkFeature {
  /**
   * Backend name {string}
   */
  backend?: string;

  /**
   * Backend description {string}
   */
  description?: string;

  /**
   * Backend config {DynamicFormConfig[]}
   */
  config?: DynamicFormConfig[];
}
