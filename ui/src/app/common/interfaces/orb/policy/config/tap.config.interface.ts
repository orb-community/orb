/**
 * Agent Policy / Tap Config Interface
 *
 * [Policies Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 * [Agent Taps](https://github.com/ns1labs/pktvisor/blob/develop/RFCs/2021-04-16-75-taps.md)
 */

/**
 * @interface TapConfig
 */
export interface TapConfig {
  /**
   * version {string}
   */
  version?: string;

  /**
   * info reserved
   */
  info?: any;

  /**
   * json object with configs
   */
  config?: {[propName: string]: any};
}

