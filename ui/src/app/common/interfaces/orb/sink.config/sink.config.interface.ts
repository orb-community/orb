/**
 * Base Sink Config Interface
 * for more details:
 * https://github.com/ns1labs/orb/wiki/Architecture:-Sinks
 */
export interface SinkConfig<T> {
  [propName: string]: T;
}
