/**
 * Base Sink Config Interface
 *
 * [Sinks Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Sinks}
 */

/**
 * @interface SinkConfig
 */
export interface SinkConfig<T> {
  /**
   * propName {string}: <T>value
   */
  [propName: string]: T;
}
